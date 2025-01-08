package logger

var instance AppLogger

func GetDefaultLogger() AppLogger {
	return instance
}

func init() {
	instance = NewLogrusLogger(nil)
}

func Config(config *AppLoggerConfig) {
	instance.Config(config)
}

func Trace(args ...interface{}) {
	instance.Trace(args...)
}

func Debug(args ...interface{}) {
	instance.Debug(args...)
}

func Info(args ...interface{}) {
	instance.Info(args...)
}

func Warn(args ...interface{}) {
	instance.Warn(args...)
}

func Error(args ...interface{}) {
	instance.Error(args...)
}

func Fatal(args ...interface{}) {
	instance.Fatal(args...)
}

func Panic(args ...interface{}) {
	instance.Panic(args...)
}

func Tracef(format string, args ...interface{}) {
	instance.Tracef(format, args...)
}

func Debugf(format string, args ...interface{}) {
	instance.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	instance.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	instance.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	instance.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	instance.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	instance.Panicf(format, args...)
}

func WithFields(fields map[string]interface{}) AppLogger {
	return instance.WithFields(fields)
}
