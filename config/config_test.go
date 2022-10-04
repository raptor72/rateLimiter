package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var limit = CoolDownTime{
	SecLimit: 60,
}

var LoginLimit = CountLimit{
	Count: 10,
}

var PasswordLimit = CountLimit{
	Count: 100,
}

var IPLimit = CountLimit{
	Count: 1000,
}

func getDefaultConfig() *Config {
	return &Config{
		Port:             8080,
		Verbose:          6,
		RedisAddress:     "localhost:6379",
		RedisPassword:    "",
		RedisDB:          0,
		LoginLimit:       LoginLimit,
		PasswordLimit:    PasswordLimit,
		IPLimit:          IPLimit,
		LoginCoolDown:    limit,
		PasswordCoolDown: limit,
		IPCoolDown:       limit,
		ConnectionString: "postgres://limiter:123456@127.0.0.1:15432/limitdb?sslmode=disable",
		MaxOpenConn:      5,
	}
}

func TestNewConfig(t *testing.T) {
	expectedConfig := getDefaultConfig()
	actualConfig, err := New()
	require.NoError(t, err)
	require.Equal(t, expectedConfig, actualConfig)
}
