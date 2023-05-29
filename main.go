package main

import (
	"CryptoStats/api/context"
	"CryptoStats/config"
	"CryptoStats/log"
	"strings"
)

func main() {
	var l *log.Logger
	l, _ = log.New("dev", "debug")
	cfg, err := config.New()
	if err != nil {
		l.Fatal("error getting config", log.Err(err))
	}
	svr, err := NewServer(l, cfg)
	if err != nil {
		l.Fatal("error creating server", log.Err(err))
	}
	ctx := context.ConsoleCancelableContext(l)
	err = svr.Run(ctx, "80")

	if err != nil && !strings.Contains(err.Error(), "Server closed") {
		l.Fatal("error running server", log.Err(err))
	}

	l.Info(`Server closed cleanly. \o/`)
}
