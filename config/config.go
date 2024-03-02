package config

import (
	"context"
	"fmt"
)

var Configuration Config

type Config struct {
	appConf   AppConfig
	infraConf InfraConfig
}

func Setup(ctx context.Context) error {
	if err := loadInfraConfiguration(); err != nil {
		return fmt.Errorf("error loading infrastructure config: %v", err)
	}

	if err := loadApplicationConfig(ctx); err != nil {
		return fmt.Errorf("error loading application config: %v", err)
	}

	return nil
}

func (c *Config) ApplicationConfiguration() AppConfig {
	return c.appConf
}

func (c *Config) SetApplicationConfiguration(ac AppConfig) {
	c.appConf = ac
}

func (c *Config) InfrastructureConfiguration() InfraConfig {
	return c.infraConf
}

func (c *Config) SetInfrastructureConfiguration(ic InfraConfig) {
	c.infraConf = ic
}
