package log

import (
  logger "github.com/sirupsen/logrus"
)

func Init() {
  logger.SetLevel(logger.DebugLevel)
}

func Debug(err ...interface{}) {
  logger.Debug(err...)
}

func Info(err ...interface{}) {
  logger.Info(err...)
}

func Warn(err ...interface{}) {
  logger.Warn(err...)
}

func Error(err ...interface{}) {
  logger.Error(err...)
}

func Fatal(err ...interface{}) {
  logger.Fatal(err...)
}

func Panic(err ...interface{}) {
  logger.Panic(err...)
}

