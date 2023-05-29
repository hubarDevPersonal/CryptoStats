package context

import (
	"CryptoStats/log"
	"context"
	"os"
	"os/signal"
	"syscall"
)

func ConsoleCancelableContext(log *log.Logger) context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalCh
		if log != nil {
			log.Info("Ctrl-C hit")
		}
		cancel()
	}()

	return ctx
}
