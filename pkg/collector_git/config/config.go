/*
 * @Author: 7erry
 * @Date: 2024-08-31 03:44:46
 * @LastEditTime: 2024-12-14 16:48:05
 * @Description: Config file for the git metrics collector
 */
package config

import "github.com/sirupsen/logrus"

const (
	// Log Config
	LOG_LEVEL = logrus.InfoLevel

	// I/O Config
	INPUT_CSV_PATH  string = "./input/input.csv"
	OUTPUT_CSV_PATH string = "./output/output.csv"

	// Database Config
	BATCH_SIZE                = 256
	PSQL_HOST          string = ""
	PSQL_USER          string = ""
	PSQL_PASSWORD      string = ""
	PSQL_DATABASE_NAME string = ""
	PSQL_PORT          string = ""
	PSQL_SSL_MODE      string = ""

	SQLITE_DATABASE_PATH string = "./output/test.db"
	SQLITE_USER          string = ""
	SQLITE_PASSWORD      string = ""
)

// The cloned repositories will be stored in STORAGE_PATH
var STORAGE_PATH string = "./storage"

func SetStoragetPath(path string) {
	STORAGE_PATH = path
}
