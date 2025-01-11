// This package contains Sql utility functions for database operations
package sqlutil

import (
	"database/sql"
	"fmt"
	"iter"
	"reflect"
	"regexp"
	"strings"

	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	"github.com/samber/lo"
)

func camelToSnake(s string) string {
	pattern := regexp.MustCompile("(\\p{Lu}+\\P{Lu}*)")
	s = pattern.ReplaceAllString(s, "${1}_")
	s, _ = strings.CutSuffix(strings.ToLower(s), "_")
	return s
}

func getFieldColumnName(f reflect.StructField) string {
	column := f.Tag.Get("column")
	if column == "" {
		// transform camel case to snake case
		column = camelToSnake(f.Name)
	}
	return column
}

var typeColumnToFieldIdxMapCache = make(map[reflect.Type]map[string]int)

func getTypeColumnToFieldIdxMap(t reflect.Type) map[string]int {
	// if t is pointer, get the element type
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	// first check cache
	if ret, ok := typeColumnToFieldIdxMapCache[t]; ok {
		return ret
	}

	// build map
	ret := make(map[string]int)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.IsExported() != false && field.Tag.Get("ignore") != "true" {
			ret[getFieldColumnName(field)] = i
		}
	}

	// save to cache
	typeColumnToFieldIdxMapCache[t] = ret
	return ret
}

var typePrimaryKeyCache = make(map[reflect.Type][]string)

func getTypePrimaryKey(t reflect.Type) []string {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if ret, ok := typePrimaryKeyCache[t]; ok {
		return ret
	}

	pkColumn := make([]string, 0)
	fallBack := ""
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.IsExported() == false || field.Tag.Get("ignore") == "true" {
			continue
		}
		if field.Tag.Get("pk") == "true" {
			column := field.Tag.Get("column")
			if column == "" {
				column = camelToSnake(field.Name)
			}
			pkColumn = append(pkColumn, column)
		}
		if field.Name == "ID" || field.Name == "Id" {
			fallBack = field.Name
		}
	}

	if len(pkColumn) == 0 && fallBack != "" {
		pkColumn = append(pkColumn, camelToSnake(fallBack))
	}

	typePrimaryKeyCache[t] = pkColumn
	return pkColumn
}

func rowsToEntity[T any](rows *sql.Rows) (*T, error) {
	var val T
	reflectType := reflect.TypeOf(val)
	reflectVal := reflect.New(reflectType).Elem()
	addressList := make([]interface{}, 0)

	cToFMap := getTypeColumnToFieldIdxMap(reflectType)

	columnNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	for _, columnName := range columnNames {
		fieldIdx, ok := cToFMap[columnName]
		if !ok { // ignore column
			var ignore interface{}
			addressList = append(addressList, &ignore)
			continue
		}
		fieldVal := reflectVal.Field(fieldIdx)
		newObj := reflect.New(fieldVal.Type().Elem())
		fieldVal.Set(newObj)
		addressList = append(addressList, newObj.Interface())
	}

	err = rows.Scan(addressList...)
	if err != nil {
		return nil, err
	}
	return reflectVal.Addr().Interface().(*T), nil
}

func createIterator[T any](rows *sql.Rows) iter.Seq[*T] {
	return func(yield func(*T) bool) {
		defer rows.Close()
		for rows.Next() {
			pval, err := rowsToEntity[T](rows)
			if err != nil {
				return
			}
			if !yield(pval) {
				return
			}
		}
	}
}

func getInsertQueryAndArgs[T any](tableName string, data *T) (string, []interface{}, error) {
	reflectType := reflect.TypeOf(*data)
	reflectVal := reflect.ValueOf(data).Elem()

	columns := make([]string, 0)
	values := make([]interface{}, 0)

	cToFMap := getTypeColumnToFieldIdxMap(reflectType)

	for k, v := range cToFMap {
		if reflectVal.Field(v).IsNil() {
			continue
		}
		columns = append(columns, k)
		values = append(values, reflectVal.Field(v).Elem().Interface())
	}

	if len(columns) == 0 {
		return "", nil, fmt.Errorf("no column to insert")
	}

	columnsStr := strings.Join(columns, ", ")
	// value str $1, $2
	valuesArr := make([]string, 0)
	for i := 1; i <= len(values); i++ {
		valuesArr = append(valuesArr, fmt.Sprintf("$%d", i))
	}
	valuesStr := strings.Join(valuesArr, ", ")

	insertSentenceTemplate := `INSERT INTO %s (%s) VALUES (%s)`
	insertSentence := fmt.Sprintf(insertSentenceTemplate, tableName, columnsStr, valuesStr)

	return insertSentence, values, nil
}

