package logger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

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
	Level      string
	FormatType LoggerFormatType
	Output     LoggerOutput
	// OutputPath is the path to the log file
	//
	// If Output is LoggerOutputFile, OutputPath is the path to the log file, and
	// if Output is LoggerOutputStdout or LoggerOutputStderr, OutputPath is ignored, can be empty
	OutputPath string
}

type LoggerOutput int

const (
	LoggerOutputStdout LoggerOutput = iota
	LoggerOutputStderr
	LoggerOutputFile
	// Not implemented yet
	LoggerOutputDatabase
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

type cliFormatter struct{}

func (f *cliFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	if entry.Level == logrus.WarnLevel {
		return []byte(fmt.Sprintf("\033[33m%s\033[0m\n", entry.Message)), nil
	} else if entry.Level <= logrus.ErrorLevel {
		return []byte(fmt.Sprintf("\033[31m%s\033[0m\n", entry.Message)), nil
	}

	return []byte(fmt.Sprintf("%s\n", entry.Message)), nil
}

type logrusLogger struct {
	logger *logrus.Logger
	fields map[string]interface{}
}

func NewLogrusLogger(config *AppLoggerConfig) AppLogger {
	if config == nil {
		config = &AppLoggerConfig{
			Level:      string(LoggerLevelInfo),
			FormatType: LoggerFormatText,
			Output:     LoggerOutputStderr,
		}
	}

	logger := &logrusLogger{
		logger: logrus.New(),
		fields: nil,
	}
	logger.Config(config)
	return logger
}

func (l *logrusLogger) Config(config *AppLoggerConfig) error {
	level, err := logrus.ParseLevel(config.Level)
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
		file, err := os.OpenFile(config.OutputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Sprintf("Failed to open log file: %v", err))
		}
		l.logger.SetOutput(file)
	case LoggerOutputStderr:
	default:
		l.logger.SetOutput(os.Stderr)
	}
	return nil
}

func (l *logrusLogger) Trace(args ...interface{}) {
	if l.fields != nil {
		l.logger.WithFields(l.fields).Trace(args...)
	} else {
		l.logger.Trace(args...)
	}
}

func (l *logrusLogger) Debug(args ...interface{}) {
	if l.fields != nil {
		l.logger.WithFields(l.fields).Debug(args...)
	} else {
		l.logger.Debug(args...)
	}
}

func (l *logrusLogger) Info(args ...interface{}) {
	if l.fields != nil {
		l.logger.WithFields(l.fields).Info(args...)
	} else {
		l.logger.Info(args...)
	}
}

func (l *logrusLogger) Warn(args ...interface{}) {
	if l.fields != nil {
		l.logger.WithFields(l.fields).Warn(args...)
	} else {
		l.logger.Warn(args...)
	}
}

func (l *logrusLogger) Error(args ...interface{}) {
	if l.fields != nil {
		l.logger.WithFields(l.fields).Error(args...)
	} else {
		l.logger.Error(args...)
	}
}

func (l *logrusLogger) Fatal(args ...interface{}) {
	if l.fields != nil {
		l.logger.WithFields(l.fields).Fatal(args...)
	} else {
		l.logger.Fatal(args...)
	}
}

func (l *logrusLogger) Panic(args ...interface{}) {
	if l.fields != nil {
		l.logger.WithFields(l.fields).Panic(args...)
	} else {
		l.logger.Panic(args...)
	}
}

func (l *logrusLogger) Print(args ...interface{}) {
	if l.fields != nil {
		l.logger.WithFields(l.fields).Print(args...)
	} else {
		l.logger.Print(args...)
	}
}

func (l *logrusLogger) Tracef(format string, args ...interface{}) {
	if l.fields != nil {
		l.logger.WithFields(l.fields).Tracef(format, args...)
	} else {
		l.logger.Tracef(format, args...)
	}
}

func (l *logrusLogger) Debugf(format string, args ...interface{}) {
	if l.fields != nil {
		l.logger.WithFields(l.fields).Debugf(format, args...)
	} else {
		l.logger.Debugf(format, args...)
	}
}

func (l *logrusLogger) Infof(format string, args ...interface{}) {
	if l.fields != nil {
		l.logger.WithFields(l.fields).Infof(format, args...)
	} else {
		l.logger.Infof(format, args...)
	}
}

func (l *logrusLogger) Warnf(format string, args ...interface{}) {
	if l.fields != nil {
		l.logger.WithFields(l.fields).Warnf(format, args...)
	} else {
		l.logger.Warnf(format, args...)
	}
}

func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	if l.fields != nil {
		l.logger.WithFields(l.fields).Errorf(format, args...)
	} else {
		l.logger.Errorf(format, args...)
	}
}

func (l *logrusLogger) Fatalf(format string, args ...interface{}) {
	if l.fields != nil {
		l.logger.WithFields(l.fields).Fatalf(format, args...)
	} else {
		l.logger.Fatalf(format, args...)
	}
}

func (l *logrusLogger) Panicf(format string, args ...interface{}) {
	if l.fields != nil {
		l.logger.WithFields(l.fields).Panicf(format, args...)
	} else {
		l.logger.Panicf(format, args...)
	}
}

func (l *logrusLogger) Printf(format string, args ...interface{}) {
	if l.fields != nil {
		l.logger.WithFields(l.fields).Printf(format, args...)
	} else {
		l.logger.Printf(format, args...)
	}
}

func (l *logrusLogger) Debugln(args ...interface{}) {
	if l.fields != nil {
		l.logger.WithFields(l.fields).Debugln(args...)
	} else {
		l.logger.Debugln(args...)
	}
}

func (l *logrusLogger) Traceln(args ...interface{}) {
	if l.fields != nil {
		l.logger.WithFields(l.fields).Traceln(args...)
	} else {
		l.logger.Traceln(args...)
	}
}

func (l *logrusLogger) Infoln(args ...interface{}) {
	if l.fields != nil {
		l.logger.WithFields(l.fields).Infoln(args...)
	} else {
		l.logger.Infoln(args...)
	}
}

func (l *logrusLogger) Warnln(args ...interface{}) {
	if l.fields != nil {
		l.logger.WithFields(l.fields).Warnln(args...)
	} else {
		l.logger.Warnln(args...)
	}
}

func (l *logrusLogger) Errorln(args ...interface{}) {
	if l.fields != nil {
		l.logger.WithFields(l.fields).Errorln(args...)
	} else {
		l.logger.Errorln(args...)
	}
}

func (l *logrusLogger) Fatalln(args ...interface{}) {
	if l.fields != nil {
		l.logger.WithFields(l.fields).Fatalln(args...)
	} else {
		l.logger.Fatalln(args...)
	}
}

func (l *logrusLogger) Panicln(args ...interface{}) {
	if l.fields != nil {
		l.logger.WithFields(l.fields).Panicln(args...)
	} else {
		l.logger.Panicln(args...)
	}
}

func (l *logrusLogger) Println(args ...interface{}) {
	if l.fields != nil {
		l.logger.WithFields(l.fields).Println(args...)
	} else {
		l.logger.Println(args...)
	}
}

func (l *logrusLogger) WithFields(fields map[string]interface{}) AppLogger {
	return &logrusLogger{
		logger: l.logger,
		fields: fields,
	}
}
