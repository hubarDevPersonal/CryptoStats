package service

import (
	"CryptoStats/log"
	"encoding/json"
	"io"
	"net/http"
)

type Price struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type RateService interface {
	GetPrice(currency string) (*Price, error)
}

type rateService struct {
	l *log.Logger
}

func NewRateService(l *log.Logger) RateService {
	return &rateService{
		l: l,
	}
}
func (s *rateService) GetPrice(symbol string) (*Price, error) {
	client := http.DefaultClient
	req, err := http.NewRequest(http.MethodGet, "https://api.binance.com/api/v3/ticker/price?symbol="+symbol, nil)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	bytes, err := bodyBytes(response.Body)
	if err != nil {
		return nil, err
	}
	price, err := priceResponse(bytes)
	return price, nil
}

func bodyBytes(body io.ReadCloser) ([]byte, error) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(body)
	bytes, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func priceResponse(bytes []byte) (*Price, error) {
	var response *Price
	err := json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
