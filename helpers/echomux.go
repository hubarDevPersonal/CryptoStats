package helpers

import (
	"CryptoStats/api/context"
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
)

var paramRegex = regexp.MustCompile("{(.*?)}")

type EchoMux struct {
	e *echo.Echo
}

func NewEchoMux(engine *echo.Echo) *EchoMux {
	return &EchoMux{e: engine}
}

func (mux *EchoMux) Handle(method, pattern string, handler http.HandlerFunc) {
	echoRoute := mux.convertToEchoRoute(pattern)
	mux.e.Add(method, echoRoute, func(c echo.Context) error {
		req := mux.populateContext(c)
		handler(c.Response(), req)
		return nil
	})
}

// ServeHTTP dispatches the request to the handler whose method
// matches the request method and whose pattern most closely
// matches the request URL.
func (mux *EchoMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux.e.ServeHTTP(w, r)
}

func (mux *EchoMux) Vars(req *http.Request) map[string]string {
	ctx := req.Context()
	if p, ok := ctx.Value(context.ParamKey).(map[string]string); ok {
		return p
	}
	return map[string]string{}
}

func (mux *EchoMux) convertToEchoRoute(pattern string) string {
	return paramRegex.ReplaceAllString(pattern, ":$1")
}

func (mux *EchoMux) populateContext(c echo.Context) *http.Request {
	req := c.Request()
	ctx := req.Context()
	return req.WithContext(ctx)
}
