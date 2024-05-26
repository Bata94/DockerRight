package config

import (
	"encoding/json"
	"errors"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/bata94/DockerRight/internal/log"
)

var (
	ConfigPath string
	Conf       Config
)

func Init(configPath string) Config {
	log.Info("Initializing Config Module")
	var err error
	ConfigPath = configPath

	if ConfigPath == "" {
		ConfigPath = "./config.json"
	}

	Conf = Config{}
	err = Conf.SetDefaults()
	if err != nil {
		log.Fatal("Error setting defaults: ", err)
	}

	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		log.Info("Config file not found generating defaults!")
		err = Conf.Save()
		if err != nil {
			log.Fatal("Error saving Config: ", err)
		}
	}

	err = Conf.Load()
	if err != nil {
		log.Fatal("Error loading Config: ", err)
	}
	err = Conf.Save()
	if err != nil {
		log.Fatal("Error saving Config: ", err)
	}
	log.Info("Config loaded")

	return Conf
}

type Config struct {
	Version                      string
	EnableBackup                 bool
	EnableMonitor                bool
	MonitorIntervalSeconds       int
	MonitorRetries               int
	BackupHours                  []int
	RetentionHours               int
	LogRetentionDays             int
	ConcurrentBackupContainer    int
	BackupPath                   string
	LogsPath                     string
	BeforeBackupCMD              string
	AfterBackupCMD               string
	LogLevel                     string
	BackupOnStartup              bool
	CreateTestContainerOnStartup bool
}

func (c *Config) SetDefaults() error {
	log.Info("Config SetDefaults")

	c.Version = "unknown"
	c.EnableBackup = false
	c.EnableMonitor = false
	c.RetentionHours = 24 * 5
	c.LogRetentionDays = 7
	c.MonitorIntervalSeconds = 60
	c.MonitorRetries = 5
	c.BackupHours = []int{}
	c.ConcurrentBackupContainer = (runtime.NumCPU() / 2)
	c.BackupPath = "/opt/DockerRight/backup"
	c.LogsPath = "/opt/DockerRight/logs"
	c.BeforeBackupCMD = ""
	c.AfterBackupCMD = ""
	c.LogLevel = "info"
	c.BackupOnStartup = false
	c.CreateTestContainerOnStartup = true

	return nil
}

func (c *Config) Load() error {
	log.Info("Config Loading")
	err := Conf.LoadFromFile()
	if err != nil {
		return err
	}
	err = Conf.LoadFromEnv()
	if err != nil {
		return err
	}
	log.Info("Config: ", log.FormatStruct(Conf))
	return nil
}

