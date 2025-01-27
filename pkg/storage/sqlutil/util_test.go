package sqlutil

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/samber/lo"
)

// private function test

func TestCamelToSnake(t *testing.T) {
	tests := []struct {
		arg  string
		want string
	}{
		{arg: "CamelToSnake", want: "camel_to_snake"},
		{arg: "ID", want: "id"},
		{arg: "SQLite", want: "sqlite"},
		{arg: "WindowsVista", want: "windows_vista"},
		{arg: "Windows7", want: "windows7"},
		{arg: "test_EventEmitter", want: "test_event_emitter"},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("TestingCamelToSnake%d", i), func(t *testing.T) {
			if got := camelToSnake(tt.arg); got != tt.want {
				t.Errorf("camelToSnake() = %v, want %v", got, tt.want)
			}
		})
	}
}

type a struct {
	ID               *int
	Name             *string
	SomeStrangeField *string `column:"abcdeSSSS"`
}

type b struct {
	ID    *int `generated:"true"`
	Type  *int `column:"type" pk:"true"`
	Event *string
}

type queryBuilderWant struct {
	query string
	args  []interface{}
}

func assertQueryBuilderWant(t *testing.T, query string, args []interface{}, err error, wantErr bool, want []queryBuilderWant) {
	if wantErr && err != nil {
		return
	}

	succeed := false
	for _, want := range want {
		if query == want.query {
			for i, arg := range args {
				if arg != want.args[i] {
					t.Errorf("GetUpdateQueryAndArgs(), args = %v, want %v", args, want.args)
					return
				}
			}
			succeed = true
			break
		}
	}
	if !succeed {
		t.Errorf("GetUpdateQueryAndArgs() = %v,%v,%v, want %v", query, args, err, want[0])
	}
}

