# About

An concurrent and efficient Go program which serves to collect Git Repo Metrics for criticality score, based on [go-git](https://github.com/go-git/go-git)

## Layout

```plaintext
├─cmd
│  ├─Cli                    Cli collector
│  ├─Clone                  Clone repos
│  ├─Collect                Clone and Collect remote repos or Collect local repos
│  ├─CountDownloaded        Count how many of the repos are downloaded
│  └─integrate              CI unit
├─config                    Basic Config file
├─input                     Default input directory
├─internal                  Libraries
│  ├─collector              Collector for git repos
│  ├─io
│  │  ├─database            I/O with database
│  │  └─file                I/O with files like .json, .csv, .yaml
│  ├─logger                 Log runtime message
│  └─parser
│      ├─git                Parser for git repos
│      └─url                Parser for given urls
├─output                    Default output directory
└─storage                   Default storage directory
```

## Quick Start

### Configure

Before use the following functions, update your `config/config.go`

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

If needed, you can also update `internal/io/database/database.go` and `internal/parser/parser.go` to make further customization.

### Cli

```sh
./Cli [url]
```

With this command, collector will collect  `url` and print metrics to standard output

### Clone

```sh
./Clone [csv_file_path]
```

The collector will read urls from the csv file and download them.

### Collect

```sh
./Collect [csv_file_path]
```

The collector will read urls/path from the csv file and collect metrics from corresponding repo. And the collected metrics will be stored at the database
