package di

import (
	"net/http"
	"time"

	"exchanger/internal/config"
	"exchanger/internal/exchange"
	"exchanger/internal/repository/cryptocurrency"
)

type Container struct {
	conf    *config.Config
	storage *cryptocurrency.Storage
}

var container = &Container{
	conf:    nil,
	storage: nil,
}

func LoadConfig() (*config.Config, error) {
	if container.conf != nil {
		return container.conf, nil
	}

	c, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	container.conf = c

	return c, nil
}

func DefaultHttpClient() *http.Client {
	return &http.Client{Timeout: 30 * time.Second}
}

func DefaultExchangeClient() (*exchange.OpenExchangeRatesClient, error) {
	c, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	return exchange.NewOpenExchangeRatesClient(c.OpenExchangeAPIKey, DefaultHttpClient())
}

func DefaultCryptoCurrencyStorage() *cryptocurrency.Storage {
	if container.storage != nil {
		return container.storage
	}

	container.storage = cryptocurrency.NewStorage()

	return container.storage
}
