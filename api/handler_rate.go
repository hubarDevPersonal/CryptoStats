package api

import (
	"CryptoStats/api/gen/rate"
	"CryptoStats/log"
	"CryptoStats/service"
	"context"
	"io"
)

var _ rate.Service = (*handlerRate)(nil)

type handlerRate struct {
	l    *log.Logger
	rate service.RateService
}

func NewHandlerRate(
	l *log.Logger,
) *handlerRate {
	return &handlerRate{
		l:    l,
		rate: service.NewRateService(l),
	}
}

func (h *handlerRate) Rate(ctx context.Context) (*rate.CustomJSONResponse, io.ReadCloser, error) {
	res, err := h.rate.GetPrice("BTCUAH")
	if err != nil {
		h.l.Error("error getting rate", log.Err(err))
		return nil, nil, err
	}
	body, err := jsonResponse(h.l, "rate", res)
	return &rate.CustomJSONResponse{
		ContentType: "application/json",
	}, body, nil
}
