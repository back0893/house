package common

import "github.com/spf13/viper"

var GlobalConfig *Config

type Config struct {
	*viper.Viper
}

func NewConfig() *Config {
	return &Config{
		viper.GetViper(),
	}

}
func (c *Config) Load() {
	c.SetConfigType("yaml")
	c.SetConfigFile("./env.yml")
}