func getUpdateQueryAndArgs[T any](tableName string, data *T) (string, []interface{}, error) {
	reflectType := reflect.TypeOf(*data)
	reflectVal := reflect.ValueOf(data).Elem()

	columns := make([]string, 0)
	values := make([]interface{}, 0)

	whereColumns := make([]string, 0)
	whereValues := make([]interface{}, 0)

	cToFMap := getTypeColumnToFieldIdxMap(reflectType)
	pkColumns := getTypePrimaryKey(reflectType)

	for k, v := range cToFMap {
		if reflectVal.Field(v).IsNil() {
			continue
		}
		if lo.IndexOf(pkColumns, k) != -1 {
			// primary key column should not be updated
			continue
		}
		columns = append(columns, k)
		values = append(values, reflectVal.Field(v).Elem().Interface())
	}

	if len(columns) == 0 {
		return "", nil, fmt.Errorf("no column to update")
	}

	queryPlacement := 1
	valuesArr := make([]string, 0)
	for _, col := range columns {
		valuesArr = append(valuesArr, fmt.Sprintf("%s = $%d", col, queryPlacement))
		queryPlacement++
	}
	valuesStr := strings.Join(valuesArr, ", ")

	whereArr := make([]string, 0)
	for _, pkColumn := range pkColumns {
		whereArr = append(whereArr, fmt.Sprintf("%s = $%d", pkColumn, queryPlacement))
		queryPlacement++
	}
	whereStr := strings.Join(whereArr, " AND ")

	for _, pkColumn := range pkColumns {
		fieldIdx, ok := cToFMap[pkColumn]
		if !ok {
			return "", nil, fmt.Errorf("primary key column %s not found in struct", pkColumn)
		}
		whereColumns = append(whereColumns, pkColumn)
		if reflectVal.Field(fieldIdx).IsNil() {
			return "", nil, fmt.Errorf("primary key column %s is nil", pkColumn)
		}
		whereValues = append(whereValues, reflectVal.Field(fieldIdx).Elem().Interface())
	}

	values = append(values, whereValues...)

	updateSentenceTemplate := `UPDATE %s SET %s WHERE %s`
	updateSentence := fmt.Sprintf(updateSentenceTemplate, tableName, valuesStr, whereStr)

	return updateSentence, values, nil
}

func getDeleteQueryAndArgs[T any](tableName string, data *T) (string, []interface{}, error) {
	reflectType := reflect.TypeOf(*data)
	reflectVal := reflect.ValueOf(data).Elem()

	whereColumns := make([]string, 0)
	whereValues := make([]interface{}, 0)

	cToFMap := getTypeColumnToFieldIdxMap(reflectType)
	pkColumns := getTypePrimaryKey(reflectType)

	queryPlacement := 1
	whereArr := make([]string, 0)
	for _, pkColumn := range pkColumns {
		whereArr = append(whereArr, fmt.Sprintf("%s = $%d", pkColumn, queryPlacement))
		queryPlacement++
	}
	whereStr := strings.Join(whereArr, " AND ")

	for _, pkColumn := range pkColumns {
		fieldIdx, ok := cToFMap[pkColumn]
		if !ok {
			return "", nil, fmt.Errorf("primary key column %s not found in struct", pkColumn)
		}
		whereColumns = append(whereColumns, pkColumn)
		if reflectVal.Field(fieldIdx).IsNil() {
			return "", nil, fmt.Errorf("primary key column %s is nil", pkColumn)
		}
		whereValues = append(whereValues, reflectVal.Field(fieldIdx).Elem().Interface())
	}

	deleteSentenceTemplate := `DELETE FROM %s WHERE %s`
	deleteSentence := fmt.Sprintf(deleteSentenceTemplate, tableName, whereStr)

	return deleteSentence, whereValues, nil
}

