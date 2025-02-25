package mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"exchanger/internal/exchange"
)

type ExchangeClient struct {
	mock.Mock
}

func (m *ExchangeClient) GetRates(ctx context.Context) (*exchange.ApiResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).(*exchange.ApiResponse), args.Error(1)
}

func (m *ExchangeClient) GetFilteredRates(ctx context.Context, f ...exchange.RatesFilter) (map[string]float64, error) {
	args := m.Called(ctx, f)
	return args.Get(0).(map[string]float64), args.Error(1)
}
