package main

import (
	"fmt"
	"os"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/bata94/DockerRight/internal/config"
	"github.com/bata94/DockerRight/internal/docker"
	"github.com/bata94/DockerRight/internal/log"
	"github.com/bata94/DockerRight/internal/notify"
)

func init() {
	log.Info("Initializing DockerRight")
	log.TempInit()

	err := os.MkdirAll("./config", 0o644)
	if err != nil {
		log.Fatal("Error creating config directory: ", err)
	}
	config.Init("./config/config.json")
	log.Init(config.Conf.LogLevel, config.Conf.LogsPath, config.Conf.Log2File)
	docker.Init()
	notify.Init(config.Conf.NotifyLevel, config.Conf.TelegramBotToken, config.Conf.TelegramChatIDs)
}

func main() {
	log.Info("Starting DockerRight")

	if !config.Conf.EnableBackup && !config.Conf.EnableMonitor {
		log.Warn("DockerRight is disabled! Edit the config file and restart :)")
		return
	} else if !config.Conf.EnableBackup {
		log.Warn("Backup functionality is disabled!")
	} else if !config.Conf.EnableMonitor {
		log.Warn("Monitor functionality is disabled!")
	}

	c := cron.New()

	log.Info("Timezone: ", c.Location().String())

	if config.Conf.Log2File {
		log.Debug("Started LogFile Rotation GoRoutine")
		_, err := c.AddFunc("@midnight", func() {
			log.Info("Wakey it's a new Day -> new LogFile :)")
			log.SetLoggerFile(config.Conf.LogsPath)
			log.DeleteOldLogs(config.Conf.LogRetentionDays, config.Conf.LogsPath)
		})
		if err != nil {
			log.Panic("Error adding LogFile rotation cronjob: ", err)
		}
	}

	if config.Conf.EnableMonitor {
		go monitorLoop(config.Conf.MonitorIntervalSeconds, config.Conf.MonitorRetries)
	}

	if config.Conf.EnableBackup {
		lastBackup := ""
		if config.Conf.BackupOnStartup {
			log.Info("Running DockerRight on startup")
			err := docker.BackupContainers()
			if err != nil {
				log.Error(err)
			} else {
				lastBackup = time.Now().Format("2006-01-02T15")
			}
		}

		for _, hour := range config.Conf.BackupHours {
			cSchedule := fmt.Sprintf("5 %v * * *", hour)
			_, err := c.AddFunc(cSchedule, func() {
				curBackup := time.Now().Format("2006-01-02T15")
				if curBackup != lastBackup {
					log.Debug("Running backup at hour: ", hour)
					err := docker.BackupContainers()
					if err != nil {
						log.Error(err)
					}
				} else {
					log.Warn("Backup already ran at hour: ", hour, "\n", "This should only happen on startup and if you are running a backup on startup!")
				}
			})
			if err != nil {
				log.Panic("Error adding backup cronjob: ", err)
			}
		}
	}

	c.Start()
	log.Info("Number of current registered Cronjobs, it should show the daily LogFile rotation job as well as the Backups (as configured): ", len(c.Entries()))

	select {}
}

func monitorLoop(intervalSec, monitorRetries int) {
	containerInfos := []docker.ContainerInfo{}
	ticker := time.NewTicker(time.Duration(intervalSec) * time.Second)
	for {
		err := docker.MonitorContainers(&containerInfos)
		if err != nil {
			log.MonitorMsg(err)
		}

		for i, ci := range containerInfos {
			if len(ci.States) < monitorRetries {
				containerInfos[i].MonitorState = "unknown"
				break
			}

			isRunning := false
			isRunningCount := 0

			for _, s := range ci.States[len(ci.States)-monitorRetries:] {
				if s == "running" || s == "healthy" {
					isRunningCount++
				}
			}

			if isRunningCount >= monitorRetries {
				isRunning = true
			}

			if isRunning {
				if ci.MonitorState == "unknown" || ci.MonitorState == "" {
					containerInfos[i].MonitorState = "running"
				} else if ci.MonitorState == "stopped" || ci.MonitorState == "unhealthy" || ci.MonitorState == "exited" {
					containerInfos[i].MonitorState = "running"
					log.MonitorMsg(ci.Name, " is UP and running again :)")
				}
			} else if isRunningCount == 0 {
				if ci.MonitorState != "stopped" && ci.MonitorState != "unhealthy" && ci.MonitorState != "exited" {
					curState := ci.States[len(ci.States)-1]
					containerInfos[i].MonitorState = curState
					log.MonitorMsg(ci.Name, " is ", curState, "!")
				} else {
					log.Debug("Container stopped but not changed!")
				}
			}

			if len(ci.States) >= monitorRetries*4 {
				containerInfos[i].States = ci.States[monitorRetries*2:]
			}
		}

		log.Info("Sleeping for ", intervalSec, "...")
		<-ticker.C
	}
}
