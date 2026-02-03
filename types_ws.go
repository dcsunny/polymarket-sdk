// types_ws.go 模块
package polymarket

import (
	"encoding/json"
	"strconv"
)

// FlexInt can unmarshal from string or number.
type FlexInt int64

func (f *FlexInt) UnmarshalJSON(data []byte) error {
	var i int64
	if err := json.Unmarshal(data, &i); err == nil {
		*f = FlexInt(i)
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	parsed, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	*f = FlexInt(parsed)
	return nil
}

func (f FlexInt) Int64() int64 {
	return int64(f)
}

// WSSChannelType indicates channel type.
type WSSChannelType string

const (
	WSSChannelTypeUser   WSSChannelType = "user"
	WSSChannelTypeMarket WSSChannelType = "market"
)

// WSSAuth represents auth for user channel.
type WSSAuth struct {
	APIKey     string `json:"apiKey"`
	Secret     string `json:"secret"`
	Passphrase string `json:"passphrase"`
}

// WSSSubscription represents subscription request.
type WSSSubscription struct {
	Auth      *WSSAuth       `json:"auth,omitempty"`
	Type      WSSChannelType `json:"type"`
	Markets   []string       `json:"markets,omitempty"`
	AssetsIDs []string       `json:"assets_ids,omitempty"`
}

// WSSMessageHandler handles raw message.
type WSSMessageHandler func(data json.RawMessage) error

// WSSTradeEvent represents a trade event from user channel.
type WSSTradeEvent struct {
	EventType    string          `json:"event_type"`
	ID           string          `json:"id"`
	Market       string          `json:"market"`
	AssetID      string          `json:"asset_id"`
	Owner        string          `json:"owner"`
	Price        string          `json:"price"`
	Side         string          `json:"side"`
	Size         string          `json:"size"`
	Status       string          `json:"status"`
	Timestamp    FlexInt         `json:"timestamp"`
	TakerOrderID string          `json:"taker_order_id"`
	MakerOrders  []MakerOrder    `json:"maker_orders"`
	RawData      json.RawMessage `json:"-"`
}

// WSSOrderEvent represents an order event from user channel.
type WSSOrderEvent struct {
	EventType        string   `json:"event_type"`
	ID               string   `json:"id"`
	Market           string   `json:"market"`
	AssetID          string   `json:"asset_id"`
	OrderOwner       string   `json:"order_owner"`
	Price            string   `json:"price"`
	Side             string   `json:"side"`
	OriginalSize     string   `json:"original_size"`
	SizeMatched      string   `json:"size_matched"`
	Timestamp        FlexInt  `json:"timestamp"`
	Type             string   `json:"type"`
	AssociatedTrades []string `json:"associated_trades,omitempty"`
}

// WSSBookMessage represents order book snapshot.
type WSSBookMessage struct {
	EventType string            `json:"event_type"`
	Market    string            `json:"market"`
	AssetID   string            `json:"asset_id"`
	Timestamp FlexInt           `json:"timestamp"`
	Hash      string            `json:"hash"`
	Bids      []WSSOrderSummary `json:"bids"`
	Asks      []WSSOrderSummary `json:"asks"`
}

// WSSOrderSummary represents an order book level in WSS.
type WSSOrderSummary struct {
	Price string `json:"price"`
	Size  string `json:"size"`
}

// WSSPriceChangeMessage represents price change events.
type WSSPriceChangeMessage struct {
	EventType string             `json:"event_type"`
	Market    string             `json:"market"`
	AssetID   string             `json:"asset_id"`
	Timestamp FlexInt            `json:"timestamp"`
	Changes   []PriceLevelChange `json:"changes"`
}

// PriceLevelChange represents a single price level change.
type PriceLevelChange struct {
	Side    string `json:"side"`
	Price   string `json:"price"`
	Size    string `json:"size"`
	BestBid string `json:"best_bid"`
	BestAsk string `json:"best_ask"`
	Hash    string `json:"hash"`
}

// WSSTickSizeChangeMessage represents tick size change.
type WSSTickSizeChangeMessage struct {
	EventType        string  `json:"event_type"`
	Market           string  `json:"market"`
	AssetID          string  `json:"asset_id"`
	Timestamp        FlexInt `json:"timestamp"`
	PreviousTickSize string  `json:"previous_tick_size"`
	CurrentTickSize  string  `json:"current_tick_size"`
}

// WSSLastTradePriceMessage represents last trade price event.
type WSSLastTradePriceMessage struct {
	EventType string  `json:"event_type"`
	Market    string  `json:"market"`
	AssetID   string  `json:"asset_id"`
	Timestamp FlexInt `json:"timestamp"`
	Price     string  `json:"price"`
	Size      string  `json:"size"`
	Side      string  `json:"side"`
	FeeRate   string  `json:"fee_rate"`
}
