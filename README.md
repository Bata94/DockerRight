# DockerRight

* [Description](#description)
* [Limitations/Warnings](#limitations-warnings)
* [How to use](#how-to-use)
* [Configuration](#configuration)
* [TODOs](#todos)
* [License](#license)

## Description

A simple Docker Container that allows you to monitor your Docker Containers and backup your Docker Volumes, with notifictaions.

It creates a tarball (.tar) per configured volume, at the configured time, and stores it in the configured location :)
In the best case the Output directory is mapped to a Network Drive or another Host.

To see whats working now and whats planned in the near future see [TODOs](#todos).

## Limitations/Warnings

Not recommended to use for DBs or complex applications, but still better than no backup at all :D

It is a very young and developing Project, so use at your own risk!

## How to use

Minimal Docker run CMD:
``` bash
docker run -d 
    -v /var/run/docker.sock:/var/run/docker.sock
    -v /path/to/config:/opt/DockerRight/config
    -v /path/to/backupDir:/opt/DockerRight/backup
    -e TZ=Europe/Berlin
    --restart always
    --name dockerright 
    ghcr.io/bata94/dockerright:latest
```

Minimal Docker-Compose File:
``` yaml
services:
  dockerright:
    image: ghcr.io/bata94/dockerright:latest
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - /path/to/config:/opt/DockerRight/config
      - /path/to/backupDir:/opt/DockerRight/backup
    environment:
      TZ: Europe/Berlin
    restart: always
```

Start the Container, it will stop after a few seconds on it's own.

Now you edit the created config.json. 

### Configuration

Parameters that can be set in the config.json. To reset them, delete the config.json and restart the container.

| Parameter                     | Default                    | Type     | Description                                            |
|-------------------------------|----------------------------|----------|--------------------------------------------------------|
| Enabled                       | false                      | Bool     | Enable Service                                         |
| MonitorIntervalSeconds        | 60                         | Int      | Interval in seconds                                    |
| MonitorReties                 | 5                          | Int      | Retries before sending notification                    |
| BackupHours                   | []                         | []Int    | Backup at these hours                                  |
| RetentionHours                | 120                        | Int      | Retention in hours (24 * 5)                            |
| ConcurrentBackupContainer     | numCPUs/2                  | Int      | How many mounts should be backed up at once            |
| BackupPath                    | "/opt/dockerBackup/backup" | String   | Backup Path inside Container (shouldn't be changed)    |
| BeforeBackupCMD               | ""                         | String   | CMD to execute before backup                           |
| AfterBackupCMD                | ""                         | String   | CMD to execute after backup                            |
| LogLevel                      | "info"                     | String   | Log Level                                              |
| BackupOnStartup               | false                      | Bool     | Backup on Startup                                      |
| CreateTestContainerOnStartup  | true                       | Bool     | Create Test Container on Startup, to check docker.sock |

## TODOs

What's planned in the near future? If you have any ideas, please [open an issue](https://github.com/bata94/dockerRight/issues) or [create a PR](https://github.com/bata94/dockerRight/pulls) :)

- [X] Create Backups per mount
- [ ] Telegram Notifications
- [ ] Discord Notifications
- [ ] Mail Notifications
- [ ] Logs to File
- [ ] BackupContainer Output to File
- [ ] Make config parameters settable by environment variables
- [ ] Fine grain settings via Container Labels (like traefik for example)
- [ ] Add tests

## License

This project is licensed under the [Unlicense](https://unlicense.org/), so do what you want, but don't blame me :D 
