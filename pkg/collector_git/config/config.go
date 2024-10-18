/*
 * @Author: 7erry
 * @Date: 2024-08-31 03:44:46
 * @LastEditTime: 2024-09-29 17:20:51
 * @Description:
 */
package config

const (
	// I/O Config
	INPUT_CSV_PATH  string = "./input/input.csv"
	OUTPUT_CSV_PATH string = "./output/output.csv"

	// Database Config

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

var STORAGE_PATH string = "./storage"

func SetStoragetPath(path string) {
	STORAGE_PATH = path
}
