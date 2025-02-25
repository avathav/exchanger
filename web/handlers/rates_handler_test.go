package handlers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"exchanger/internal/exchange"
	internalmock "exchanger/internal/mock"
	"exchanger/web/handlers"
)

func setupRouter(client exchange.Client) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/rates", handlers.GetRatesHandler(client))
	return r
}

func TestGetRatesHandler(t *testing.T) {
	type expected struct {
		status    int
		resultLen int
		emptyBody bool
	}

	type testCase struct {
		name        string
		requestFunc *http.Request
		mockF       func(client *internalmock.ExchangeClient)
		expected    expected
	}

	testCases := []testCase{
		{
			name:        "no currencies",
			requestFunc: httptest.NewRequest(http.MethodGet, "/rates", nil),
			mockF:       func(_ *internalmock.ExchangeClient) {},
			expected: expected{
				status:    http.StatusBadRequest,
				emptyBody: true,
			},
		},
		{
			name:        "one valid currency",
			requestFunc: httptest.NewRequest(http.MethodGet, "/rates?currencies=USD", nil),
			mockF:       func(_ *internalmock.ExchangeClient) {},
			expected: expected{
				status:    http.StatusBadRequest,
				emptyBody: true,
			},
		},
		{
			name:        "exchange service error",
			requestFunc: httptest.NewRequest(http.MethodGet, "/rates?currencies=USD,GBP", nil),
			mockF: func(client *internalmock.ExchangeClient) {
				client.
					On("GetFilteredRates", mock.Anything, mock.Anything).
					Return(map[string]float64{}, errors.New("exchange service error"))
			},
			expected: expected{
				status:    http.StatusBadRequest,
				emptyBody: true,
			},
		},
		{
			name:        "two valid currencies",
			requestFunc: httptest.NewRequest(http.MethodGet, "/rates?currencies=USD,GBP", nil),
			mockF: func(client *internalmock.ExchangeClient) {
				client.
					On("GetFilteredRates", mock.Anything, mock.Anything).
					Return(map[string]float64{"USD": 1.0, "GBP": 0.8}, nil)
			},
			expected: expected{
				status:    http.StatusOK,
				resultLen: 2,
			},
		},
		{
			name:        "only one valid currency",
			requestFunc: httptest.NewRequest(http.MethodGet, "/rates?currencies=USD,unknown", nil),
			mockF: func(client *internalmock.ExchangeClient) {
				client.
					On("GetFilteredRates", mock.Anything, mock.Anything).
					Return(map[string]float64{"USD": 1.0}, nil)
			},
			expected: expected{
				status:    http.StatusBadRequest,
				emptyBody: true,
			},
		},
		{
			name:        "more than to valid currencies",
			requestFunc: httptest.NewRequest(http.MethodGet, "/rates?currencies=USD,GBP,EUR,unknown", nil),
			mockF: func(client *internalmock.ExchangeClient) {
				client.
					On("GetFilteredRates", mock.Anything, mock.Anything).
					Return(map[string]float64{"USD": 1.0, "GBP": 0.8, "EUR": 1.2}, nil)
			},
			expected: expected{
				status:    http.StatusOK,
				resultLen: 2,
			},
		},
		{
			name:        "call with missing currency",
			requestFunc: httptest.NewRequest(http.MethodGet, "/rates?currencies=USD,GBP,EUR", nil),
			mockF: func(client *internalmock.ExchangeClient) {
				client.
					On("GetFilteredRates", mock.Anything, mock.Anything).
					Return(map[string]float64{"USD": 1.0, "EUR": 1.2}, nil)
			},
			expected: expected{
				status:    http.StatusOK,
				resultLen: 2,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := new(internalmock.ExchangeClient)
			tc.mockF(mockClient)

			r := setupRouter(mockClient)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, tc.requestFunc)

			assert.Equal(t, tc.expected.status, w.Code, "invalid status code")

			if tc.expected.emptyBody {
				assert.Empty(t, w.Body.String(), "expected empty body")
			}

			if !tc.expected.emptyBody && tc.expected.status == http.StatusOK {
				var result []handlers.RateResponse
				err := json.Unmarshal(w.Body.Bytes(), &result)
				assert.NoError(t, err, "")

				assert.True(t, len(result) >= tc.expected.resultLen,
					"expected at least %d results, got %d",
					tc.expected.resultLen, len(result))
			}

			mockClient.AssertExpectations(t)
		})
	}
}
