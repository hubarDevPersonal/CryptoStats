package api

import (
	"CryptoStats/config"
	"CryptoStats/log"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "goa.design/goa/v3/pkg"
)

type (
	Handler struct {
	}
)

func New(
	l *log.Logger,
	cfg *config.Config,
) (*Handler, error) {
	h := Handler{}
	return &h, nil
}

func (h *Handler) MountREST(e *echo.Echo, errHandler func(ctx context.Context, w http.ResponseWriter, err error)) {
	//dec := goahttp.RequestDecoder
	//enc := goahttp.ResponseEncoder
	//mux := helpers.NewEchoMux(e)

}

func jsonResponse(l *log.Logger, label string, obj interface{}) (io.ReadCloser, error) {
	pr, pw := io.Pipe()
	enc := json.NewEncoder(pw)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	go func() {
		err := enc.Encode(obj)
		if err != nil {
			l.Error(fmt.Sprintf("error encoding %s", label), log.Err(err))
		}
		err = pw.Close()
		if err != nil {
			l.Error("error closing json encoder pipe", log.Err(err))
		}
	}()

	return pr, nil
}
