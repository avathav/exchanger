package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"exchanger/internal/exchange"
	"exchanger/internal/exchange/filters"
)

type RateResponse struct {
	From string  `json:"from"`
	To   string  `json:"to"`
	Rate float64 `json:"rate"`
}

func GetRatesHandler(client exchange.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		currencies, cErr := fetchCurrenciesParam(c.Query("currencies"))
		if cErr != nil {
			log.Warningln(cErr)
			c.Status(http.StatusBadRequest)
			return
		}

		rates, exchangeErr := client.GetFilteredRates(c, filters.NewRatesByCurrencies(currencies))
		if exchangeErr != nil {
			log.Warningln(exchangeErr)
			c.Status(http.StatusBadRequest)
			return
		}

		if len(rates) < 2 {
			log.Warningln("at least two currencies are required")
			c.Status(http.StatusBadRequest)
			return
		}

		var response []RateResponse
		for i := 0; i < len(currencies); i++ {
			if _, ok := rates[currencies[i]]; !ok {
				log.Warningf("currency %s not found in rates", currencies[i])
				continue
			}

			for j := i + 1; j < len(currencies); j++ {
				if _, ok := rates[currencies[j]]; !ok {
					log.Warningf("currency %s not found in rates", currencies[j])
					continue
				}

				response = append(response, CalculateRate(currencies[i], currencies[j], rates))
				response = append(response, CalculateRate(currencies[j], currencies[i], rates))
			}
		}

		c.JSON(http.StatusOK, response)
	}
}

func CalculateRate(from, to string, rates map[string]float64) RateResponse {
	return RateResponse{
		From: from,
		To:   to,
		Rate: rates[to] / rates[from],
	}
}

func fetchCurrenciesParam(currencies string) ([]string, error) {
	list := strings.Split(currencies, ",")
	if len(list) < 2 {
		return nil, fmt.Errorf("at least two currencies are required")
	}

	return list, nil
}
