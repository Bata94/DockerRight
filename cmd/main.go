package main

import (
	"os"
	"time"

	"github.com/bata94/DockerRight/internal/config"
	"github.com/bata94/DockerRight/internal/docker"
	"github.com/bata94/DockerRight/internal/log"
	"github.com/bata94/DockerRight/internal/notify"
)

func init() {
	log.Info("Initializing DockerRight")
	log.TempInit()
	os.Mkdir("./config", 0o755)
	config.Init("./config/config.json")
	log.Init(config.Conf.LogLevel)
	docker.Init()
	notify.Init()
}

func main() {
	log.Info("Starting DockerRight")

	if config.Conf.EnableBackup == false && config.Conf.EnableMonitor == false {
		log.Warn("DockerRight is disabled! Edit the config file and restart :)")
		return
	}

	if config.Conf.BackupOnStartup && config.Conf.EnableBackup {
		log.Info("Running DockerRight on startup")
		docker.BackupContainers()
	}

	if config.Conf.EnableMonitor {
		// TODO: Not working if Backup is disabled
		go monitorLoop()
	}

	lastBackupHour := -1
	for config.Conf.EnableBackup {
		curHour := time.Now().Hour()
		log.Info("Current hour: ", curHour)
		log.Info("BackupHours: ", config.Conf.BackupHours)
		for _, hour := range config.Conf.BackupHours {
			if hour == curHour && lastBackupHour != hour {
				log.Info("Running backup at hour: ", hour)
				err := docker.BackupContainers()
				if err != nil {
					log.Error(err)
				} else {
					lastBackupHour = hour
				}
			}
		}
		minutes2FullHour := 60 - time.Now().Minute()
		log.Info("minutes2FullHour: ", minutes2FullHour)
		if minutes2FullHour < 0 {
			log.Warn("minutes2FullHour < 0, setting minutes2FullHour to 2")
			minutes2FullHour = 2
		} else if minutes2FullHour >= 10 {
			log.Info("minutes2FullHour >= 10, setting minutes2FullHour to 10")
			minutes2FullHour = 10
		}
		sleepDur := time.Duration(minutes2FullHour) * time.Second * 60
		log.Info("Sleeping for ", sleepDur, "...")
		time.Sleep(sleepDur)
	}
}

func monitorLoop() {
	sleepDur := time.Duration(config.Conf.MonitorIntervalSeconds) * time.Second
	for true {
		docker.MonitorContainers()
		log.Info("Sleeping for ", sleepDur, "...")
		time.Sleep(sleepDur)
	}
}
