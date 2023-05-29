package design

import (
	. "goa.design/goa/v3/dsl"
	cors "goa.design/plugins/v3/cors/dsl"
)

var (
	_ = Service("rate", func() {
		Description("The rate service provides BTC to UAH exchange rate")
		cors.Origin("*", func() {
			cors.Methods("*")
			cors.Headers("*")
		})
		HTTP(func() {
			Path("/rate")
		})

		Method("rate", func() {
			Description("This request should return the current BTC to UAH exchange rate using any third-party service with a public API")
			HTTP(func() {
				GET("/")
				Response(StatusOK, CustomJSONResponse)
				SkipResponseBodyEncodeDecode()
			})
			Result(CustomJSONResult)
		})
	})
)
