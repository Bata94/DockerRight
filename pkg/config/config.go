package config

import (
  "DockerRight/pkg/log"

  "os"
  "encoding/json"
  "runtime"
)

var ConfigPath string
var Conf Config

func Init(configPath string) Config {
  log.Info("Initializing Config Module") 
  ConfigPath = configPath

  if ConfigPath == "" {
    ConfigPath = "./config.json"
  }

  Conf = Config{}
  Conf.SetDefaults()

  if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
    log.Info("Config file not found generating defaults!")
    Conf.Save()
  }

  Conf.Load()
  Conf.Save()
  log.Info("Config loaded")

  return Conf
}

type Config struct {
  IncludeContainers []string
  ExcludeContainers []string
  IncludeVolumes []string
  ExcludeVolumes []string
  IncludeMounts []string
  ExcludeMounts []string
  MonitorIntervalSeconds int
  MonitorRetries int
  BackupHours []int
  RetentionHours int
  ConcurrentBackupContainer int
  BackupPath string
  BeforeBackupCMD string
  AfterBackupCMD string
  LogLevel string
  BackupOnStartup bool
  CreateTestContainerOnStartup bool
}

func (c *Config) SetDefaults() error {
  log.Info("Config SetDefaults")

  c.IncludeContainers = []string{}
  c.ExcludeContainers = []string{}
  c.IncludeVolumes = []string{}
  c.ExcludeVolumes = []string{}
  c.IncludeMounts = []string{}
  c.ExcludeMounts = []string{}
  c.RetentionHours = 24 * 5
  c.MonitorIntervalSeconds = 60
  c.MonitorRetries = 5
  c.BackupHours = []int{2}
  c.ConcurrentBackupContainer = (runtime.NumCPU() / 2)
  c.BackupPath = "./testDirs/backupDir/"
  c.BeforeBackupCMD = ""
  c.AfterBackupCMD = ""
  c.LogLevel = "debug"
  c.BackupOnStartup = false
  c.CreateTestContainerOnStartup = true

  return nil 
}

func (c *Config) Load() error {
  log.Info("Config Loading")
  Conf.LoadFromFile()
  Conf.LoadFromEnv()
  return nil
}

func (c *Config) LoadFromEnv() error {
  log.Info("Config LoadFromEnv")

  // TODO: Implement this!

  // if os.Getenv("INCLUDE_CONTAINERS") != "" {
  //   c.IncludeContainers = os.Getenv("INCLUDE_CONTAINERS").Split(",")
  // }
  // if os.Getenv("EXCLUDE_CONTAINERS") != "" {
  //   c.ExcludeContainers = os.Getenv("EXCLUDE_CONTAINERS").Split(",")
  // }
  // if os.Getenv("INCLUDE_VOLUMES") != "" {
  //   c.IncludeVolumes = os.Getenv("INCLUDE_VOLUMES").Split(",")
  // }
  // if os.Getenv("EXCLUDE_VOLUMES") != "" {
  //   c.ExcludeVolumes = os.Getenv("EXCLUDE_VOLUMES").Split(",")
  // }
  // if os.Getenv("INCLUDE_MOUNTS") != "" {
  //   c.IncludeMounts = os.Getenv("INCLUDE_MOUNTS").Split(",")
  // }
  // if os.Getenv("EXCLUDE_MOUNTS") != "" {
  //   c.ExcludeMounts = os.Getenv("EXCLUDE_MOUNTS").Split(",")
  // }
  // if os.Getenv("RETENTION_HOURS") != "" {
  //   c.RetentionHours = os.Getenv("RETENTION_HOURS").ToInt()
  // }
  // if os.Getenv("CONCURRENT_BACKUP_CONTAINER") != "" {
  //   c.ConcurrentBackupContainer = os.Getenv("CONCURRENT_BACKUP_CONTAINER").ToInt()
  // }
  if os.Getenv("BACKUP_PATH") != "" {
    c.BackupPath = os.Getenv("BACKUP_PATH")
  }
  if os.Getenv("BEFORE_BACKUP_CMD") != "" {
    c.BeforeBackupCMD = os.Getenv("BEFORE_BACKUP_CMD")
  }
  if os.Getenv("AFTER_BACKUP_CMD") != "" {
    c.AfterBackupCMD = os.Getenv("AFTER_BACKUP_CMD")
  }
  // if os.Getenv("BACKUP_ON_STARTUP") != "" {
  //   c.BackupOnStartup = os.Getenv("BACKUP_ON_STARTUP").ToBool()
  // }
  // if os.Getenv("CREATE_TEST_CONTAINER_ON_STARTUP") != "" {
  //   c.CreateTestContainerOnStartup = os.Getenv("CREATE_TEST_CONTAINER_ON_STARTUP").ToBool()
  // }

  return nil
}

func (c *Config) LoadFromFile() error {
  log.Info("Config LoadFromFile")

  confFile, err := os.ReadFile(ConfigPath)
  if err != nil {
    log.Fatal("Error reading config file" + err.Error())
  }
  err = json.Unmarshal(confFile, &Conf)
  if err != nil {
    log.Fatal("Error unmarshalling config file" + err.Error())
  }

  return nil
}

func (c *Config) Save() error {
  log.Info("Config Saving")
  
  confFile, err := json.MarshalIndent(Conf, "", " ")
  if err != nil {
    log.Fatal("Error marshalling config file" + err.Error())
  }
  err = os.WriteFile(ConfigPath, confFile, 0644)
  if err != nil {
    log.Fatal("Error writing config file" + err.Error())
  }

  return nil
}
