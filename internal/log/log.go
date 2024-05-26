package log

// A Wrapper for the logrus logger
// maybe useful in the future :D

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"time"

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

func Init(logLvlStr, logPath string) {
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

	SetLoggerFile(logPath)
}

func SetLoggerFile(logPath string) {
	if !strings.HasSuffix(logPath, "/") {
		logPath = logPath + "/"
	}
	logFilePath := logPath + time.Now().Format("2006-01-02") + ".log"
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o644)
	if err != nil {
		logger.Panic("Can't write to log File! Panicing....")
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	logger.SetOutput(mw)

	Info("Logger File Output Set!")
}

func DeleteOldLogs(retentionDays int, logPath string) {
	Debug("Deleting old Logs...")

	if !strings.HasSuffix(logPath, "/") {
		logPath = logPath + "/"
	}
	logFiles, err := os.ReadDir(logPath)
	if err != nil {
		Error("Error reading LogPath path: ", err)
		return
	}

	for _, l := range logFiles {
		logDate, err := time.Parse("2006-01-02", strings.Replace(l.Name(), ".log", "", 1))
		if err != nil {
			Error("Error Parsing LogFile Name to Date: ", err)
			continue
		}

		if (time.Now().Day() - logDate.Day()) > retentionDays {
			err := os.Remove(logPath + l.Name())
			if err != nil {
				Error("Error removeing old Log: ", l.Name(), err)
				continue
			}
		}
	}

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
