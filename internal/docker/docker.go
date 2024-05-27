package docker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/bata94/DockerRight/internal/config"
	"github.com/bata94/DockerRight/internal/log"
	"github.com/bata94/DockerRight/internal/workpool"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

var (
	ctx      = context.Background()
	cli      *client.Client
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
	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		log.Fatal(err)
	}

	if config.Conf.CreateTestContainerOnStartup {
		log.Info("Creating test container")
		_, err := RunContainer(RunContainerParams{
			ContainerName: "TestContainer",
			ImageName:     defImage,
			Cmd:           []string{"echo", "Running inside TestContainer"},
			Remove:        true,
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
		if err != nil {
			return errors.New("Error pulling image: " + err.Error())
		}
		defer reader.Close()
		_, err = io.Copy(os.Stdout, reader)
		if err != nil {
			return errors.New("Error pulling image: " + err.Error())
		}
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
	ImageName     string
	Cmd           []string
	Remove        bool
	Volumes       map[string]struct{}
	VolumesFrom   []string
	Mounts        []mount.Mount
}

func RunContainer(p RunContainerParams) ([]byte, error) {
	log.Debug("Running Container")

	err := PullImage(p.ImageName)
	if err != nil {
		log.Fatal(err)
	}

	ctr, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:           p.ImageName,
			Tty:             false,
			Cmd:             p.Cmd,
			NetworkDisabled: true,
			Volumes:         p.Volumes,
		},
		&container.HostConfig{
			VolumesFrom: p.VolumesFrom,
			Mounts:      p.Mounts,
		},
		&network.NetworkingConfig{},
		nil,
		p.ContainerName,
	)
	if err != nil {
		log.Error("Error creating container: ")
		log.Error(err)
		_ = RemoveContainer(ctr.ID)
		return nil, err
	}

	err = cli.ContainerStart(ctx, ctr.ID, container.StartOptions{})
	if err != nil {
		log.Error("Error starting container: ")
		log.Error(err)
		_ = RemoveContainer(ctr.ID)
		return nil, err
	}

	statusCh, errCh := cli.ContainerWait(ctx, ctr.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			log.Error(err)
			log.Error(err)
			_ = RemoveContainer(ctr.ID)
			return nil, err
		}
	case <-statusCh:
		log.Debug("Container finished!")
	}

	out, err := cli.ContainerLogs(ctx, ctr.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		log.Error("Error getting logs: ", err)
	}
	defer out.Close()
	if err != nil {
		log.Error(err)
		_ = RemoveContainer(ctr.ID)
		return nil, err
	}

	logs, err := io.ReadAll(out)
	if err != nil {
		log.Error(err)
		_ = RemoveContainer(ctr.ID)
		return nil, err
	}

	log.Debug("Container output:", "\n", string(logs))

	if p.Remove {
		err = RemoveContainer(ctr.ID)
		if err != nil {
			log.Fatal(err)
		}
	}

	return logs, nil
}

func MonitorContainers() error {
	log.Info("MonitorContainers")

	return nil
}

func GetHostBackupPath(containers []types.Container) string {
	hostBackupPath := ""

	for _, ctr := range containers {
		log.Debug(ctr.Names)
		if strings.Contains(strings.ToLower(ctr.Names[0]), "dockerright") {
			for _, m := range ctr.Mounts {
				log.Debug(m)
				log.Debug("BackupPathConf: " + config.Conf.BackupPath)
				log.Debug("MountDestination: " + m.Destination)
				if strings.EqualFold(m.Destination, config.Conf.BackupPath) {
					hostBackupPath = m.Source
				}
			}
		}
	}

	return hostBackupPath
}

func RunOSCmd(cmdType, cmd string) ([]byte, error) {
	if cmd != "" {
		runCmd := exec.Command("sh", "-c", cmd)

		output, err := runCmd.Output()
		if err != nil {
			log.Error("Error running BeforeBackupCMD: ")
			log.Error(err)
			return nil, err
		}

		logPath := config.Conf.LogsPath
		if !strings.HasSuffix(logPath, "/") {
			logPath = logPath + "/"
		}
		err = os.WriteFile(logPath+cmdType+"-"+time.Now().Format("2006-01-02-15:04:05")+".log", output, 0644)
		if err != nil {
			log.Error("Error writing log: ", err)
		}

		return output, nil
	} else {
		return nil, nil
	}
}

func BackupContainers() error {
	log.Info("BackupContainers")

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		log.Error("Error listing containers: ")
		log.Error(err)
		return err
	}

	hostBackupPath := GetHostBackupPath(containers)
	if hostBackupPath == "" {
		err := errors.New("Error finding backup path! Maybe you changed the Dockerright Container name to something else then 'dockerright'?")
		log.Error(err)
		return err
	}

	log.Info("Running BeforeBackupCMD", "\n", config.Conf.BeforeBackupCMD)
	output, err := RunOSCmd("BeforeBackupCMD", config.Conf.BeforeBackupCMD)
	if err != nil {
		log.Error("Error running BeforeBackupCMD: ", err)
	} else {
		log.Info("BeforeBackupCMD ran successfully, Output:", "\n", string(output))
	}

	var wg workpool.WaitGroupCount
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

		wg.Add(1)
		log.Info("Current concurrent BackupRunners: ", wg.GetCount())
		for wg.GetCount() > config.Conf.ConcurrentBackupContainer {
			time.Sleep(time.Duration(time.Millisecond * 250))
		}
		go func(ctr types.Container) {
			defer wg.Done()
			backupErr := RunBackupHelperForContainer(ctr, hostBackupPath)

			if backupErr != nil {
				log.Error("Error in concurrent backup runner ", ctr.Names[0], " Error: ", backupErr)
			}
		}(ctr)
	}
	wg.Wait()

	log.Info("Running AfterBackupCMD", "\n", config.Conf.AfterBackupCMD)
	output, err = RunOSCmd("AfterBackupCMD", config.Conf.AfterBackupCMD)
	if err != nil {
		log.Error("Error running AfterBackupCMD: ", err)
	} else {
		log.Info("AfterBackupCMD ran successfully, Output:", "\n", string(output))
	}

	log.Info("BackupContainers done")
	err = DeleteOldBackups()
	if err != nil {
		return err
	}

	return nil
}

