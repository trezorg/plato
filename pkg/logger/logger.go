package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

type LogOption func(*logrus.Logger)

func WithLogLevel(level string) LogOption {
	return func(log *logrus.Logger) {
		l, err := logrus.ParseLevel(level)
		if err != nil {
			logrus.Fatalf("Cannot parse log level: %s", level)
		}
		log.SetLevel(l)
	}
}

func WithOutput(writers ...io.Writer) LogOption {
	return func(log *logrus.Logger) {
		if len(writers) == 0 {
			writers = append(writers, os.Stdout)
		}
		multiWriter := io.MultiWriter(writers...)
		log.SetOutput(multiWriter)
	}
}

func WithReportCaller(caller bool) LogOption {
	return func(log *logrus.Logger) {
		log.SetReportCaller(caller)
	}
}

func initLogger(options ...LogOption) {

	const (
		defaultLogLevel       = "INFO"
		defaultLogPrettyPrint = true
		defaultReportCaller   = false
	)

	var (
		defaultFormatter = &logrus.TextFormatter{
			PadLevelText: true,
			ForceColors:  false,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "time",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
			TimestampFormat: "2006-01-02 15:04:05 -0700",
		}
	)

	level, err := logrus.ParseLevel(defaultLogLevel)
	if err != nil {
		logrus.Fatalf("Cannot parse log level: %s", defaultLogLevel)
	}

	log = logrus.New()
	log.SetFormatter(defaultFormatter)
	log.SetReportCaller(defaultReportCaller)
	log.SetOutput(os.Stdout)
	log.SetLevel(level)

	for _, option := range options {
		option(log)
	}

}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Init(options ...LogOption) *logrus.Logger {
	initLogger(options...)
	return log
}

func init() {
	initLogger()
}
