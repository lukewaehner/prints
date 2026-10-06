package ws

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"time"

	"github.com/coder/websocket"
)

const (
	baseDelay = 1 * time.Second
	maxDelay  = 8 * time.Second
	healthy   = 10 * time.Second
)

type Config struct {
	URL       string
	OnConnect func(ctx context.Context, conn *websocket.Conn) error // send substribe msg
	Handle    func(msg []byte) error                                // Per message handler
}

func Run(ctx context.Context, cfg Config) error {
	attempt := 0
	for {
		err := dial(ctx, cfg)
		start := time.Now()
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if time.Since(start) > healthy {
			attempt = 0
		}
		delay := jitter(baseDelay, maxDelay, attempt)
		log.Printf("websocket connection closed, reconnecting: %v", err)
		attempt++

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func dial(ctx context.Context, cfg Config) error {
	dialCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(dialCtx, cfg.URL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to websocket: %w", err)
	}
	defer conn.CloseNow()

	if err := cfg.OnConnect(ctx, conn); err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}

	for {
		_, msg, err := conn.Read(ctx)
		if err != nil {
			return fmt.Errorf("read: %w", err)
		}

		if err := cfg.Handle(msg); err != nil {
			return fmt.Errorf("handle: %w", err)
		}
	}
}

func jitter(base, capa time.Duration, attempt int) time.Duration {
	temp := capa
	if attempt < 30 {
		if d := base << attempt; d > 0 && d < capa {
			temp = d
		}
	}
	half := temp / 2
	return half + rand.N(half+1)
}
