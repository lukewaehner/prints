package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lukewaehner/prints/internal/coinbase"
	"github.com/lukewaehner/prints/internal/trade"
)

func main() {
	cb := coinbase.Connector{Products: []string{"BTC-USD", "ETH-USD", "SOL-USD"}}
	out := make(chan trade.Trade)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var runErr error
	go func() {
		e := cb.Run(ctx, out)
		if e != nil {
			runErr = e
		}
		close(out)
	}()

	for t := range out {
		fmt.Println(t.String())
	}

	if runErr != nil {
		if errors.Is(runErr, context.Canceled) {
			return
		}
		log.Fatalf("Error running connector: %v\n", runErr)
	}
}
