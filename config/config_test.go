package config

import (
	"github.com/stretchr/testify/require"
	"testing"
)

var limit = CoulDownTime{
	SecLimit: 60,
} 

var LoginLimit = CountLimit{
	Count: 10,
}

var PasswordLimit = CountLimit{
	Count: 100,
}

var IpLimit = CountLimit{
	Count: 1000,
}

func getDefaultConfig() *Config {
	return &Config{
		Port:          8080,
		Verbose:       6,
		RedisAddress:  "localhost:6379",
		RedisPassword: "",
		RedisDB:       0,
		LoginLimit:    LoginLimit,
		PasswordLimit: PasswordLimit,
		IpLimit:       IpLimit,
		LoginCouldown: limit,
		PasswordCouldown: limit,
		IpCouldown: limit,
	}
}

func TestNewConfig(t *testing.T) {
	expectedConfig := getDefaultConfig()
	actualConfig, err := New()
	require.NoError(t, err)
	require.Equal(t, expectedConfig, actualConfig)
}