func RunBackupHelperForContainer(container types.Container, hostBackupPath string) error {
	log.Info("RunBackupHelperForContainer" + container.Names[0])
	log.Info(fmt.Sprintf("%s %s %s (status: %s)\n", container.ID, container.Names, container.Image, container.Status))

	if container.Mounts == nil || len(container.Mounts) == 0 {
		log.Info("Container has no mounts")
		return nil
	}

	containerNameBase := "DockerRight-BackupRunner-" + strings.ReplaceAll(container.Names[0], "/", "")
	now := time.Now()

	backupPathBase := config.Conf.BackupPath
	if !strings.HasSuffix(backupPathBase, "/") {
		backupPathBase = backupPathBase + "/"
	}
	backupPath := strings.ReplaceAll(container.Names[0], "/", "") + "/" + now.Format("2006-01-02-15-04-05") + "/"
	err := os.MkdirAll(backupPathBase+"/"+backupPath, 0o755)
	if err != nil {
		log.Error(err)
		return err
	}

	err = os.WriteFile(backupPathBase+"/"+backupPath+"/ContainerInfo.txt", []byte(log.FormatStruct(container)), 0o644)
	if err != nil {
		log.Error("Unable to save ContainerInfoFile for container ", container.Names[0], " Error: ", err)
	}

	for i, m := range container.Mounts {
		containerName := fmt.Sprint(containerNameBase, "-m", i, "-", strings.ReplaceAll(m.Destination, "/", "_"))
		log.Info(fmt.Sprintf("Creating container %s", containerName))

		// TODO: Move those to a Parameter
		if strings.HasSuffix(m.Destination, ".sock") || strings.HasSuffix(m.Source, ".sock") {
			log.Warn(fmt.Sprintf("Skipping mount %s : %s for Container %s because it contains a socket!", m.Source, m.Destination, containerName))
			continue
		} else if m.Source == "/" {
			log.Warn(fmt.Sprintf("Skipping mount %s : %s for Container %s because it is the root directory!", m.Source, m.Destination, containerName))
			continue
		} else if strings.Contains(m.Destination, "/var/lib/docker/volumes") {
			log.Warn(fmt.Sprintf("Skipping mount %s : %s for Container %s because it is /var/lib/docker/volumes!", m.Source, m.Destination, containerName))
			continue
		}

		mountInfoFileName := fmt.Sprint(m.Type) + strings.Replace(m.Destination, "/", "_", -1)
		// backupPathBase := strings.ReplaceAll(container.Names[0], "/", "") + "/" + now.Format("2006-01-02-15-04-05") + "/" + mountInfoFileName

		cmd := []string{"tar", "cvf", backupPathBase + "/" + backupPath + mountInfoFileName + ".tar", m.Destination}
		log.Debug(cmd)
		out, err := RunContainer(RunContainerParams{
			ContainerName: containerName,
			ImageName:     defImage,
			Cmd:           cmd,
			Remove:        true,
			VolumesFrom:   []string{container.ID},
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: hostBackupPath,
					Target: backupPathBase,
				},
			},
		})
		if err != nil {
			return err
		}

		err = os.WriteFile(backupPathBase+"/"+backupPath+mountInfoFileName+".log", out, 0o644)
		if err != nil {
			log.Error("Unable to save backup logfile for container ", containerName, " Error: ", err)
		}
	}
	time.Sleep(time.Second * 5)

	return nil
}

func DeleteOldBackups() error {
	log.Info("DeleteOldBackups")

	containerDirs, err := os.ReadDir(config.Conf.BackupPath)
	if err != nil {
		log.Error("Error reading backup path: ", err)
		return err
	}

	log.Info("Found ", len(containerDirs), " containerDirs:", "\n", containerDirs)

	for _, c := range containerDirs {
		if c.IsDir() {
			log.Info("ContainerDir: ", c.Name())
			backupDirs, err := os.ReadDir(config.Conf.BackupPath + "/" + c.Name())
			if err != nil {
				log.Error("Error reading backup path: ", err)
				continue
			}
			log.Info("Found ", len(backupDirs), " backupDirs:", "\n", backupDirs)
			for _, b := range backupDirs {
				if b.IsDir() {
					log.Info("BackupDir: ", b.Name())
					backupTime, err := time.Parse("2006-01-02-15-04-05", b.Name())
					if err != nil {
						log.Error("Error parsing backup time: ", err)
						continue
					}
					timeSinceBackup := time.Since(backupTime).Hours()
					log.Info("timeSinceBackup: ", timeSinceBackup)
					if timeSinceBackup > float64(config.Conf.RetentionHours) {
						log.Info("Removing ", config.Conf.BackupPath+"/"+c.Name()+"/"+b.Name())
						err = os.RemoveAll(config.Conf.BackupPath + "/" + c.Name() + "/" + b.Name())
						if err != nil {
							log.Error("Error removing backup: ", err)
							continue
						}
					}
				}
			}
		}
	}

	return nil
}
