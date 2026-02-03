// types_rtds.go 模块
package polymarket

import "encoding/json"

// RTDSSubscription represents subscription request.
type RTDSSubscription struct {
	Action        string                   `json:"action"`
	Subscriptions []RTDSSubscriptionDetail `json:"subscriptions"`
}

// RTDSSubscriptionDetail represents a single subscription detail.
type RTDSSubscriptionDetail struct {
	Topic   string      `json:"topic"`
	Type    string      `json:"type"`
	Filters interface{} `json:"filters,omitempty"`
}

// RTDSMessage represents a message received from RTDS.
type RTDSMessage struct {
	Topic     string          `json:"topic"`
	Type      string          `json:"type"`
	Timestamp int64           `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"`
}

// RTDSCryptoPricePayload represents crypto price data payload.
type RTDSCryptoPricePayload struct {
	Symbol    string  `json:"symbol"`
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

// CryptoPriceSource represents the source of crypto price data.
type CryptoPriceSource string

const (
	CryptoPriceSourceBinance   CryptoPriceSource = "crypto_prices"
	CryptoPriceSourceChainlink CryptoPriceSource = "crypto_prices_chainlink"
)

// RTDSMessageHandler handles RTDS messages.
type RTDSMessageHandler func(msg *RTDSMessage) error
