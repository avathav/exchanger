package main

import (
	"log"

	"exchanger/internal/di"
	"exchanger/web"
)

func main() {
	exchangeClient, clientErr := di.DefaultExchangeClient()
	if clientErr != nil {
		log.Fatalf("failed to create exchange client: %v", clientErr)
	}

	s := web.NewServer(exchangeClient, di.DefaultCryptoCurrencyStorage())
	s.Routes()

	if err := s.Run(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
