package docker

import (
	"DockerRight/pkg/config"
	"DockerRight/pkg/log"

  "time"
	"context"
	"fmt"
	"io"
	"os"
  "os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	// "github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

var (
  ctx = context.Background()
  cli *client.Client
  defImage = "debian:latest"
)

func Init() {
  log.Info("Initializing Docker Module")
  var err error

  cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
  if err != nil {
    log.Error("Error initializing Docker client: ")
    log.Fatal(err)
  }
  defer cli.Close()

  reader, err := cli.ImagePull(ctx, defImage, image.PullOptions{})
  if err != nil {
    log.Error("Error pulling image: ")
    log.Fatal(err)
  }
  defer reader.Close()
  io.Copy(os.Stdout, reader)

  if config.Conf.CreateTestContainerOnStartup {
    log.Info("Creating test container")
    err := RunContainer(RunContainerParams{
      ContainerName: "TestContainer",
      ImageName: defImage,
      Cmd: []string{"echo", "Running inside TestContainer"},
      Remove: true,
    })
    if err != nil {
      log.Error("Error creating test container: ")
      log.Fatal(err)
    }
  }

  log.Info("Docker initialized")
}

func PullImage(imageName string) error {
  log.Debug("Check if Image needs to be pulled")
  images, err := cli.ImageList(ctx, image.ListOptions{})
  if err != nil {
    log.Error("Error listing images: ")
    log.Error(err)
    return err
  }
  pull := true
  for _, img := range images {
    for _, imgTag := range img.RepoTags {
      if imageName == imgTag {
        log.Debug("Image already pulled")
        pull = false
        break
      }
    }
    if !pull {
      break
    }
  }

  if pull {
    log.Debug("Pulling Image")
    reader, err := cli.ImagePull(ctx, imageName, image.PullOptions{})
    defer reader.Close()
    if err != nil {
      log.Error("Error pulling image: ")
      return err
    }
    io.Copy(os.Stdout, reader)
  }

  return nil
}

func RemoveContainer(containerId string) error {
  log.Debug("Removing Container")
  err := cli.ContainerRemove(ctx, containerId, container.RemoveOptions{
    Force: true,
  })
  if err != nil {
    log.Error("Error removing container: ")
    return err
  }
  return nil
}

type RunContainerParams struct {
  ContainerName string
  ImageName string
  Cmd []string
  Remove bool
  Volumes map[string]struct{}
  VolumesFrom []string
  Mounts []mount.Mount
}

func RunContainer(p RunContainerParams) error {
  log.Debug("Running Container")

  err := PullImage(p.ImageName)
  if err != nil {
    log.Fatal(err)
  }

  ctr, err := cli.ContainerCreate(
    ctx,
    &container.Config{
      Image: p.ImageName,
      Tty: false,
      Cmd: p.Cmd,
      NetworkDisabled: true,
      Volumes: p.Volumes,
    },
    &container.HostConfig{
      VolumesFrom: p.VolumesFrom,
      Mounts: p.Mounts,
    },
    &network.NetworkingConfig{},
    nil,
    p.ContainerName,
    )
  if err != nil {
    log.Error("Error creating container: ")
    log.Error(err)
    _ = RemoveContainer(ctr.ID)
    return err
  }

  err = cli.ContainerStart(ctx, ctr.ID, container.StartOptions{})
  if err != nil {
    log.Error("Error starting container: ")
    log.Error(err)
    _ = RemoveContainer(ctr.ID)
    return err
  }

  statusCh, errCh := cli.ContainerWait(ctx, ctr.ID, container.WaitConditionNotRunning)
  select {
  case err := <-errCh:
    if err != nil {
      log.Error(err)
      log.Error(err)
      _ = RemoveContainer(ctr.ID)
      return err
    }
  case <-statusCh:
    log.Debug("Container finished!")
}

  out, err := cli.ContainerLogs(ctx, ctr.ID, container.LogsOptions{
    ShowStdout: true,
    ShowStderr: true,
    Follow: true,
  })
  defer out.Close()
  if err != nil {
    log.Error(err)
    _ = RemoveContainer(ctr.ID)
    return err
  }

  logs, err := io.ReadAll(out)
  if err != nil {
    log.Error(err)
    _ = RemoveContainer(ctr.ID)
    return err
  }

  log.Debug("Container output:", "\n", string(logs))

  if p.Remove {
    err = RemoveContainer(ctr.ID)
    if err != nil {
      log.Fatal(err)
    }
  }

  return nil
}

func MonitorContainers() error {
  log.Info("MonitorContainers")

  return nil
}

func BackupContainers() error {
  log.Info("BackupContainers")

  if config.Conf.BeforeBackupCMD != "" {
    log.Info("Running BeforeBackupCMD", "\n", config.Conf.BeforeBackupCMD)
    runCmd := exec.Command("sh", "-c", config.Conf.BeforeBackupCMD)

    output, err := runCmd.Output()
    if err != nil {
      log.Error("Error running BeforeBackupCMD: ")
      log.Error(err)
      return err
    }

    log.Info("BeforeBackupCMD ran successfully, Output:", "\n", string(output))
  }

  containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
  if err != nil {
  	log.Error("Error listing containers: ")
  	log.Error(err)
  	return err
  }

  for _, ctr := range containers {
    // Skip dockerright named containers
    skip := false
    for _, containerName := range ctr.Names {
      if strings.Contains(strings.ToLower(containerName), "dockerright") {
        skip = true
        continue
      }
    }
    if skip {
      continue
    }

    backupErr := RunBackupHelperForContainer(ctr)
    if backupErr != nil {
      continue
    }
  }
  
  if config.Conf.AfterBackupCMD != "" {
    log.Info("Running AfterBackupCMD", "\n", config.Conf.AfterBackupCMD)
    runCmd := exec.Command("sh", "-c", config.Conf.AfterBackupCMD)

    output, err := runCmd.Output()
    if err != nil {
      log.Error("Error running AfterBackupCMD: ")
      log.Error(err)
      return err
    }

    log.Info("AfterBackupCMD ran successfully, Output:", "\n", string(output))
  }

  return nil
}

func RunBackupHelperForContainer(container types.Container) error {
  log.Debug("RunBackupHelperForContainer" + container.Names[0])
  log.Debug(fmt.Sprintf("%s %s %s (status: %s)\n", container.ID, container.Names, container.Image, container.Status))

  if container.Mounts == nil || len(container.Mounts) == 0 {
    log.Debug("Container has no mounts")
    return nil
  }

  containerName := "DockerRight-BackupRunner-" + container.ID
  now := time.Now()
  backupPath := ""
  curPWD, err := os.Getwd()
  if err != nil {
    log.Error(err)
    return err
  }
  backupPath = curPWD + "/backup" + container.Names[0] + "/" + now.Format("2006-01-02-15-04-05")
  err = os.MkdirAll(backupPath, 0755)
  if err != nil {
    log.Error(err)
    return err
  }


  for _, m := range container.Mounts {
    log.Debug(fmt.Sprintf("Creating container %s", containerName))
    mountInfoFileName := fmt.Sprint(m.Type) + strings.Replace(m.Destination, "/", "_", -1)
    cmd := []string{"tar", "cvf", "/opt/backup" + "/" + mountInfoFileName + ".tar", m.Destination}
    log.Debug(cmd)
    err := RunContainer(RunContainerParams{
      ContainerName: containerName,
      ImageName: defImage,
      Cmd: cmd,
      Remove: true,
      VolumesFrom: []string{container.ID},
      Mounts: []mount.Mount{
        {
          Type:   mount.TypeBind,
          Source: backupPath,
          Target: "/opt/backup",
        },
      },
    })
    if err != nil {
      return err
    }
  }

  return nil
}
