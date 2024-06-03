# DockerRight

[![Go Report Card](https://goreportcard.com/badge/github.com/bata94/DockerRight)](https://goreportcard.com/report/github.com/bata94/DockerRight)

* [Description](#description)
* [Limitations/Warnings](#limitations-warnings)
* [How to use](#how-to-use)
* [VersionTags](#available-versiontags)
* [Configuration](#configuration)
* [TODOs](#todos)
* [Patchnotes](#patchnotes)
* [License](#license)

## Description

A simple Docker Container that allows you to monitor your other Docker Containers and backup your Docker Volumes, with notifictaions.

It creates a tarball (.tar) per configured volume, at the configured time, and stores it in the configured location :)
In the best case the Output directory is mapped to a Network Drive or another Host.

To see whats working now and whats planned in the near future see [TODOs](#todos).

Tested and developed for Linux. Windows, WSL, MacOS might be working, but not tested/designed for!

## Limitations/Warnings

Not recommended to use DockerRight as a solo backup solution for DBs or complex applications! But it's still better than no backup at all :D

It is a very young and developing Project, so use at your own risk!

## How to use

Minimal Docker run CMD:
``` bash
docker run -d 
    -v /var/run/docker.sock:/var/run/docker.sock
    -v /path/to/config:/opt/DockerRight/config
    -v /path/to/backupDir:/opt/DockerRight/backup
    # Optional: -v /path/to/logsDir:/opt/DockerRight/logs
    -e TZ=Europe/Berlin
    --name dockerright 
    ghcr.io/bata94/dockerright:latest
```

Minimal Docker-Compose File:
``` yaml
services:
  dockerright:
    container_name: dockerright
    image: ghcr.io/bata94/dockerright:latest
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - /path/to/config:/opt/DockerRight/config
      - /path/to/backupDir:/opt/DockerRight/backup
      # Optional: - /path/to/logsDir:/opt/DockerRight/logs
    environment:
      TZ: Europe/Berlin
```

To run the service successfully the containername needs to be "dockerright"!
Start the Container, it will stop after a few seconds on it's own.

Now you edit the created config.json. 

A restart option might be a good idea, for such a service, but it's easier to set it after the first run, to only generate the config.json.

If you don't set Before- or AfterBackupCMDs that would need network access, you should set the Docker Network to "none".

If you don't want to use the ContainerRegistry, just git clone this repository and build your image. The Dockerfile and an example Composefile are already there.

### Available VersionTags

| Tag       | Description                                   |
|-----------|-----------------------------------------------|
| `latest`  | The latest version                            |
| `X`       | The latest Major version                      |
| `0.X`     | The latest Minor version                      |
| `0.0.X`   | Specific patch (see releases for versions)    |

DockerRight uses Semantic Versioning, so you can lock the Image Version to a specific tag, as shown above.

When all ToDos are completed I will move to Major 1 :)\
Most ToDos should increase the Minor by 1.

[The Patchnotes can be found here](#patchnotes)

### Configuration

Parameters that can be set in the config.json. To reset them, delete the config.json and restart the container.

The Parameters in the config.json are typed in UpperCamelCase. To use those Parameters in environment variables, the must be typed as UPPER_SNAKE_CASE.

Parameters are evaluated as follows:

    1. EnvironmentVariables
    2. config.json
    3. default values

If you change a Parameter you will need to restart the DockerRightContainer to apply the change.

| Parameter (config.json)       | Parameter (EnvVar)               | Default                    | Type     | Description                                                            |
|-------------------------------|----------------------------------|----------------------------|----------|------------------------------------------------------------------------|
| EnableBackup                  | ENABLE_BACKUP                    | false                      | Bool     | Enable backup service                                                  |
| EnableMonitor                 | ENABLE_MONITOR                   | false                      | Bool     | Enable monitor service                                                 |
| MonitorIntervalSeconds        | MONITOR_INTERVAL_SECONDS         | 60                         | Int      | Interval in seconds                                                    |
| MonitorReties                 | MONITOR_RETIES                   | 5                          | Int      | Retries before sending notification                                    |
| BackupHours                   | BACKUP_HOURS                     | []                         | []Int    | Backup at these hours                                                  |
| RetentionHours                | RETENTION_HOURS                  | 120                        | Int      | Backup Retention in hours (24h * 5d)                                   |
| LogRetentionDays              | LOG_RETENTION_DAYS               | 7                          | Int      | Log Retention in days                                                  |
| ConcurrentBackupContainer     | CONCURRENT_BACKUP_CONTAINER      | numCPUs/2                  | Int      | How many mounts should be backed up at once                            |
| BackupPath                    | BACKUP_PATH                      | "/opt/DockerRight/backup"  | String   | Backup Path inside container (shouldn't be changed)                    |
| LogsPath                      | LOGS_PATH                        | "/opt/DockerRight/logs"    | String   | Logs Path inside container (shouldn't be changed)                      |
| BeforeBackupCMD               | BEFORE_BACKUP_CMD                | ""                         | String   | CMD to execute before backup                                           |
| AfterBackupCMD                | AFTER_BACKUP_CMD                 | ""                         | String   | CMD to execute after backup                                            |
| Log2File                      | LOG2FILE                         | false                      | Bool     | Toggle Log to File                                                     |
| LogLevel                      | LOG_LEVEL                        | "info"                     | String   | Set LogLevel (debug, info, warn, error, fatal, panic)                  |
| NotifyLevel                   | NOTIFY_LEVEL                     | "error"                    | String   | Set NotificationLevel (debug, info, warn, error, fatal, panic)         |
| BackupOnStartup               | BACKUP_ON_STARTUP                | false                      | Bool     | Start a Backup on startup, won't run again if started in a BackupHour  |
| CreateTestContainerOnStartup  | CREATE_TEST_CONTAINER_ON_STARTUP | true                       | Bool     | Create a TestContainer on startup, to check docker.sock                |
| NotifyLevel                   | NOTIFY_LEVEL                     | "warn"                     | String   | Set NotificationLevel (debug, info, warn, error, fatal, panic, none)   |
| TelegramBotToken              | TELEGRAM_BOT_TOKEN               | ""                         | String   | Telegram Bot Token [TelegramConf](#notifytelegram)                     |
| TelegramChatIDs               | TELEGRAM_CHAT_IDS                | []                         | []Int    | Telegram Chat IDs [TelegramConf](#notifytelegram)                      |

#### Notifications

If you want to get Notifications you will need to set the desired NotifyLevel, so all Logs in that Level (and above) will be send to the configured NotifyClients (i.e Telegram).

Container Monitoring will always be sent to all available clients.

##### NotifyTelegram

To enable Telegram notifications you will need to setup a BotToken and at least one ChatID.

To get a BotToken got to [BotFather](https://t.me/BotFather), create a new Bot and start it. It is necessary to create a new Bot for every DockerRightInstance!

The ChatID is more or less your ID, to get it you have two options:

    1. DockerRight is not running -> send a message to your created Bot and than visit >>https://api.telegram.org/bot<HIER_DEIN_BOT_TOKEN>/getUpdates<<
    2. DockerRight is running -> send a message to your created Bot and than watch the DockerRight logs. Under the WARN Flag there should pop up a LogMessage with your ID

## TODOs

What's planned in the near future? If you have any ideas, feature requests, suggestions or bug reports, please [open an issue](https://github.com/bata94/dockerRight/issues) or [create a PR](https://github.com/bata94/dockerRight/pulls) :)

Those points are roughly in order of importance (for me):

- [X] Create Backups per mount
- [X] Delete old Backups
- [X] Make config parameters settable by environment variables
- [X] Add VersionTag to startup console output
- [X] Better DockerBackupHelperNames, that reflect ContainerName and Mount
- [X] Enable concurrent backups
- [X] BackupContainer Output to File
- [X] Backup Docker Compose Files/Run Parameters
- [X] Logs to File
- [X] Monitor Docker Containers
- [X] Telegram Notifications
- [X] Fix Monitor only Loop
- [ ] Mount Container Volumes/Binds as read only, for safety
- [ ] Find reason for high CPU usage on low end CPUs (Zimaboard 40% in "idle"...)
- [X] Add Parameter to enable/disable log to File
- [ ] Deleting EnvVars do not overwrite config.json... (not sure how to fix it right now...)
- [ ] config.json FilePermissions
- [ ] Add FormatWrapper for Notify Package
- [ ] Refactor!!!
- [ ] Mail Notifications
- [ ] Discord Notifications
- [ ] Fine grain settings via Container Labels (like traefik for example)
- [ ] Restore Backups
- [ ] Image specific backup CMDs (i.e. for DBs, Nextcloud, Zammad, Mailcow etc.)
- [ ] SSH, SFTP, S3, NFS, SMB Backup location options
- [ ] Either configure Watchtower Container from DockerRight or program Watchtower functionality into DockerRight
- [ ] Refactor!!!
- [ ] WebUI for Configuration, Monitoring and Dashboard
- [ ] Add tests :D

## Patchnotes

If a specific Version is not listed here, eventhough it was released, it might only be a refactor or super minor change, without changes for the user.
As the development is rapid I might skip the patchnotes for a version!

### 0.2.1

- Telegram Notifications implemented
- DockerContainer Monitor implemented
- Added Parameter to toggle File logging

### 0.1.1

- Parameter now settable via Environment Variables
- Improved Struct logging (if "\n" are display instead of linebreaks, pls check your Docker LogDriver!)
- Disabled BackupRunner to run twice in one hour, if DockerRight is started in a defined BackupHour
- Changed BackupRunner ContainerNaming to better show what is running
- DockerRight will now print the current VersionTag to the startup console
- Added cuncurrent running of BackupRunners
- ContainerInfo is now in a file, inside the backup directory
- After-/BeforeBackupCMD Output as well as regular Logs are now written to /opt/DockerRight/logs/ (default value) directory
- BackupErrors will be found in the BackupDir as well

### 0.0.11

- Better logging
- Fixed a bug where the backup is blocked if only one backup time is defined

### 0.0.6

- Some refactoring and optimization
- Bugfixes, specially with finding the BackupPath of the DockerRightContainer on Host
- Build pipeline changes and CD testing

### 0.0.4

- Deletes old backups now

### 0.0.1

- Initial Release

## License

This project is licensed under the [Unlicense](https://unlicense.org/), so do what you want, but don't blame me :D 
