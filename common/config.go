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
func (c *Config) Load(configType, file string) error {
	c.SetConfigType(configType)
	c.SetConfigFile(file)
	if err := c.ReadInConfig(); err != nil {
		return err
	}
	return nil
}
