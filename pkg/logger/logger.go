package logger

type AppLogger interface {
	Config(config *AppLoggerConfig) error

	Trace(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
	Print(args ...interface{})

	Traceln(args ...interface{})
	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Warnln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})
	Println(args ...interface{})

	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Printf(format string, args ...interface{})

	WithFields(fields map[string]interface{}) AppLogger
}

type AppLoggerConfig struct {
	Level      LoggerLevel
	FormatType LoggerFormatType
	Output     LoggerOutput
	// the path to the log file, take effect only when output = LoggerOutputFile
	OutputPath string
	// the elastic search url, take effect only when output = LoggerOutputElasticSearch
	OutputEsURL string
	// the elastic search index, take effect only when output = LoggerOutputElasticSearch
	OutputEsIndex string
}

type LoggerOutput int

const (
	LoggerOutputStdout LoggerOutput = iota
	LoggerOutputStderr
	LoggerOutputFile
	LoggerOutputElasticSearch
)

type LoggerLevel string

const (
	LoggerLevelTrace LoggerLevel = "trace"
	LoggerLevelDebug             = "debug"
	LoggerLevelInfo              = "info"
	LoggerLevelWarn              = "warn"
	LoggerLevelError             = "error"
	LoggerLevelFatal             = "fatal"
	LoggerLevelPanic             = "panic"
)

type LoggerFormatType int

const (
	LoggerFormatText LoggerFormatType = iota
	LoggerFormatCliTool
	LoggerFormatJSON
)
