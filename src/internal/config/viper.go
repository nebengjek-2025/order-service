package config

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

var (
	cfg  *viper.Viper
	once sync.Once
)

func NewViper() *viper.Viper {
	once.Do(func() {
		c := viper.New()
		if raw := os.Getenv("ORDER_CONFIG_JSON"); raw != "" {
			c.SetConfigType("json")
			if err := c.ReadConfig(strings.NewReader(raw)); err != nil {
				panic(fmt.Errorf("failed to read config from ORDER_CONFIG_JSON: %w", err))
			}
		} else {
			c.SetConfigName("config")
			c.SetConfigType("json")
			c.AddConfigPath("./../")
			c.AddConfigPath("./")
			if err := c.ReadInConfig(); err != nil {
				panic(fmt.Errorf("Fatal error config file: %w \n", err))
			}
		}

		c.AutomaticEnv()
		cfg = c
	})
	return cfg
}

func GetConfig() *viper.Viper {
	return NewViper()
}
