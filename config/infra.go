package config

import (
	"fmt"

	"github.com/caarlos0/env/v9"
)

type InfraConfig struct {
	Database Database
	Otel     Otel
}

type Database struct {
	Host     string `env:"DATABASE_HOST"`
	Port     int    `env:"DATABASE_PORT"`
	Name     string `env:"DATABASE_NAME"`
	User     string `env:"DATABASE_USER" json:"-"`
	Password string `env:"DATABASE_PASSWORD" json:"-"`
}

type Otel struct {
	Host string `env:"OTEL_EXPORTER_OTLP_HOST"`
	Port int    `env:"OTEL_EXPORTER_OTLP_PORT"`
}

func loadInfraConfiguration() error {
	inf := InfraConfig{}
	if err := env.Parse(&inf); err != nil {
		return err
	}

	Configuration.infraConf = inf
	return nil
}

func (c *Config) DBHost() string {
	return c.infraConf.Database.Host
}

func (c *Config) DBPort() int {
	return c.infraConf.Database.Port
}

func (c *Config) DBName() string {
	return c.infraConf.Database.Name
}

func (c *Config) DBUser() string {
	return c.infraConf.Database.User
}

func (c *Config) DBPassword() string {
	return c.infraConf.Database.Password
}

func (c *Config) OtelURL() string {
	return fmt.Sprintf("%s:%d", c.infraConf.Otel.Host, c.infraConf.Otel.Port)
}
