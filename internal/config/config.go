package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Log       LogConfig       `mapstructure:"log"`
	API       APIConfig       `mapstructure:"api"`
	Cron      CronConfig      `mapstructure:"cron"`
	Retention RetentionConfig `mapstructure:"retention"`
	DB        DBConfig        `mapstructure:"db"`
	Storage   StorageConfig   `mapstructure:"storage"`
	Admin     AdminConfig     `mapstructure:"admin"`
	Token     TokenConfig     `mapstructure:"token"`
	Feature   FeatureConfig   `mapstructure:"feature"`
}

type ServerConfig struct {
	Port    int    `mapstructure:"port"`
	BaseURL string `mapstructure:"base_url"`
}

type LogConfig struct {
	Level string `mapstructure:"level"`
}

type APIConfig struct {
	Mode string `mapstructure:"mode"` // local | redirect
}

type CronConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	DailySpec string `mapstructure:"daily_spec"`
}

type RetentionConfig struct {
	Days int `mapstructure:"days"`
}

type DBConfig struct {
	Type string `mapstructure:"type"` // sqlite/mysql/postgres
	DSN  string `mapstructure:"dsn"`
}

type StorageConfig struct {
	Type   string       `mapstructure:"type"` // local/s3/webdav
	Local  LocalConfig  `mapstructure:"local"`
	S3     S3Config     `mapstructure:"s3"`
	WebDAV WebDAVConfig `mapstructure:"webdav"`
}

type LocalConfig struct {
	Root string `mapstructure:"root"`
}

type S3Config struct {
	Endpoint        string `mapstructure:"endpoint"`
	Region          string `mapstructure:"region"`
	Bucket          string `mapstructure:"bucket"`
	AccessKey       string `mapstructure:"access_key"`
	SecretKey       string `mapstructure:"secret_key"`
	PublicURLPrefix string `mapstructure:"public_url_prefix"`
	ForcePathStyle  bool   `mapstructure:"force_path_style"`
}

type WebDAVConfig struct {
	URL             string `mapstructure:"url"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	PublicURLPrefix string `mapstructure:"public_url_prefix"`
}

type AdminConfig struct {
	PasswordBcrypt string `mapstructure:"password_bcrypt"`
}

type TokenConfig struct {
	DefaultTTL string `mapstructure:"default_ttl"`
}

type FeatureConfig struct {
	WriteDailyFiles bool `mapstructure:"write_daily_files"`
}

// Bing 默认配置 (内置)
const (
	BingMkt     = "zh-CN"
	BingFetchN  = 8
	BingAPIBase = "https://www.bing.com/HPImageArchive.aspx"
)

var (
	GlobalConfig *Config
	configLock   sync.RWMutex
	v            *viper.Viper
)

func Init(configPath string) error {
	v = viper.New()
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath("./data")
		v.AddConfigPath(".")
	}

	v.SetDefault("server.port", 8080)
	v.SetDefault("log.level", "info")
	v.SetDefault("api.mode", "local")
	v.SetDefault("cron.enabled", true)
	v.SetDefault("cron.daily_spec", "0 10 * * *")
	v.SetDefault("retention.days", 30)
	v.SetDefault("db.type", "sqlite")
	v.SetDefault("db.dsn", "data/bing_paper.db")
	v.SetDefault("storage.type", "local")
	v.SetDefault("storage.local.root", "data/picture")
	v.SetDefault("token.default_ttl", "168h")
	v.SetDefault("feature.write_daily_files", true)
	v.SetDefault("admin.password_bcrypt", "$2a$10$fYHPeWHmwObephJvtlyH1O8DIgaLk5TINbi9BOezo2M8cSjmJchka") // 默认密码: admin123

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
		// 如果文件不存在，我们使用默认值并尝试创建一个默认配置文件
		fmt.Println("Config file not found, creating default config at ./data/config.yaml")
		if err := v.SafeWriteConfigAs("./data/config.yaml"); err != nil {
			fmt.Printf("Warning: Failed to create default config file: %v\n", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return err
	}

	GlobalConfig = &cfg

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		var newCfg Config
		if err := v.Unmarshal(&newCfg); err == nil {
			configLock.Lock()
			GlobalConfig = &newCfg
			configLock.Unlock()
		}
	})
	v.WatchConfig()

	return nil
}

func GetConfig() *Config {
	configLock.RLock()
	defer configLock.RUnlock()
	return GlobalConfig
}

func SaveConfig(cfg *Config) error {
	v.Set("server", cfg.Server)
	v.Set("log", cfg.Log)
	v.Set("api", cfg.API)
	v.Set("cron", cfg.Cron)
	v.Set("retention", cfg.Retention)
	v.Set("db", cfg.DB)
	v.Set("storage", cfg.Storage)
	v.Set("admin", cfg.Admin)
	v.Set("token", cfg.Token)
	v.Set("feature", cfg.Feature)
	return v.WriteConfig()
}

func GetRawViper() *viper.Viper {
	return v
}

func GetTokenTTL() time.Duration {
	ttl, err := time.ParseDuration(GetConfig().Token.DefaultTTL)
	if err != nil {
		return 168 * time.Hour
	}
	return ttl
}
