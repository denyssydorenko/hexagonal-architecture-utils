package config

import (
	"context"
	"errors"
	"fmt"
	"hexagonal-architexture-utils/internal/pkg/logging"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const (
	externalAppConfig = "/etc/config/app/"
	embeddedAppConfig = "config/"
)

var BuildVersion string

type AppConfig struct {
	Http Http
	Db   DB
}

type Http struct {
	HostAddress     string        `mapstructure:"HOST_ADDRESS"`
	BasePath        string        `mapstructure:"BASE_PATH"`
	ShutdownTimeout time.Duration `mapstructure:"SHUTDOWN_TIMEOUT"`
	ServiceName     string        `mapstructure:"SERVICE_NAME"`
	Timeout         time.Duration `mapstructure:"TIMEOUT"`
}

type DB struct {
	MinOpenConns int `mapstructure:"MIN_OPEN_CONNECTIONS"`
}

func loadApplicationConfig(ctx context.Context) error {
	application := AppConfig{}

	if err := loadAppConfiguration(ctx, "config", &application, true); err != nil {
		return err
	}

	Configuration.appConf = application
	return nil
}

func loadAppConfiguration(ctx context.Context, name string, configType interface{}, useEnv bool) error {
	if _, err := os.Stat(fmt.Sprintf("%v%v.yaml", externalAppConfig, name)); errors.Is(err, os.ErrNotExist) {
		logging.Global.Info(ctx, fmt.Sprintf("no environment config found, using default config"))
		viper.AddConfigPath(embeddedAppConfig)
	} else {
		logging.Global.Info(ctx, "using environment config")
		viper.AddConfigPath(externalAppConfig)
	}

	viper.SetConfigName(name)
	viper.SetConfigType("yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if useEnv {
		viper.AutomaticEnv()
	}

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&configType)
	if err != nil {
		return err
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		err = viper.Unmarshal(&configType)
		if err != nil {
			logging.Global.Error(ctx, fmt.Sprintf("failed loading new config: %v", err))
		} else {
			logging.Global.Info(ctx, "updated config successfully...")
		}
	})
	viper.WatchConfig()

	return nil
}

func (c *Config) ApplicationConfig() *AppConfig {
	return &c.appConf
}

func (c *Config) HostAddress() string {
	return c.appConf.Http.HostAddress
}

func (c *Config) BasePath() string {
	return c.appConf.Http.BasePath
}

func (c *Config) ShutDownTimeout() time.Duration {
	return c.appConf.Http.ShutdownTimeout
}

func (c *Config) ServiceName() string {
	return c.appConf.Http.ServiceName
}

func (c *Config) Timeout() time.Duration {
	return c.appConf.Http.Timeout
}

func (c *Config) MinOpenConns() int {
	return c.appConf.Db.MinOpenConns
}
