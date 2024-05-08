package main

import (
  "time"

	"DockerRight/pkg/config"
	"DockerRight/pkg/docker"
  "DockerRight/pkg/log"
  "DockerRight/pkg/notify"
)

func init() {
  log.Info("Initializing DockerRight")
  log.TempInit()
  config.Init("./config.json")
  log.Init(config.Conf.LogLevel)
  docker.Init()
  notify.Init()
}

func main() {
  log.Info("Starting DockerRight")

  if config.Conf.BackupOnStartup {
    log.Info("Running DockerRight on startup")
    docker.BackupContainers()
  }

  go monitorLoop()

  lastBackupHour := -1
  for true {
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
