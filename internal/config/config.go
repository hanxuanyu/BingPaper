package config

import (
	"fmt"
	"os"
	"sort"
	"strings"
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
	Web       WebConfig       `mapstructure:"web"`
}

type ServerConfig struct {
	Port    int    `mapstructure:"port"`
	BaseURL string `mapstructure:"base_url"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`     // 业务日志文件名
	DBFilename string `mapstructure:"db_filename"`  // 数据库日志文件名
	MaxSize    int    `mapstructure:"max_size"`     // 每个日志文件最大大小 (MB)
	MaxBackups int    `mapstructure:"max_backups"`  // 保留旧日志文件最大个数
	MaxAge     int    `mapstructure:"max_age"`      // 保留旧日志文件最大天数
	Compress   bool   `mapstructure:"compress"`     // 是否压缩旧日志文件
	LogConsole bool   `mapstructure:"log_console"`  // 是否同时输出到控制台
	ShowDBLog  bool   `mapstructure:"show_db_log"`  // 是否在控制台显示数据库日志
	DBLogLevel string `mapstructure:"db_log_level"` // 数据库日志级别: debug, info, warn, error
}

func (c LogConfig) GetLevel() string      { return c.Level }
func (c LogConfig) GetFilename() string   { return c.Filename }
func (c LogConfig) GetDBFilename() string { return c.DBFilename }
func (c LogConfig) GetMaxSize() int       { return c.MaxSize }
func (c LogConfig) GetMaxBackups() int    { return c.MaxBackups }
func (c LogConfig) GetMaxAge() int        { return c.MaxAge }
func (c LogConfig) GetCompress() bool     { return c.Compress }
func (c LogConfig) GetLogConsole() bool   { return c.LogConsole }
func (c LogConfig) GetShowDBLog() bool    { return c.ShowDBLog }
func (c LogConfig) GetDBLogLevel() string { return c.DBLogLevel }

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

type WebConfig struct {
	Path string `mapstructure:"path"`
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

	// OnDBConfigChange 当数据库配置发生变更时的回调函数
	OnDBConfigChange func(newCfg *Config)
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
	v.SetDefault("log.filename", "data/logs/app.log")
	v.SetDefault("log.db_filename", "data/logs/db.log")
	v.SetDefault("log.max_size", 100)
	v.SetDefault("log.max_backups", 3)
	v.SetDefault("log.max_age", 7)
	v.SetDefault("log.compress", true)
	v.SetDefault("log.log_console", true)
	v.SetDefault("log.show_db_log", false)
	v.SetDefault("log.db_log_level", "info")
	v.SetDefault("api.mode", "local")
	v.SetDefault("cron.enabled", true)
	v.SetDefault("cron.daily_spec", "20 8-23/4 * * *")
	v.SetDefault("retention.days", 0)
	v.SetDefault("db.type", "sqlite")
	v.SetDefault("db.dsn", "data/bing_paper.db")
	v.SetDefault("storage.type", "local")
	v.SetDefault("storage.local.root", "data/picture")
	v.SetDefault("token.default_ttl", "168h")
	v.SetDefault("feature.write_daily_files", true)
	v.SetDefault("web.path", "web")
	v.SetDefault("admin.password_bcrypt", "$2a$10$fYHPeWHmwObephJvtlyH1O8DIgaLk5TINbi9BOezo2M8cSjmJchka") // 默认密码: admin123

	// 绑定环境变量
	v.SetEnvPrefix("BINGPAPER")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		// 如果指定了配置文件但读取失败（且不是找不到文件的错误），或者没指定但也没找到
		_, isNotFound := err.(viper.ConfigFileNotFoundError)
		// 如果显式指定了文件，viper 报错可能不是 ConfigFileNotFoundError 而是 os.PathError
		if !isNotFound && configPath != "" {
			if _, statErr := os.Stat(configPath); os.IsNotExist(statErr) {
				isNotFound = true
			}
		}

		if !isNotFound {
			return err
		}

		// 如果文件不存在，我们使用默认值并尝试创建一个默认配置文件
		targetConfigPath := configPath
		if targetConfigPath == "" {
			targetConfigPath = "data/config.yaml"
		}
		fmt.Printf("Config file not found, creating default config at %s\n", targetConfigPath)
		if err := v.SafeWriteConfigAs(targetConfigPath); err != nil {
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
			oldDBConfig := GlobalConfig.DB
			GlobalConfig = &newCfg
			newDBConfig := newCfg.DB
			configLock.Unlock()

			// 检查数据库配置是否发生变更
			if oldDBConfig.Type != newDBConfig.Type || oldDBConfig.DSN != newDBConfig.DSN {
				// 触发数据库迁移逻辑
				// 这里由于循环依赖问题，我们可能需要通过回调或者一个统一的 Reload 函数来处理
				if OnDBConfigChange != nil {
					OnDBConfigChange(&newCfg)
				}
			}
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
	v.Set("web", cfg.Web)
	return v.WriteConfig()
}

func GetRawViper() *viper.Viper {
	return v
}

// GetAllSettings 返回所有生效配置项
func GetAllSettings() map[string]interface{} {
	return v.AllSettings()
}

// GetFormattedSettings 以 key: value 形式返回所有配置项的字符串
func GetFormattedSettings() string {
	keys := v.AllKeys()
	sort.Strings(keys)
	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("%s: %v\n", k, v.Get(k)))
	}
	return sb.String()
}

// GetEnvOverrides 返回环境变量覆盖详情（已排序）
func GetEnvOverrides() []string {
	var overrides []string
	keys := v.AllKeys()
	sort.Strings(keys)
	for _, key := range keys {
		// 根据 viper 的配置生成对应的环境变量名
		// Prefix: BINGPAPER, KeyReplacer: . -> _
		envKey := strings.ToUpper(fmt.Sprintf("BINGPAPER_%s", strings.ReplaceAll(key, ".", "_")))
		if val, ok := os.LookupEnv(envKey); ok {
			overrides = append(overrides, fmt.Sprintf("%s: %s=%s", key, envKey, val))
		}
	}
	return overrides
}

func GetTokenTTL() time.Duration {
	ttl, err := time.ParseDuration(GetConfig().Token.DefaultTTL)
	if err != nil {
		return 168 * time.Hour
	}
	return ttl
}
