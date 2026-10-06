// Package trade provides a struct for a single trade in a bucket of trades from a tick
package trade

import (
	"fmt"
	"time"
)

type Side string

const (
	Buy  Side = "buy"
	Sell Side = "sell"
)

// Trade - single trade event in a bucket of trades from a tick
type Trade struct {
	Exchange  string    `json:"exchange"`
	Ticker    string    `json:"ticker"`
	TradeID   string    `json:"trade_id"`
	Price     string    `json:"price"`
	Size      string    `json:"size"`
	TimeStamp time.Time `json:"timestamp"`
	Side      Side      `json:"side"`
}

func (t Trade) String() string {
	return fmt.Sprintf("[%s]\nTicker: %s, TradeID: %s\nPrice: %s, Size: %s, TimeStamp: %s, Side: %s",
		t.Exchange, t.Ticker, t.TradeID, t.Price, t.Size, t.TimeStamp.Format(time.RFC3339Nano), t.Side)
}
