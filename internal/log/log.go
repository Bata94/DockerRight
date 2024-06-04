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

	"github.com/bata94/DockerRight/internal/notify"
	logger "github.com/sirupsen/logrus"
)

var loggerLvl int

func FormatListOfStructs(vL ...interface{}) string {
	retStr := "[\n"
	for _, v := range vL {
		retStr += FormatStruct(v) + ",\n"
	}
	retStr += "]"

	return retStr
}

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
	loggerLvl = 3 // logger.WarnLevel
	logger.SetLevel(logger.DebugLevel)
}

func Init(logLvlStr, logPath string, log2File bool) {
	logger.Info("Initializing final Logger Module")
	var logLvl logger.Level

	switch strings.ToLower(logLvlStr) {
	case "debug":
		loggerLvl = 5
		logLvl = logger.DebugLevel
	case "info":
		loggerLvl = 4
		logLvl = logger.InfoLevel
	case "warn":
		loggerLvl = 3
		logLvl = logger.WarnLevel
	case "error":
		loggerLvl = 2
		logLvl = logger.ErrorLevel
	case "fatal":
		loggerLvl = 1
		logLvl = logger.FatalLevel
	case "panic":
		loggerLvl = 1
		logLvl = logger.PanicLevel
	default:
		loggerLvl = 4
		logLvl = logger.InfoLevel
	}

	logger.Info("Setting log level to: ", logLvl, loggerLvl)
	logger.SetLevel(logLvl)

	err := os.MkdirAll(logPath, 0o644)
	if err != nil {
		logger.Panic("Error creating LogPath: ", err)
	}

	if log2File {
		SetLoggerFile(logPath)
		Info("Logging to STDOUT and File")
	} else {
		Info("Logging to STDOUT")
	}
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
		fileNameCleaned := strings.Replace(l.Name(), ".log", "", 1)

		if fileNameCleaned[0] != '2' {
			fileNameCleaned = strings.TrimPrefix(fileNameCleaned, fileNameCleaned[:strings.Index(fileNameCleaned, "-")+1])
		}

		var logDate time.Time

		if len(fileNameCleaned) > 10 {
			logDate, err = time.Parse("2006-01-02-05:04:05", fileNameCleaned)
			if err != nil {
				Error("Error Parsing LogFile Name to Date: ", err)
				continue
			}
		} else {
			logDate, err = time.Parse("2006-01-02", fileNameCleaned)
			if err != nil {
				Error("Error Parsing LogFile Name to Date: ", err)
				continue
			}
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
	notify.Notifier(5, err)
	logger.Debug(err...)
}

func Info(err ...interface{}) {
	notify.Notifier(4, err)
	logger.Info(err...)
}

func Warn(err ...interface{}) {
	notify.Notifier(3, err)
	logger.Warn(err...)
}

func MonitorMsg(err ...interface{}) {
	notify.Notifier(-1, err)
	logger.Warn(err...)
}

func Error(err ...interface{}) {
	notify.Notifier(2, err)
	logger.Error(err...)
}

func Fatal(err ...interface{}) {
	notify.Notifier(1, err)
	logger.Fatal(err...)
}

func Panic(err ...interface{}) {
	notify.Notifier(1, err)
	logger.Panic(err...)
}
