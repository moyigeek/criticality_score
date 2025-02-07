package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/sirupsen/logrus"
	"go.elastic.co/ecslogrus"
	"gopkg.in/go-extras/elogrus.v8"
)

type cliFormatter struct{}

func (f *cliFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	if entry.Level == logrus.WarnLevel {
		return []byte(fmt.Sprintf("\033[33m%s\033[0m\n", entry.Message)), nil
	} else if entry.Level <= logrus.ErrorLevel {
		return []byte(fmt.Sprintf("\033[31m%s\033[0m\n", entry.Message)), nil
	}

	return []byte(fmt.Sprintf("%s\n", entry.Message)), nil
}

func ECSLogMessageModifierFunc(formatter *ecslogrus.Formatter) func(*logrus.Entry, *elogrus.Message) any {
	return func(entry *logrus.Entry, message *elogrus.Message) any {
		var data json.RawMessage
		data, err := formatter.Format(entry)
		if err != nil {
			return entry // in case of an error just preserve the original entry
		}
		return data
	}

}

type logrusLogger struct {
	logger *logrus.Logger
	entry  *logrus.Entry
}

func NewLogrusLogger(config *AppLoggerConfig) AppLogger {
	if config == nil {
		config = &AppLoggerConfig{
			Level:      LoggerLevelInfo,
			FormatType: LoggerFormatText,
			Output:     LoggerOutputStderr,
		}
	}

	l := logrus.New()

	logger := &logrusLogger{
		logger: l,
		entry:  logrus.NewEntry(l),
	}
	logger.Config(config)
	return logger
}

func (l *logrusLogger) Config(config *AppLoggerConfig) error {
	level, err := logrus.ParseLevel(string(config.Level))
	if err != nil {
		level = logrus.InfoLevel
	}

	l.logger.SetLevel(level)
	var formatter logrus.Formatter

	switch config.FormatType {
	case LoggerFormatText:
		formatter = &logrus.TextFormatter{}
	case LoggerFormatCliTool:
		formatter = &cliFormatter{}
	case LoggerFormatJSON:
		formatter = &logrus.JSONFormatter{}
	default:
		formatter = &logrus.TextFormatter{}
	}

	l.logger.SetFormatter(formatter)

	switch config.Output {
	case LoggerOutputStdout:
		l.logger.SetOutput(os.Stdout)
	case LoggerOutputFile:
		file, err := os.OpenFile(config.OutputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
		if err != nil {
			panic(fmt.Sprintf("Failed to open log file: %v", err))
		}
		l.logger.SetOutput(file)
	case LoggerOutputStderr:
		l.logger.SetOutput(os.Stderr)
	case LoggerOutputElasticSearch:
		var caCert []byte = nil
		if config.OutputEsCert != "" {
			caCert, err = os.ReadFile(config.OutputEsCert)
			if err != nil {
				return err
			}
		}

		client, err := elasticsearch.NewClient(elasticsearch.Config{
			Addresses: []string{config.OutputEsURL},
			Username:  config.OutputEsUser,
			Password:  config.OutputEsPassword,
			CACert:    caCert,
		})
		if err != nil {
			return err
		}
		hook, err := elogrus.NewAsyncElasticHook(client, "localhost", level, config.OutputEsIndex)
		if err != nil {
			return err
		}
		hook.MessageModifierFunc = ECSLogMessageModifierFunc(&ecslogrus.Formatter{})
		l.logger.SetReportCaller(true)
		l.logger.AddHook(hook)
		l.logger.SetOutput(io.Discard)
	default:
		l.logger.SetOutput(os.Stderr)
	}
	return nil
}

func (l *logrusLogger) Trace(args ...interface{}) {
	l.entry.Trace(args...)
}

func (l *logrusLogger) Debug(args ...interface{}) {
	l.entry.Debug(args...)
}

func (l *logrusLogger) Info(args ...interface{}) {
	l.entry.Info(args...)
}

func (l *logrusLogger) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}

func (l *logrusLogger) Error(args ...interface{}) {
	l.entry.Error(args...)
}

func (l *logrusLogger) Fatal(args ...interface{}) {
	l.entry.Fatal(args...)
}

func (l *logrusLogger) Panic(args ...interface{}) {
	l.entry.Panic(args...)
}

func (l *logrusLogger) Print(args ...interface{}) {
	l.entry.Print(args...)
}

func (l *logrusLogger) Tracef(format string, args ...interface{}) {
	l.entry.Tracef(format, args...)
}

func (l *logrusLogger) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

func (l *logrusLogger) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

func (l *logrusLogger) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

func (l *logrusLogger) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

func (l *logrusLogger) Panicf(format string, args ...interface{}) {
	l.entry.Panicf(format, args...)
}

func (l *logrusLogger) Printf(format string, args ...interface{}) {
	l.entry.Printf(format, args...)
}

func (l *logrusLogger) Debugln(args ...interface{}) {
	l.entry.Debugln(args...)
}

func (l *logrusLogger) Traceln(args ...interface{}) {
	l.entry.Traceln(args...)
}

func (l *logrusLogger) Infoln(args ...interface{}) {
	l.entry.Infoln(args...)
}

func (l *logrusLogger) Warnln(args ...interface{}) {
	l.entry.Warnln(args...)
}

func (l *logrusLogger) Errorln(args ...interface{}) {
	l.entry.Errorln(args...)
}

func (l *logrusLogger) Fatalln(args ...interface{}) {
	l.entry.Fatalln(args...)
}

func (l *logrusLogger) Panicln(args ...interface{}) {
	l.entry.Panicln(args...)
}

func (l *logrusLogger) Println(args ...interface{}) {
	l.entry.Println(args...)
}

func (l *logrusLogger) WithFields(fields map[string]interface{}) AppLogger {
	return &logrusLogger{
		logger: l.logger,
		entry:  l.entry.WithFields(fields),
	}
}
