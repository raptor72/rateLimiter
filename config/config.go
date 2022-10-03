package config

import (
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

type CoulDownTime struct {
	SecLimit int
}

type CountLimit struct {
	Count int
}

type Config struct {
	Port             int
	Verbose          int
	RedisAddress     string
	RedisPassword    string
	RedisDB          int
	LoginLimit       CountLimit
	PasswordLimit    CountLimit
	IpLimit          CountLimit
	LoginCouldown    CoulDownTime
	PasswordCouldown CoulDownTime
	IpCouldown       CoulDownTime
	connectionString string
	maxOpenConn      int
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
	limit := CoulDownTime{
		SecLimit: 60,
	}

	LoginLimit := CountLimit{
		Count: 10,
	}

	PasswordLimit := CountLimit{
		Count: 100,
	}

	IpLimit := CountLimit{
		Count: 1000,
	}

	viper.SetDefault("Port", 8080)
	viper.SetDefault("Verbose", 6) // 6 Debug // 4 info // 3 warning
	viper.SetDefault("RedisAddress", "localhost:6379")
	viper.SetDefault("RedisPassword", "")
	viper.SetDefault("RedisDB", 0)
	viper.SetDefault("LoginLimit", LoginLimit)
	viper.SetDefault("PasswordLimit", PasswordLimit)
	viper.SetDefault("IpLimit", IpLimit)
	viper.SetDefault("LoginCouldown", limit)
	viper.SetDefault("PasswordCouldown", limit)
	viper.SetDefault("IpCouldown", limit)
	viper.SetDefault("connectionString", "postgresql://limiter:123456@127.0.0.1:15432/limitdb")
	viper.SetDefault("maxOpenConn", 5)

	cfg, err := InitConfig()
	if err != nil {
		log.WithError(err).Error("Could not init Config")
		return nil, err
	}

	log.SetLevel(log.Level(cfg.Verbose))
	log.WithField("level", log.GetLevel().String()).Info("Setting log level to")

	// Log formating
	log.SetFormatter(
		&log.TextFormatter{
			DisableColors:   false,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		},
	)

	return cfg, err
}

func InitConfig() (*Config, error) {
	c := &Config{}
	err := viper.Unmarshal(&c)
	return c, err
}

func NewDB(cfg *Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", cfg.connectionString)
	if err == nil {
		db.SetMaxOpenConns(cfg.maxOpenConn)
		db.SetConnMaxIdleTime(5 * time.Minute)
		db.SetConnMaxLifetime(30 * time.Minute)
		return db, nil
	}
	fmt.Println(err)
	return nil, errors.New("could not connect to database")
}
