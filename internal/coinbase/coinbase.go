package coinbase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/coder/websocket"
	"github.com/lukewaehner/prints/internal/trade"
	"github.com/lukewaehner/prints/internal/ws"
)

type Connector struct {
	Products []string
}

type SubscribeMessage struct {
	Type       string   `json:"type"`
	ProductIDs []string `json:"product_ids"`
	Channel    string   `json:"channel"`
}

func (c *Connector) Run(ctx context.Context, out chan<- trade.Trade) error {
	return ws.Run(ctx, ws.Config{
		URL:       "wss://advanced-trade-ws.coinbase.com",
		OnConnect: c.subscribe,
		Handle: func(msg []byte) error {
			trades, err := parse(msg)
			if err != nil {
				return fmt.Errorf("unable to parse message: %w", err)
			}
			for _, t := range trades {
				select {
				case out <- t:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			return nil
		},
	})
}

func (c *Connector) subscribe(ctx context.Context, conn *websocket.Conn) error {
	subscribeMessage := SubscribeMessage{
		Type:       "subscribe",
		ProductIDs: c.Products,
		Channel:    "market_trades",
	}

	subscription, err := json.Marshal(subscribeMessage)
	if err != nil {
		return fmt.Errorf("unable to marshal subscriptionMessage: %w", err)
	}
	err = conn.Write(ctx, websocket.MessageText, subscription)
	if err != nil {
		return fmt.Errorf("could not subscribe to feed: %w", err)
	}
	return nil
}