func (c *Config) LoadFromEnv() error {
	log.Info("Config LoadFromEnv")

	// TODO: Refactor into functions
	// Bool Values
	if os.Getenv("ENABLE_BACKUP") != "" {
		enableBackupVar := strings.ToLower(os.Getenv("ENABLE_BACKUP"))
		if enableBackupVar == "true" {
			c.EnableBackup = true
		} else if enableBackupVar == "false" {
			c.EnableBackup = false
		} else {
			log.Error("Environment Variable 'ENABLE_BACKUP' could not be parsed... value read: ", enableBackupVar)
			log.Warn("Falling back to value in 'config.json' or to default value!")
		}
	}
	if os.Getenv("ENABLE_MONITOR") != "" {
		enableMonitorVar := strings.ToLower(os.Getenv("ENABLE_MONITOR"))
		if enableMonitorVar == "true" {
			c.EnableMonitor = true
		} else if enableMonitorVar == "false" {
			c.EnableMonitor = false
		} else {
			log.Error("Environment Variable 'ENABLE_BACKUP' could not be parsed... value read: ", enableMonitorVar)
			log.Warn("Falling back to value in 'config.json' or to default value!")
		}
	}
	if os.Getenv("BACKUP_ON_STARTUP") != "" {
		backupOnStartupVar := strings.ToLower(os.Getenv("BACKUP_ON_STARTUP"))
		if backupOnStartupVar == "true" {
			c.BackupOnStartup = true
		} else if backupOnStartupVar == "false" {
			c.BackupOnStartup = false
		} else {
			log.Error("Environment Variable 'ENABLE_MONITOR' could not be parsed... value read: ", backupOnStartupVar)
			log.Warn("Falling back to value in 'config.json' or to default value!")
		}
	}
	if os.Getenv("CREATE_TEST_CONTAINER_ON_STARTUP") != "" {
		createTestContainerOnStartup := strings.ToLower(os.Getenv("CREATE_TEST_CONTAINER_ON_STARTUP"))
		if createTestContainerOnStartup == "true" {
			c.BackupOnStartup = true
		} else if createTestContainerOnStartup == "false" {
			c.BackupOnStartup = false
		} else {
			log.Error("Environment Variable 'CREATE_TEST_CONTAINER_ON_STARTUP' could not be parsed... value read: ", createTestContainerOnStartup)
			log.Warn("Falling back to value in 'config.json' or to default value!")
		}
	}

	// String Values
	if os.Getenv("BACKUP_PATH") != "" {
		c.BackupPath = os.Getenv("BACKUP_PATH")
	}
	if os.Getenv("LOGS_PATH") != "" {
		c.LogsPath = os.Getenv("LOGS_PATH")
	}
	if os.Getenv("BEFORE_BACKUP_CMD") != "" {
		c.BeforeBackupCMD = os.Getenv("BEFORE_BACKUP_CMD")
	}
	if os.Getenv("AFTER_BACKUP_CMD") != "" {
		c.AfterBackupCMD = os.Getenv("AFTER_BACKUP_CMD")
	}
	if os.Getenv("LOG_LEVEL") != "" {
		c.LogLevel = os.Getenv("LOG_LEVEL")
	}

	// List Values
	if os.Getenv("BACKUP_HOURS") != "" {
		backupHoursVar := os.Getenv("BACKUP_HOURS")
		backupHours := []int{}

		// String cleanup
		backupHoursVar = strings.ReplaceAll(backupHoursVar, "[", "")
		backupHoursVar = strings.ReplaceAll(backupHoursVar, "]", "")
		backupHoursVar = strings.ReplaceAll(backupHoursVar, " ", "")

		backupHoursStrList := strings.Split(backupHoursVar, ",")

		for _, s := range backupHoursStrList {
			hourInt, err := strconv.Atoi(s)

			if err != nil {
				log.Debug(err, s)

				log.Error("Environment Variable 'BACKUP_HOURS' could not be parsed... value read: ", backupHoursVar)
				log.Warn("Falling back to value in 'config.json' or to default value!")
				break
			}

			backupHours = append(backupHours, hourInt)
		}

		if len(backupHours) != 0 {
			c.BackupHours = backupHours
		}
	}

	// Int Values
	if os.Getenv("MONITOR_INTERVAL_SECONDS") != "" {
		val := os.Getenv("MONITOR_INTERVAL_SECONDS")
		valInt, err := strconv.Atoi(val)

		if err != nil {
			log.Debug(err)
			log.Error("Environment Variable 'MONITOR_INTERVAL_SECONDS' could not be parsed... value read: ", val)
			log.Warn("Falling back to value in 'config.json' or to default value!")
		} else {
			c.MonitorIntervalSeconds = valInt
		}
	}
	if os.Getenv("MONITOR_RETIES") != "" {
		val := os.Getenv("MONITOR_RETIES")
		valInt, err := strconv.Atoi(val)

		if err != nil {
			log.Debug(err)
			log.Error("Environment Variable 'MONITOR_RETIES' could not be parsed... value read: ", val)
			log.Warn("Falling back to value in 'config.json' or to default value!")
		} else {
			c.MonitorRetries = valInt
		}
	}
	if os.Getenv("RETENTION_HOURS") != "" {
		val := os.Getenv("RETENTION_HOURS")
		valInt, err := strconv.Atoi(val)

		if err != nil {
			log.Debug(err)
			log.Error("Environment Variable 'RETENTION_HOURS' could not be parsed... value read: ", val)
			log.Warn("Falling back to value in 'config.json' or to default value!")
		} else {
			c.RetentionHours = valInt
		}
	}
	if os.Getenv("LOG_RETENTION_DAYS") != "" {
		val := os.Getenv("LOG_RETENTION_DAYS")
		valInt, err := strconv.Atoi(val)

		if err != nil {
			log.Debug(err)
			log.Error("Environment Variable 'LOG_RETENTION_DAYS' could not be parsed... value read: ", val)
			log.Warn("Falling back to value in 'config.json' or to default value!")
		} else {
			c.LogRetentionDays = valInt
		}
	}
	if os.Getenv("CONCURRENT_BACKUP_CONTAINER") != "" {
		val := os.Getenv("CONCURRENT_BACKUP_CONTAINER")
		valInt, err := strconv.Atoi(val)

		if err != nil {
			log.Debug(err)
			log.Error("Environment Variable 'CONCURRENT_BACKUP_CONTAINER' could not be parsed... value read: ", val)
			log.Warn("Falling back to value in 'config.json' or to default value!")
		} else {
			c.ConcurrentBackupContainer = valInt
		}
	}

	return nil
}

func (c *Config) SetVersion() {
	if os.Getenv("VERSION") != "" {
		c.Version = os.Getenv("VERSION")
	} else {
		c.Version = "unknown"
	}

	log.Info("")
	log.Info("")
	log.Info("####-------------------------####")
	log.Info("")
	log.Info("    DockerRight Version: ", c.Version)
	log.Info("")
	log.Info("####-------------------------####")
	log.Info("")
	log.Info("")
}

func (c *Config) LoadFromFile() error {
	log.Info("Config LoadFromFile")

	confFile, err := os.ReadFile(ConfigPath)
	if err != nil {
		return errors.New("Error reading config file: " + err.Error())
	}
	err = json.Unmarshal(confFile, &Conf)
	if err != nil {
		return errors.New("Error unmarshalling config file: " + err.Error())
	}

	return nil
}

func (c *Config) Save() error {
	log.Info("Config Saving")

	Conf.SetVersion()

	confFile, err := json.MarshalIndent(Conf, "", " ")
	if err != nil {
		return errors.New("Error marshalling config file: " + err.Error())
	}
	err = os.WriteFile(ConfigPath, confFile, 0o644)
	if err != nil {
		return errors.New("Error writing config file: " + err.Error())
	}

	return nil
}
