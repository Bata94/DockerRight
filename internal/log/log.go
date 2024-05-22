package log

// A Wrapper for the logrus logger
// maybe useful in the future :D

import (
	"fmt"
	"reflect"
	"strings"

	logger "github.com/sirupsen/logrus"
)

func FormatStruct(v interface{}) string {
	val := reflect.ValueOf(v)
	typ := val.Type()

	if typ.Kind() != reflect.Struct {
		return "Provided value is not a struct"
	}

	retStr := "{\n"
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		retStr += fmt.Sprintf("  %s: %v, \n", field.Name, value)
	}
	retStr += "}"

	return retStr
}

func TempInit() {
	logger.Info("Initializing temp Logger Module")
	logger.SetLevel(logger.DebugLevel)
}

func Init(logLvlStr string) {
	logger.Info("Initializing final Logger Module")
	var logLvl logger.Level

	switch strings.ToLower(logLvlStr) {
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

	logger.Info("Setting log level to: ", logLvl)
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
