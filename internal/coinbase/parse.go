package coinbase

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/lukewaehner/prints/internal/trade"
)

type marketTradeMsg struct {
	Channel     string    `json:"channel"`
	Timestamp   time.Time `json:"timestamp"`
	SequenceNum int64     `json:"sequence_num"`
	Events      []event   `json:"events"`
}

type event struct {
	Type   string    `json:"type"`
	Trades []cbTrade `json:"trades"`
}

type cbTrade struct {
	ProductID string    `json:"product_id"`
	TradeID   string    `json:"trade_id"`
	Price     string    `json:"price"`
	Size      string    `json:"size"`
	Time      time.Time `json:"time"`
	Side      string    `json:"side"`
}

func parse(msg []byte) ([]trade.Trade, error) {
	marketTrade := marketTradeMsg{}
	err := json.Unmarshal(msg, &marketTrade)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal %w", err)
	}

	if marketTrade.Channel != "market_trades" {
		return nil, nil
	}

	trades := []trade.Trade{}

	for _, ev := range marketTrade.Events {
		// Only proc update events, eventually include snapshots though
		// TODO: add snapshot handling
		if ev.Type != "update" {
			continue
		}

		for _, cbt := range ev.Trades {
			s, err := toSide(cbt.Side)
			if err != nil {
				log.Printf("Error converting side: %v", err)
				continue
			}
			trd := trade.Trade{
				Exchange:  "coinbase",
				Ticker:    cbt.ProductID,
				TradeID:   cbt.TradeID,
				Price:     cbt.Price,
				Size:      cbt.Size,
				TimeStamp: cbt.Time,
				Side:      s,
			}
			trades = append(trades, trd)
		}
	}
	return trades, nil
}

func toSide(side string) (trade.Side, error) {
	switch strings.ToLower(side) {
	case "buy":
		return trade.Buy, nil
	case "sell":
		return trade.Sell, nil
	}
	return "", fmt.Errorf("unknown side %q", side)
}
