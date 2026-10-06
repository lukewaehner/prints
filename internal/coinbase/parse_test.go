package coinbase

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/lukewaehner/prints/internal/trade"
)

func TestParse_update_with_two_trades_returns_both_normalized(t *testing.T) {
	msg := []byte(`{"channel":"market_trades","timestamp":"2026-09-29T21:02:15.784787554Z","sequence_num":34,"events":[{"type":"update","trades":[
		{"product_id":"BTC-USD","trade_id":"1100160197","price":"83561.13","size":"0.00029619","time":"2026-09-29T21:02:15.751986Z","side":"SELL"},
		{"product_id":"BTC-USD","trade_id":"1100160196","price":"83561.12","size":"0.0000156","time":"2026-09-29T21:02:15.679346Z","side":"BUY"}]}]}`)
	want := []trade.Trade{
		{
			Exchange:  "coinbase",
			Ticker:    "BTC-USD",
			TradeID:   "1100160197",
			Price:     "83561.13",
			Size:      "0.00029619",
			TimeStamp: time.Date(2026, 9, 29, 21, 2, 15, 751986000, time.UTC),
			Side:      trade.Sell,
		},
		{
			Exchange:  "coinbase",
			Ticker:    "BTC-USD",
			TradeID:   "1100160196",
			Price:     "83561.12",
			Size:      "0.0000156",
			TimeStamp: time.Date(2026, 9, 29, 21, 2, 15, 679346000, time.UTC),
			Side:      trade.Buy,
		},
	}

	got, err := parse(msg)
	if err != nil {
		t.Fatalf("parse() error = %v, want nil", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("parse() =\n%+v\nwant\n%+v", got, want)
	}
}

func TestParse_non_trade_messages_return_no_trades_and_no_error(t *testing.T) {
	tests := []struct {
		name string
		msg  string
	}{
		{
			name: "subscriptions confirmation",
			msg:  `{"channel":"subscriptions","timestamp":"2026-09-29T21:02:14Z","sequence_num":0,"events":[{"subscriptions":{"market_trades":["BTC-USD"]}}]}`,
		},
		{
			name: "heartbeat",
			msg:  `{"channel":"heartbeats","timestamp":"2026-09-29T21:02:14Z","sequence_num":1,"events":[{"current_time":"2026-09-29T21:02:14Z","heartbeat_counter":1}]}`,
		},
		{
			name: "snapshot event is skipped",
			msg:  `{"channel":"market_trades","timestamp":"2026-09-29T21:02:14Z","sequence_num":2,"events":[{"type":"snapshot","trades":[{"product_id":"BTC-USD","trade_id":"1","price":"1","size":"1","time":"2026-09-29T21:02:14Z","side":"BUY"}]}]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// act
			got, err := parse([]byte(tt.msg))
			// assert
			if err != nil {
				t.Fatalf("parse() error = %v, want nil", err)
			}
			if len(got) != 0 {
				t.Errorf("parse() returned %d trades, want 0: %+v", len(got), got)
			}
		})
	}
}

func TestParse_unknown_side_skips_that_trade_and_keeps_the_rest(t *testing.T) {
	msg := []byte(`{"channel":"market_trades","timestamp":"2026-09-29T21:02:15Z","sequence_num":5,"events":[{"type":"update","trades":[
		{"product_id":"ETH-USD","trade_id":"10","price":"2500.00","size":"1","time":"2026-09-29T21:02:15Z","side":"HOLD"},
		{"product_id":"ETH-USD","trade_id":"11","price":"2500.01","size":"2","time":"2026-09-29T21:02:15Z","side":"SELL"}]}]}`)

	got, err := parse(msg)
	if err != nil {
		t.Fatalf("parse() error = %v, want nil", err)
	}
	if len(got) != 1 || got[0].TradeID != "11" {
		t.Errorf("parse() = %+v, want only trade 11", got)
	}
}

func TestParse_invalid_json_returns_wrapped_syntax_error(t *testing.T) {
	msg := []byte(`{"channel":`)

	_, err := parse(msg)

	var syntaxErr *json.SyntaxError
	if !errors.As(err, &syntaxErr) {
		t.Errorf("parse() error = %v, want it to wrap *json.SyntaxError", err)
	}
}
