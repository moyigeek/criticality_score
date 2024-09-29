<!--
 * @Author: 7erry
 * @Date: 2024-08-31 03:16:18
 * @LastEditTime: 2024-09-29 17:52:24
 * @Description: 
-->
# About

An concurrent and efficient Go program which serves to collect Git Repo Metrics for [criticality score](https://github.com/ossf/criticality_score).

## Layout

```plaintext
├─cmd
│  ├─Cli --- Cli for collector
│  ├─Clone --- Clone repos from corresponding .csv file
│  ├─Collect --- Collect metrics from corresponding .csv file and store them into postgreSQL
│  └─CountDownloaded --- Count How many repos downloaded
├─config --- Customize relative configs
├─docs --- To Do
├─examples --- To Do
├─images --- To Do
├─input --- Default input path
│  └─tmp
├─internal --- modules of this collector
│  ├─collector --- Collect / Clone related
│  ├─io
│  │  ├─database
│  │  │  ├─psql --- i/o for psql
│  │  │  └─sqlite --- i/o for sqlite
│  │  └─file
│  │   ├─csv --- i/o for csv
│  │   ├─json --- i/o for json
│  │   └─yaml --- i/o for yaml
│  ├─logger
│  ├─parser --- parser of collector
│  │  ├─git --- parse git repo
│  │  └─url --- parse url
│  ├─utils --- basic utils for collector
│  └─workerpool --- workerpool implementation of collector
├─output -- default output path
├─scripts -- some scripts to help run this collector in multi-process way
│  └─ez_scripts -- python scripts
└─storage -- default storage path
```

## Usage

### Cli

```sh
./Cli [url]
```

With this command, collector will collect [url] and print metrics to standard output

### Update config/config.go

Before use the following functions, update your config/config.go

``` Go
// I/O Config
INPUT_CSV_PATH  string = "./input/input.csv"
OUTPUT_CSV_PATH string = "./output/output.csv"
STORAGE_PATHstring = "./storage"

// Database Config

PSQL_HOST  string = ""
PSQL_USER  string = ""
PSQL_PASSWORD  string = ""
PSQL_DATABASE_NAME string = ""
PSQL_PORT  string = ""
PSQL_SSL_MODE  string = ""

SQLITE_DATABASE_PATH string = "./output/test.db"
SQLITE_USER  string = ""
SQLITE_PASSWORD  string = ""
```

If needed, you can also update /internal/io/database/database.go and /internal/parser/parser.go

### Clone

```sh
./Clone [csv_file_path]
```

The collector will read urls from the csv file and download them. The log info will be printed to the standard output

.eg

```sh
./Clone ../../input/input.csv
```

### Collect

```sh
./Collect [csv_file_path]
```

The collector will read urls/path from the csv file and collect metrics from corresponding repo. The log info will be printed to the standard output

.eg

```sh
./Collect ../../input/input.csv
```

### CountDownloaded

```sh
./CountDownloaded [[csv_file_path]]
```

Count how many repos of the csv file have been downloaded. May it should be placed in /scripts

.eg

```sh
./CountDownloaded ../../input/input.csv
```
