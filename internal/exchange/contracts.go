package exchange

import "context"

type Client interface {
	GetRates(ctx context.Context) (*ApiResponse, error)
	GetFilteredRates(ctx context.Context, filters ...RatesFilter) (map[string]float64, error)
}

type RatesFilter interface {
	Filter(currencies map[string]float64) map[string]float64
}