func TestInsertSentence(t *testing.T) {
	type testArgs struct {
		tableName string
		data      any
	}

	tests := []struct {
		args    testArgs
		want    []queryBuilderWant
		wantErr bool
	}{
		{
			args: testArgs{
				tableName: "table",
				data:      &a{ID: nil, Name: nil, SomeStrangeField: nil},
			},
			wantErr: true,
		},
		{
			args: testArgs{
				tableName: "table",
				data:      &a{ID: nil, Name: nil, SomeStrangeField: lo.ToPtr("test")},
			},
			want: []queryBuilderWant{
				{"INSERT INTO table (abcdeSSSS) VALUES ($1)", []any{"test"}},
			},
		},
		{
			args: testArgs{
				tableName: "table",
				data:      &a{ID: lo.ToPtr(1), Name: lo.ToPtr("test"), SomeStrangeField: lo.ToPtr("test")},
			},
			want: []queryBuilderWant{
				{"INSERT INTO table (id, name, abcdeSSSS) VALUES ($1, $2, $3)", []any{1, "test", "test"}},
				{"INSERT INTO table (id, abcdeSSSS, name) VALUES ($1, $2, $3)", []any{1, "test", "test"}},
				{"INSERT INTO table (name, id, abcdeSSSS) VALUES ($1, $2, $3)", []any{"test", 1, "test"}},
				{"INSERT INTO table (name, abcdeSSSS, id) VALUES ($1, $2, $3)", []any{"test", "test", 1}},
				{"INSERT INTO table (abcdeSSSS, id, name) VALUES ($1, $2, $3)", []any{"test", 1, "test"}},
				{"INSERT INTO table (abcdeSSSS, name, id) VALUES ($1, $2, $3)", []any{"test", "test", 1}},
			},
		},
		{
			args: testArgs{
				tableName: "table",
				data:      &b{ID: lo.ToPtr(1), Type: lo.ToPtr(2), Event: lo.ToPtr("test")},
			},
			want: []queryBuilderWant{
				{"INSERT INTO table (type, event) VALUES ($1, $2)", []any{2, "test"}},
				{"INSERT INTO table (event, type) VALUES ($1, $2)", []any{"test", 2}},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("TestInsertSentence%d", i), func(t *testing.T) {
			var got string
			var args []interface{}
			var err error
			switch tt.args.data.(type) {
			case *a:
				got, args, err = getInsertQueryAndArgs(tt.args.tableName, tt.args.data.(*a), false)
			case *b:
				got, args, err = getInsertQueryAndArgs(tt.args.tableName, tt.args.data.(*b), false)
			}
			assertQueryBuilderWant(t, got, args, err, tt.wantErr, tt.want)
		})
	}
}

func TestUpdateSentence(t *testing.T) {
	type testArgs struct {
		tableName string
		data      *a
	}

	tests := []struct {
		args    testArgs
		want    []queryBuilderWant
		wantErr bool
	}{
		{
			args: testArgs{
				tableName: "table",
				data:      &a{ID: nil, Name: nil, SomeStrangeField: nil},
			},
			wantErr: true,
		},
		{
			args: testArgs{
				tableName: "table",
				data:      &a{ID: nil, Name: nil, SomeStrangeField: lo.ToPtr("test")},
			},
			wantErr: true,
		},
		{
			args: testArgs{
				tableName: "table",
				data:      &a{ID: lo.ToPtr(1), SomeStrangeField: lo.ToPtr("test")},
			},
			want: []queryBuilderWant{
				{"UPDATE table SET abcdeSSSS = $1 WHERE id = $2", []interface{}{"test", 1}},
			},
		},
		{
			args: testArgs{
				tableName: "table",
				data:      &a{ID: lo.ToPtr(1), Name: lo.ToPtr("testName"), SomeStrangeField: lo.ToPtr("testSomeStrangeField")},
			},
			want: []queryBuilderWant{
				{"UPDATE table SET name = $1, abcdeSSSS = $2 WHERE id = $3", []interface{}{"testName", "testSomeStrangeField", 1}},
				{"UPDATE table SET abcdeSSSS = $1, name = $2 WHERE id = $3", []interface{}{"testSomeStrangeField", "testName", 1}},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("TestUpdateSentence%d", i), func(t *testing.T) {
			query, args, err := getUpdateQueryAndArgs(tt.args.tableName, tt.args.data)
			assertQueryBuilderWant(t, query, args, err, tt.wantErr, tt.want)
		})
	}

}

func isStructEqual[T any](a, b *T) bool {
	// every field .Elem() same then equal
	reflectType := reflect.TypeOf(*a)
	reflectValA := reflect.ValueOf(a).Elem()
	reflectValB := reflect.ValueOf(b).Elem()

	for i := 0; i < reflectType.NumField(); i++ {
		if reflectValA.Field(i).Elem().Interface() != reflectValB.Field(i).Elem().Interface() {
			return false
		}
	}
	return true
}

func TestRowsToEntity(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT").WillReturnRows(mock.NewRows([]string{"id", "name", "abcdeSSSS"}).AddRow(1, "test", "test"))

	rows, err := db.Query("SELECT")
	if err != nil {
		t.Fatalf("failed to query: %v", err)
	}
	if !rows.Next() {
		t.Fatalf("no rows")
	}
	got, err := rowsToEntity[a](rows)
	if err != nil {
		t.Fatalf("failed to rowsToEntity: %v", err)
	}

	want := &a{ID: lo.ToPtr(1), Name: lo.ToPtr("test"), SomeStrangeField: lo.ToPtr("test")}
	if !isStructEqual(got, want) {
		t.Errorf("RowsToEntity() = %v, want %v", got, want)
	}
}

func TestMergeStruct(t *testing.T) {
	type testCase struct {
		oldStruct *a
		newStruct *a
		want      *a
	}

	tests := []testCase{
		{
			oldStruct: &a{ID: lo.ToPtr(1), Name: lo.ToPtr("test"), SomeStrangeField: lo.ToPtr("test")},
			newStruct: &a{ID: lo.ToPtr(2), Name: lo.ToPtr("test"), SomeStrangeField: lo.ToPtr("test")},
			want:      &a{ID: lo.ToPtr(2), Name: lo.ToPtr("test"), SomeStrangeField: lo.ToPtr("test")},
		},
		{
			oldStruct: &a{ID: nil, Name: nil, SomeStrangeField: lo.ToPtr("test")},
			newStruct: &a{ID: lo.ToPtr(2), Name: lo.ToPtr("test"), SomeStrangeField: lo.ToPtr("test2")},
			want:      &a{ID: lo.ToPtr(2), Name: lo.ToPtr("test"), SomeStrangeField: lo.ToPtr("test2")},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("TestMergeStruct%d", i), func(t *testing.T) {
			MergeStruct(tt.oldStruct, tt.newStruct)
			if !isStructEqual(tt.newStruct, tt.want) {
				t.Errorf("MergeStruct() = %v, want %v", tt.newStruct, tt.want)
			}
		})
	}
}
