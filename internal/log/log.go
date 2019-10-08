// Package log is a wrapper around logrus
// For more information see https://godoc.org/github.com/sirupsen/logrus
package log

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
}

const (
	PanicLevel = logrus.PanicLevel
	FatalLevel = logrus.FatalLevel
	ErrorLevel = logrus.ErrorLevel
	InfoLevel  = logrus.InfoLevel
	DebugLevel = logrus.DebugLevel
)

func SetOutput(out io.Writer)     { logger.SetOutput(out) }
func SetLevel(level logrus.Level) { logger.SetLevel(level) }

func Debug(args ...interface{})                 { logger.Debug(args...) }
func Debugf(format string, args ...interface{}) { logger.Debugf(format, args...) }

func Info(args ...interface{})                 { logger.Info(args...) }
func Infof(format string, args ...interface{}) { logger.Infof(format, args...) }

func Warn(args ...interface{})                 { logger.Warn(args...) }
func Warnf(format string, args ...interface{}) { logger.Warnf(format, args...) }

func Error(args ...interface{})                 { logger.Error(args...) }
func Errorf(format string, args ...interface{}) { logger.Errorf(format, args...) }

func Fatal(args ...interface{})                 { logger.Fatal(args...) }
func Fatalf(format string, args ...interface{}) { logger.Fatalf(format, args...) }
