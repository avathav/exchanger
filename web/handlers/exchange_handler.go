package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"exchanger/internal/repository/cryptocurrency"
)

type ExchangeCryptoResponse struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

func GetExchangeHandler(finder cryptocurrency.Finder) gin.HandlerFunc {
	return func(c *gin.Context) {
		from := c.Query("from")
		to := c.Query("to")
		amountStr := c.Query("amount")

		if from == "" || to == "" || amountStr == "" {
			log.Error("missing required parameters")
			c.Status(http.StatusBadRequest)
			return
		}

		fromCrypto, fromErr := finder.FindCurrency(from)
		if fromErr != nil {
			log.Error(fromErr)
			c.Status(http.StatusBadRequest)
			return
		}

		toCrypto, toErr := finder.FindCurrency(to)
		if toErr != nil {
			log.Error(toErr)
			c.Status(http.StatusBadRequest)
			return
		}

		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		usdValue := amount * fromCrypto.RateToUSD
		toValue := usdValue / toCrypto.RateToUSD
		factor := math.Pow10(toCrypto.DecimalPlaces)

		response := ExchangeCryptoResponse{
			From:   from,
			To:     to,
			Amount: math.Round(toValue*factor) / factor,
		}

		c.JSON(http.StatusOK, response)
	}
}
