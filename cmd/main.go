package main

import (
	// tea "github.com/charmbracelet/bubbletea"
	// "time"

	"DockerRight/pkg/config"
	"DockerRight/pkg/docker"
  "DockerRight/pkg/log"
)

func init() {
  log.Init()
  config.Init("./config.json")
  docker.Init()
}

func main() {
  log.Info("Starting DockerRight")
  runLoop := true
  loopCount := 0

  for runLoop {
    log.Info("Running DockerRight")
    docker.BackupContainers()
    // time.Sleep(3 * time.Second)
    loopCount++

    if loopCount >= 1 {
      runLoop = false
    }
  }
}
