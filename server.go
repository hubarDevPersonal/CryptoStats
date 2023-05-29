package main

import (
	"CryptoStats/api"
	"CryptoStats/config"
	"CryptoStats/helpers"
	"CryptoStats/log"
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	goaMiddleware "goa.design/goa/v3/middleware"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	log     *log.Logger
	httpSvr *echo.Echo
}

func NewServer(l *log.Logger, cfg *config.Config) (*Server, error) {

	middleware.DefaultLoggerConfig.Format = "${time_custom}    W3C     ${remote_ip} ${host} ${method} ${uri} ${user_agent} ${status} ${error} ${latency_human} <-${bytes_in} ->${bytes_out}\n"
	middleware.DefaultLoggerConfig.CustomTimeFormat = "2006-01-02T15:04:05.000Z0700"
	httpSvr := echo.New()
	httpSvr.HideBanner = true
	httpSvr.HidePort = true

	httpSvr.Pre(
		encodeSemicolonsInQueryString,
	)
	httpSvr.Use(
		middleware.Logger(),
		middleware.Recover(),
		middleware.Gzip(),
	)

	httpSvr.GET("/healthCheck", func(c echo.Context) error {
		w := c.Response()
		_, _ = fmt.Fprintln(w, "Hello")
		return nil
	})

	errHandler := func(ctx context.Context, w http.ResponseWriter, err error) {
		id, ok := ctx.Value(goaMiddleware.RequestIDKey).(string)
		if ok {
			_, _ = w.Write([]byte("[" + id + "] encoding: " + err.Error()))
		}
	}

	api, err := api.New(
		l,
		cfg,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating api handler: %w", err)
	}

	api.MountREST(httpSvr, errHandler)
	helpers.NewStaticFiles(http.Dir("api")).AttachTo(httpSvr, "/api/v1")

	return &Server{
		log:     l,
		httpSvr: httpSvr,
	}, nil
}

func (svr *Server) Run(ctx context.Context, httpPort string) error {
	grp, ctx := errgroup.WithContext(ctx)

	// Start REST server
	grp.Go(func() error {
		svr.log.Info("http listening", log.String("port", httpPort))
		svr.logRoutes(svr.httpSvr)

		return svr.httpSvr.Start(":" + httpPort)
	})

	// Wait to shut down REST server
	grp.Go(func() error {
		<-ctx.Done()
		svr.log.Info("shutting down REST server")
		shutDownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := svr.httpSvr.Shutdown(shutDownCtx)
		if err != nil && err.Error() != "http: Server closed" {
			return err
		}

		return nil
	})

	return grp.Wait()
}

func (svr *Server) logRoutes(e *echo.Echo) {
	routes := e.Routes()
	sort.Slice(routes, func(i, j int) bool {
		if routes[i].Path == routes[j].Path {
			return routes[i].Method < routes[j].Method
		}

		return routes[i].Path < routes[j].Path
	})
	for _, route := range routes {
		if strings.Contains(route.Name, "NotFoundHandler") {
			continue
		}
		svr.log.Info("route",
			log.String("method", route.Method),
			log.String("path", route.Path),
			log.String("name", route.Name),
		)
	}
}

func encodeSemicolonsInQueryString(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		req.URL.RawQuery = strings.Replace(req.URL.RawQuery, ";", "%3B", -1)

		return next(c)
	}
}
