package configs

import (
	"github.com/spf13/viper"
	"io"
	"log/slog"
	"os"
	"strings"
	"wisebase/internal/handler"
	"wisebase/pkg/xslog"
)

type App struct {
	Addr string `yaml:"Addr"`
	Mode string `yaml:"Mode"`
}
type DB struct {
	Dialect     string `yaml:"Dialect"`
	DSN         string `yaml:"DSN"`
	MaxIdle     int    `yaml:"MaxIdle"`
	MaxActive   int    `yaml:"MaxActive"`
	MaxLifetime int    `yaml:"MaxLifetime"`
	AutoMigrate bool   `yaml:"AutoMigrate"`
}
type Redis struct {
	Addr     string `yaml:"Addr"`
	DB       int    `yaml:"DB"`
	Password string `yaml:"Password"`
}

type Config struct {
	Log          xslog.Config `yaml:"Log"`
	App          App          `yaml:"App"`
	Level        slog.Level
	ExtraWriters []ExtraWriter
	Handler      handler.Config `yaml:"Handler"`
	ReplaceAttr  func(groups []string, a slog.Attr) slog.Attr
}

func (c *Config) IsLocalOrDebugMode() bool {
	return c.IsLocalMode() || c.IsDebugMode()
}
func (c *Config) IsLocalMode() bool {
	return c.App.Mode == "local"
}
func (c *Config) IsDebugMode() bool {
	return c.App.Mode == "debug"
}

func (c *Config) IsReleaseMode() bool {
	return c.App.Mode == "release"
}

func InitConfig() (*Config, error) {
	var cfg Config
	configPath := "configs/prod.yaml"
	mode := os.Getenv("APP_MODE")
	if mode != "" {
		configPath = "configs/" + mode + ".config.yaml"
	}
	//log.Println("config path:", configPath)
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	viper.SetConfigType("yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	if err = viper.ReadConfig(file); err != nil {
		return nil, err
	}
	if err = viper.UnmarshalExact(&cfg); err != nil {
		return nil, err
	}
	return &cfg, err
}

type ExtraWriter struct {
	Writer io.Writer
	Level  slog.Level
}
