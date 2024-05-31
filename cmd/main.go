package main

import (
	"os"
	"sync"
	"time"

	"github.com/bata94/DockerRight/internal/config"
	"github.com/bata94/DockerRight/internal/docker"
	"github.com/bata94/DockerRight/internal/log"
	"github.com/bata94/DockerRight/internal/notify"
)

func init() {
	log.Info("Initializing DockerRight")
	log.TempInit()

	err := os.MkdirAll("./config", 0o755)
	if err != nil {
		log.Fatal("Error creating config directory: ", err)
	}
	config.Init("./config/config.json")
	log.Init(config.Conf.LogLevel, config.Conf.LogsPath)
	docker.Init()
	notify.Init(config.Conf.NotifyLevel, config.Conf.TelegramBotToken, config.Conf.TelegramChatIDs)
}

func main() {
	log.Info("Starting DockerRight")
	lastBackup := ""
	mainWg := sync.WaitGroup{}

	if !config.Conf.EnableBackup && !config.Conf.EnableMonitor {
		log.Warn("DockerRight is disabled! Edit the config file and restart :)")
		return
	}

	go func() {
		log.Debug("Started LogFile Rotation GoRoutine")
		lastDate := time.Now()

		for {
			if lastDate.Day() != time.Now().Day() {
				log.Info("Wakey it's a new Day, new LogFile :)")
				log.SetLoggerFile(config.Conf.LogsPath)
				log.DeleteOldLogs(config.Conf.LogRetentionDays, config.Conf.LogsPath)
				lastDate = time.Now()
			}
			time.Sleep(time.Duration(60 - time.Now().Minute()))
		}
	}()

	// TODO: Move the functionality to the lower block
	if config.Conf.BackupOnStartup && config.Conf.EnableBackup {
		log.Info("Running DockerRight on startup")
		err := docker.BackupContainers()
		if err != nil {
			if config.Conf.EnableMonitor {
				log.Error(err)
			} else {
				log.Fatal(err)
			}
		}
		lastBackup = time.Now().Format("2006-01-02T15")
	}

	if config.Conf.EnableMonitor {
		mainWg.Add(1)
		go func() {
			defer mainWg.Done()
			monitorLoop()
		}()
	}

	for config.Conf.EnableBackup {
		curHour := time.Now().Hour()
		curBackup := time.Now().Format("2006-01-02T15")

		log.Debug("Current hour: ", curHour)
		log.Debug("Current Date: ", curBackup)
		log.Debug("Last Backup Date: ", lastBackup)
		log.Debug("BackupHours: ", config.Conf.BackupHours)

		for _, hour := range config.Conf.BackupHours {
			if hour == curHour && lastBackup != curBackup {
				log.Debug("Running backup at hour: ", hour)
				err := docker.BackupContainers()
				if err != nil {
					// TODO: Sent Error to notify api
					log.Error(err)
				} else {
					lastBackup = curBackup
				}
			}
		}
		minutes2FullHour := 60 - time.Now().Minute()
		log.Debug("minutes2FullHour: ", minutes2FullHour)
		if minutes2FullHour < 0 {
			log.Warn("minutes2FullHour < 0, setting minutes2FullHour to 2... This should not been happening, please report this as a bug!")
			minutes2FullHour = 2
		} else if minutes2FullHour >= 10 {
			log.Debug("minutes2FullHour >= 10, setting minutes2FullHour to 10")
			minutes2FullHour = 10
		}
		sleepDur := time.Duration(minutes2FullHour) * time.Second * 60
		log.Debug("Sleeping for ", sleepDur, "...")
		time.Sleep(sleepDur)
	}

	if !config.Conf.EnableBackup {
		log.Warn("Backup functionality is disabled!")
	}

	mainWg.Wait()
}

func monitorLoop() {
	sleepDur := time.Duration(config.Conf.MonitorIntervalSeconds) * time.Second
	for {
		err := docker.MonitorContainers()
		if err != nil {
			log.Error(err)
			// TODO: Handle error, notify user of crash
		}
		log.Info("Sleeping for ", sleepDur, "...")
		time.Sleep(sleepDur)
	}
}
