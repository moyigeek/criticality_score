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
	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
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
	LoggerFormatConsole LoggerFormatType = iota
	LoggerFormatMessageOnly
	LoggerFormatJSON
)

type logrusLogger struct {
	logger *logrus.Logger
	fields map[string]interface{}
}

func NewLogrusLogger(config *AppLoggerConfig) AppLogger {
	if config == nil {
		config = &AppLoggerConfig{
			Level:      string(LoggerLevelInfo),
			FormatType: LoggerFormatConsole,
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
	case LoggerFormatConsole:
		formatter = &logrus.TextFormatter{}
	case LoggerFormatMessageOnly:
		formatter = &logrus.TextFormatter{
			DisableTimestamp: true,
		}
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

func (l *logrusLogger) WithFields(fields map[string]interface{}) AppLogger {
	return &logrusLogger{
		logger: l.logger,
		fields: fields,
	}
}
