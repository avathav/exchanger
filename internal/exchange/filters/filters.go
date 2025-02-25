package filters

type RatesByCurrencies struct {
	currencies []string
}

func NewRatesByCurrencies(currencies []string) RatesByCurrencies {
	return RatesByCurrencies{currencies: currencies}
}

func (r RatesByCurrencies) Filter(rates map[string]float64) map[string]float64 {
	filtered := make(map[string]float64)
	for _, curr := range r.currencies {
		if _, ok := rates[curr]; ok {
			filtered[curr] = rates[curr]
		}
	}

	return filtered
}
