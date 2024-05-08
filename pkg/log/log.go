package log

import (
  "strings"
  logger "github.com/sirupsen/logrus"
)

func TempInit() {
  logger.Info("Initializing temp Logger Module")
  logger.SetLevel(0)
}

func Init(logLvlStr string) {
  logger.Info("Initializing final Logger Module")
  logLvl := logger.InfoLevel
  logLvlStr = strings.ToLower(logLvlStr)

  switch logLvlStr {
  case "debug":
    logLvl = logger.DebugLevel
  case "info":
    logLvl = logger.InfoLevel
  case "warn":
    logLvl = logger.WarnLevel
  case "error":
    logLvl = logger.ErrorLevel
  case "fatal":
    logLvl = logger.FatalLevel
  case "panic":
    logLvl = logger.PanicLevel
  default:
    logLvl = logger.InfoLevel
  }

  logger.SetLevel(logLvl)
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

