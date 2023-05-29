package design

import (
	. "goa.design/goa/v3/dsl"
)

var (
	_ = API("Test", func() {
		Title("BitCoin API")
		Description("An Api to receive email notifications about BitCoin prices")
		Version("1.0")
		HTTP(func() {
			Path("/api/v1")
		})
		Server("http", func() {
			Host("local", func() {
				URI("http://localhost")
			})
		})
	})
	CustomJSONResult = Type("CustomJSONResponse", func() {
		Attribute("content-type", String, func() {
			Default("application/json")
		})
		Attribute("content-length", Int)
	})
	CustomJSONResponse = func() {
		Header("content-type:Content-Type")
		Header("content-length:Content-Length")
	}
)
