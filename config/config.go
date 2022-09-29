package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Port          int
	Verbose       int
	RedisAddress  string
	RedisPassword string
	RedisDB       int
}

var config *Config

func New() (*Config, error) {
	if config == nil {
		var err error
		config, err = newConfig()
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

func newConfig() (*Config, error) {
	viper.SetDefault("port", 8088)
	viper.SetDefault("verbose", 6) // 6 Debug // 4 info // 3 warning
	viper.SetDefault("RedisAddress", "localhost:6379")
	viper.SetDefault("RedisPassword", "")
	viper.SetDefault("RedisDB", 0)
	c := &Config{}
	log.SetLevel(log.Level(c.Verbose))
	log.WithField("level", log.GetLevel().String()).Info("Setting log level to")

	// Log formating
	log.SetFormatter(
		&log.TextFormatter{
			DisableColors:   false,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		},
	)

	err := viper.Unmarshal(&c)
	return c, err
}
