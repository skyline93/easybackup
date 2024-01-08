package log

import (
	slog "log/syslog"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/syslog"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.SetFormatter(&logrus.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.DebugLevel)

	hook, err := syslog.NewSyslogHook("", "", slog.LOG_INFO, "")
	if err == nil {
		log.Hooks.Add(hook)
	}
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}
