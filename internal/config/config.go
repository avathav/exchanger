package config

import (
	"os"
)

const OpenexchangeApiKey = "OPENEXCHANGE_API_KEY"

type ErrMissingKey struct {
	Key string
}

func (e ErrMissingKey) Error() string {
	return "config: missing key: " + e.Key
}

type Config struct {
	OpenExchangeAPIKey string
}

func LoadConfig() (*Config, error) {
	apiKey := os.Getenv(OpenexchangeApiKey)
	if apiKey == "" {
		return nil, ErrMissingKey{OpenexchangeApiKey}
	}

	return &Config{
		OpenExchangeAPIKey: apiKey,
	}, nil
}
