package exchange

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const BaseUrl = "https://openexchangerates.org/api/"

var (
	ErrOpenExchangeConnectionFailed = fmt.Errorf("failed to connect to openexchangerates")
)

type OpenExchangeRatesClient struct {
	apiKey     string
	httpClient *http.Client
}

func NewOpenExchangeRatesClient(apiKey string, client *http.Client) (*OpenExchangeRatesClient, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("apiKey cannot be empty")
	}

	if client == nil {
		return nil, fmt.Errorf("httpClient cannot be nil")
	}

	return &OpenExchangeRatesClient{
		apiKey:     apiKey,
		httpClient: client,
	}, nil
}

func (s *OpenExchangeRatesClient) GetRates(ctx context.Context) (*ApiResponse, error) {
	req, rErr := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%slatest.json?app_id=%s", BaseUrl, s.apiKey), nil)
	if rErr != nil {
		return nil, fmt.Errorf("failed to create request: %w", rErr)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Errorf("error closing response body: %v", closeErr)
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		logErrorResponse(resp)

		return nil, ErrOpenExchangeConnectionFailed
	}

	data := &ApiResponse{}
	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *OpenExchangeRatesClient) GetFilteredRates(ctx context.Context, filters ...RatesFilter) (map[string]float64, error) {
	rates, err := s.GetRates(ctx)
	if err != nil {
		return nil, err
	}

	for _, f := range filters {
		rates.Rates = f.Filter(rates.Rates)
	}

	return rates.Rates, nil
}

func logErrorResponse(resp *http.Response) {
	var data ResponseError
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Errorf("failed to decode response body: %v", err)
	}

	log.WithFields(log.Fields{
		"status":      data.Code,
		"message":     data.Message,
		"description": data.Description,
	}).Error("failed to get exchange rates")
}
