package storage

import (
	"database/sql"
	"fmt"

	"github.com/HUSTSecLab/criticality_score/pkg/logger"
	_ "github.com/lib/pq"
)

type AppDatabaseContext interface {
	GetConfig() Config
	SetSQLLog(enable bool)
	NewBatchExecContext(config *BatchExecContextConfig) BatchExecContext
	GetDatabaseConnection() (*sql.DB, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Close() error
}

type appDatabaseContext struct {
	config       Config
	enableSQLLog bool
	db           *sql.DB
}

func NewAppDatabase(configPath string) (AppDatabaseContext, error) {
	config, err := loadConfig(configPath)
	if err != nil {
		fmt.Errorf("Failed to load config:", err)
		return nil, err
	}
	return &appDatabaseContext{config: config}, nil
}

func NewAppDatabaseWithDb(db *sql.DB) AppDatabaseContext {
	return &appDatabaseContext{config: Config{}, db: db}
}

func (appDb *appDatabaseContext) ensureDatabaseConnection() error {
	if appDb.db == nil {
		connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			appDb.config.Host, appDb.config.Port, appDb.config.User, appDb.config.Password, appDb.config.Database)
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			return err
		}
		appDb.db = db
	}
	return nil
}

func (appDb *appDatabaseContext) GetConfig() Config {
	return appDb.config
}

func (appDb *appDatabaseContext) SetSQLLog(enable bool) {
	appDb.enableSQLLog = enable
}

func (appDb *appDatabaseContext) NewBatchExecContext(config *BatchExecContextConfig) BatchExecContext {
	return &batchExecContext{
		appDb:          appDb,
		config:         config,
		sentences:      "",
		sentencesCount: 0,
		args:           make([]interface{}, 0),
	}
}

func (app *appDatabaseContext) GetDatabaseConnection() (*sql.DB, error) {
	err := app.ensureDatabaseConnection()
	return app.db, err
}

func (app *appDatabaseContext) Exec(query string, args ...interface{}) (sql.Result, error) {
	if app.enableSQLLog {
		logger.Info("Exec SQL: ", query)
	}

	conn, err := app.GetDatabaseConnection()
	if err != nil {
		return nil, err
	}
	return conn.Exec(query, args...)
}

func (app *appDatabaseContext) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if app.enableSQLLog {
		logger.Info("Query SQL: ", query)
	}

	conn, err := app.GetDatabaseConnection()
	if err != nil {
		return nil, err
	}
	return conn.Query(query, args...)
}

func (app *appDatabaseContext) QueryRow(query string, args ...interface{}) *sql.Row {
	if app.enableSQLLog {
		logger.Info("QueryRow SQL: ", query)
	}

	conn, err := app.GetDatabaseConnection()
	if err != nil {
		return nil
	}
	return conn.QueryRow(query, args...)
}

func (app *appDatabaseContext) Close() error {
	if app.db != nil {
		return app.db.Close()
	}
	return nil
}
