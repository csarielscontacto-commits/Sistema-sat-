package logger

import (
    "os"
    "time"

    "github.com/sirupsen/logrus"
)

type Logger struct {
    *logrus.Logger
}

func NewLogger(service string) *Logger {
    log := logrus.New()

    log.SetFormatter(&logrus.JSONFormatter{
        TimestampFormat: time.RFC3339Nano,
        FieldMap: logrus.FieldMap{
            logrus.FieldKeyTime:  "timestamp",
            logrus.FieldKeyLevel: "severity",
            logrus.FieldKeyMsg:   "message",
        },
    })

    log.SetOutput(os.Stdout)

    level := os.Getenv("LOG_LEVEL")
    if level == "" {
        level = "info"
    }
    lvl, err := logrus.ParseLevel(level)
    if err != nil {
        lvl = logrus.InfoLevel
    }
    log.SetLevel(lvl)

    log.WithFields(logrus.Fields{
        "service": service,
        "version": "1.0.0",
    })

    return &Logger{Logger: log}
}

func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
    l.WithFields(logrus.Fields{}).Info(msg)
}

func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
    l.WithFields(logrus.Fields{}).Error(msg)
}

func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
    l.WithFields(logrus.Fields{}).Warn(msg)
}

func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
    l.WithFields(logrus.Fields{}).Debug(msg)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
    l.Logger.Fatalf(format, args...)
}