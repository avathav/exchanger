package cryptocurrency

import "errors"

var ErrCurrencyNotFound = errors.New("currency not found")

type Finder interface {
	FindCurrency(name string) (Currency, error)
}

type Repository interface {
	Finder
}

type Currency struct {
	Name          string
	DecimalPlaces int
	RateToUSD     float64
}

type Storage struct {
	data map[string]Currency
}

func NewStorage() *Storage {
	return &Storage{
		data: map[string]Currency{
			"BEER":  {DecimalPlaces: 18, RateToUSD: 0.00002461},
			"FLOKI": {DecimalPlaces: 18, RateToUSD: 0.0001428},
			"GATE":  {DecimalPlaces: 18, RateToUSD: 6.87},
			"USDT":  {DecimalPlaces: 6, RateToUSD: 0.999},
			"WBTC":  {DecimalPlaces: 8, RateToUSD: 57037.22},
		},
	}
}

func (s *Storage) FindCurrency(name string) (Currency, error) {
	if c, ok := s.data[name]; ok {
		return c, nil
	}
	return Currency{}, ErrCurrencyNotFound

}
