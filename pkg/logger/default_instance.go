package logger

var instance AppLogger

func GetDefaultLogger() AppLogger {
	return instance
}

func init() {
	instance = NewLogrusLogger(nil)
}

func SetContext(ctx string) {
	switch l := instance.(type) {
	case *logrusLogger:
		l.entry = l.entry.WithField("context", ctx)
	default:
		// do nothing
	}
}

func Config(config *AppLoggerConfig) {
	instance.Config(config)
}

func ConfigAsCommandLineTool() {
	instance.Config(&AppLoggerConfig{
		Level:      LoggerLevelInfo,
		FormatType: LoggerFormatCliTool,
		Output:     LoggerOutputStdout,
	})
}

func Trace(args ...interface{}) { instance.Trace(args...) }
func Debug(args ...interface{}) { instance.Debug(args...) }
func Info(args ...interface{})  { instance.Info(args...) }
func Warn(args ...interface{})  { instance.Warn(args...) }
func Error(args ...interface{}) { instance.Error(args...) }
func Fatal(args ...interface{}) { instance.Fatal(args...) }
func Panic(args ...interface{}) { instance.Panic(args...) }
func Print(args ...interface{}) { instance.Print(args...) }

func Traceln(args ...interface{}) { instance.Traceln(args...) }
func Debugln(args ...interface{}) { instance.Debugln(args...) }
func Infoln(args ...interface{})  { instance.Infoln(args...) }
func Warnln(args ...interface{})  { instance.Warnln(args...) }
func Errorln(args ...interface{}) { instance.Errorln(args...) }
func Fatalln(args ...interface{}) { instance.Fatalln(args...) }
func Panicln(args ...interface{}) { instance.Panicln(args...) }
func Println(args ...interface{}) { instance.Println(args...) }

func Tracef(format string, args ...interface{}) { instance.Tracef(format, args...) }
func Debugf(format string, args ...interface{}) { instance.Debugf(format, args...) }
func Infof(format string, args ...interface{})  { instance.Infof(format, args...) }
func Warnf(format string, args ...interface{})  { instance.Warnf(format, args...) }
func Errorf(format string, args ...interface{}) { instance.Errorf(format, args...) }
func Fatalf(format string, args ...interface{}) { instance.Fatalf(format, args...) }
func Panicf(format string, args ...interface{}) { instance.Panicf(format, args...) }
func Printf(format string, args ...interface{}) { instance.Printf(format, args...) }

func WithFields(fields map[string]interface{}) AppLogger { return instance.WithFields(fields) }
