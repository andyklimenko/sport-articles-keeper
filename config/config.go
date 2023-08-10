package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Feed   Feed
	Poller Poller
}

func (c *Config) Load() error {
	if err := c.Feed.load("feed"); err != nil {
		return fmt.Errorf("load feed config: %w", err)
	}

	if err := c.Poller.load("poll"); err != nil {
		return fmt.Errorf("load poller config: %w", err)
	}

	return nil
}

func setupViper(envPrefix string) *viper.Viper {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetEnvPrefix(envPrefix)
	return v
}
