package config

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func getDefaultConfig() *Config {
	return &Config{
		Port:          8080,
		Verbose:       6,
		RedisAddress:  "localhost:6379",
		RedisPassword: "",
		RedisDB:       0,
		LoginLimit:    10,
		PasswordLimit: 100,
		IpLimit:       1000,
	}
}

func TestNewConfig(t *testing.T) {
	expectedConfig := getDefaultConfig()
	actualConfig, err := New()
	require.NoError(t, err)
	require.Equal(t, expectedConfig, actualConfig)
}
