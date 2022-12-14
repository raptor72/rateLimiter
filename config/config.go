package config

import (
	"errors"
	"time"

	_ "github.com/jackc/pgx/stdlib" // trust
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type CoolDownTime struct {
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
	IPLimit          CountLimit
	LoginCoolDown    CoolDownTime
	PasswordCoolDown CoolDownTime
	IPCoolDown       CoolDownTime
	ConnectionString string
	MaxOpenConn      int
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
	limit := CoolDownTime{
		SecLimit: 60,
	}

	LoginLimit := CountLimit{
		Count: 10,
	}

	PasswordLimit := CountLimit{
		Count: 100,
	}

	IPLimit := CountLimit{
		Count: 1000,
	}

	viper.SetDefault("Port", 8080)
	viper.SetDefault("Verbose", 6) // 6 Debug // 4 info // 3 warning
	viper.SetDefault("RedisAddress", "localhost:6379")
	viper.SetDefault("RedisPassword", "")
	viper.SetDefault("RedisDB", 0)
	viper.SetDefault("LoginLimit", LoginLimit)
	viper.SetDefault("PasswordLimit", PasswordLimit)
	viper.SetDefault("IPLimit", IPLimit)
	viper.SetDefault("LoginCoolDown", limit)
	viper.SetDefault("PasswordCoolDown", limit)
	viper.SetDefault("IPCoolDown", limit)
	viper.SetDefault("ConnectionString", "postgres://limiter:123456@127.0.0.1:15432/limitdb?sslmode=disable")
	viper.SetDefault("MaxOpenConn", 5)

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

func (cfg *Config) NewDB() (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", cfg.ConnectionString)
	if err == nil {
		db.SetMaxOpenConns(cfg.MaxOpenConn)
		db.SetConnMaxIdleTime(5 * time.Minute)
		db.SetConnMaxLifetime(30 * time.Minute)
		return db, nil
	}
	return nil, errors.New("could not connect to database")
}