func Query[T any](ctx storage.AppDatabaseContext, query string, args ...interface{}) (iter.Seq[*T], error) {
	rows, err := ctx.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return createIterator[T](rows), nil
}

func QueryFirst[T any](ctx storage.AppDatabaseContext, query string, args ...interface{}) (*T, error) {
	rows, err := ctx.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		return rowsToEntity[T](rows)
	}
	return nil, nil
}

func getSelectQuery[T any](from string, afterFrom string) string {
	// data is only use for get the type of T
	data := new(T)
	reflectType := reflect.TypeOf(*data)

	cToFMap := getTypeColumnToFieldIdxMap(reflectType)

	fields := make([]string, 0)

	for k, _ := range cToFMap {
		fields = append(fields, k)
	}

	query := fmt.Sprintf(`SELECT %s FROM %s %s`, strings.Join(fields, ", "), from, afterFrom)
	return query
}

func QueryCommon[T any](ctx storage.AppDatabaseContext, from string, afterFrom string, args ...interface{}) (iter.Seq[*T], error) {
	query := getSelectQuery[T](from, afterFrom)
	return Query[T](ctx, query, args...)
}

func QueryCommonFirst[T any](ctx storage.AppDatabaseContext, from string, afterFrom string, args ...interface{}) (*T, error) {
	query := getSelectQuery[T](from, afterFrom+" LIMIT 1")
	return QueryFirst[T](ctx, query, args...)
}

func Insert[T any](ctx storage.AppDatabaseContext, into string, data *T) error {
	insertSentence, values, err := getInsertQueryAndArgs[T](into, data)
	if err != nil {
		return err
	}
	_, err = ctx.Exec(insertSentence, values...)
	return err
}

func BatchInsert[T any](ctx storage.AppDatabaseContext, into string, data []*T) error {
	batchCtx := ctx.NewBatchExecContext(&storage.BatchExecContextConfig{
		AutoCommit:     true,
		AutoCommitSize: 1000,
	})
	defer batchCtx.Commit()
	for _, d := range data {
		insertSentence, args, err := getInsertQueryAndArgs[T](into, d)
		if err != nil {
			return err
		}
		batchCtx.AppendExec(insertSentence, args...)
	}
	return nil
}

func Update[T any](ctx storage.AppDatabaseContext, tableName string, data *T) error {
	updateSentence, values, err := getUpdateQueryAndArgs[T](tableName, data)
	if err != nil {
		return err
	}
	_, err = ctx.Exec(updateSentence, values...)
	return err
}

func BatchUpdate[T any](ctx storage.AppDatabaseContext, tableName string, data []*T) error {
	batchCtx := ctx.NewBatchExecContext(&storage.BatchExecContextConfig{
		AutoCommit:     true,
		AutoCommitSize: 1000,
	})
	defer batchCtx.Commit()
	for _, d := range data {
		updateSentence, args, err := getUpdateQueryAndArgs[T](tableName, d)
		if err != nil {
			return err
		}
		batchCtx.AppendExec(updateSentence, args...)
	}
	return nil
}

func Delete[T any](ctx storage.AppDatabaseContext, tableName string, data *T) error {
	deleteSentence, values, err := getDeleteQueryAndArgs[T](tableName, data)
	if err != nil {
		return err
	}
	_, err = ctx.Exec(deleteSentence, values...)
	return err
}

// MergeStruct merge old struct to dst struct
// if field in dst is nil, set it to the value of old
func MergeStruct[T any](old *T, dst *T) {
	srcVal := reflect.ValueOf(old).Elem()
	dstVal := reflect.ValueOf(dst).Elem()

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		dstField := dstVal.Field(i)
		if dstVal.IsNil() {
			dstField.Set(srcField)
		}
	}
}
