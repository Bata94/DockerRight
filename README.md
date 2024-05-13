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
    environment:
      TZ: Europe/Berlin
```

To run the service successfully the containername needs to be "dockerright"!
Start the Container, it will stop after a few seconds on it's own.

Now you edit the created config.json. 

A restart option might be a good idea, for such a service, but it's easier to set it after the first run, to only generate the config.json.

If you don't set Before- or AfterBackupCMDs that would need network access, you should set the Docker Network to "none".

### Available VersionTags

| Tag       | Description                                   |
|-----------|-----------------------------------------------|
| `latest`  | The latest stable version                     |
| `0`       | The latest Major version                      |
| `0.0`     | The latest Minor version                      |
| `0.0.X`   | Specific patch (see releases for versions)    |

DockerRight uses Semantic Versioning, so you can lock the Image Version to a specific tag, as shown above.

When all ToDos are completed I will move to Major 1 :)\
Most ToDos should increase the Minor by 1.

[The Patchnotes can be found here](#patchnotes)

### Configuration

Parameters that can be set in the config.json. To reset them, delete the config.json and restart the container.

| Parameter                     | Default                    | Type     | Description                                                   |
|-------------------------------|----------------------------|----------|---------------------------------------------------------------|
| EnableBackup                  | false                      | Bool     | Enable backup service                                         |
| EnableMonitor                 | false                      | Bool     | Enable monitor service                                        |
| MonitorIntervalSeconds        | 60                         | Int      | Interval in seconds                                           |
| MonitorReties                 | 5                          | Int      | Retries before sending notification                           |
| BackupHours                   | []                         | []Int    | Backup at these hours                                         |
| RetentionHours                | 120                        | Int      | Retention in hours (24 * 5)                                   |
| ConcurrentBackupContainer     | numCPUs/2                  | Int      | How many mounts should be backed up at once                   |
| BackupPath                    | "/opt/DockerRight/backup"  | String   | Backup Path inside container (shouldn't be changed)           |
| BeforeBackupCMD               | ""                         | String   | CMD to execute before backup                                  |
| AfterBackupCMD                | ""                         | String   | CMD to execute after backup                                   |
| LogLevel                      | "info"                     | String   | Set LogLevel (debug, info, warn, error, fatal, panic)         |
| BackupOnStartup               | false                      | Bool     | Start a Backup on startup                                     |
| CreateTestContainerOnStartup  | true                       | Bool     | Create a TestContainer on startup, to check docker.sock       |

## TODOs

What's planned in the near future? If you have any ideas, feature requests, suggestions or bug reports, please [open an issue](https://github.com/bata94/dockerRight/issues) or [create a PR](https://github.com/bata94/dockerRight/pulls) :)

Those points are roughly in order of importance (for me):

- [X] Create Backups per mount
- [X] Delete old Backups
- [ ] Monitor Docker Containers
- [ ] Telegram Notifications
- [ ] Mail Notifications
- [ ] Discord Notifications
- [ ] Logs to File
- [ ] Fix Monitor only Loop
- [ ] Enable concurrent backups
- [ ] BackupContainer Output to File
- [ ] Make config parameters settable by environment variables
- [ ] Mount Container Volumes/Binds as read only, for safety
- [ ] Fine grain settings via Container Labels (like traefik for example)
- [ ] Add tests
- [ ] Backup Docker Compose Files/Run Parameters
- [ ] Restore Backups
- [ ] Image specific backup CMDs (i.e. for DBs, Nextcloud, Zammad, Mailcow etc.)
- [ ] WebUI for Configuration, Monitoring and Dashboard

## Patchnotes

If a specific Version is not listed here, eventhough it was released, it might only be a refactor or super minor change, without changes for the user.
As the development is rapid I might skip the patchnotes for a version!

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
