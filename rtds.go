// rtds.go 模块
package polymarket

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	rtdsPingInterval = 5 * time.Second
)

// RTDSClient 处理 Polymarket RTDS 流式传输。
type RTDSClient struct {
	cfg Config

	mu       sync.RWMutex
	conn     *websocket.Conn
	handlers map[string]RTDSMessageHandler

	ctx    context.Context
	cancel context.CancelFunc
}

func NewRTDSClient(cfg Config) *RTDSClient {
	ctx, cancel := context.WithCancel(context.Background())
	return &RTDSClient{
		cfg:      cfg,
		handlers: make(map[string]RTDSMessageHandler),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Connect 打开 RTDS 连接。
func (c *RTDSClient) Connect() error {
	if c.cfg.RTDSURL == "" {
		return errors.New("RTDS URL is required")
	}
	conn, _, err := websocket.DefaultDialer.Dial(c.cfg.RTDSURL, nil)
	if err != nil {
		return err
	}

	c.mu.Lock()
	if c.conn != nil {
		_ = c.conn.Close()
	}
	c.conn = conn
	c.mu.Unlock()

	go c.readLoop()
	go c.pingLoop()
	return nil
}

// SubscribeCryptoPrices 订阅加密货币价格更新。
func (c *RTDSClient) SubscribeCryptoPrices(source CryptoPriceSource, symbols []string, handler RTDSMessageHandler) error {
	var filters interface{}
	switch source {
	case CryptoPriceSourceBinance:
		if len(symbols) > 0 {
			filters = stringsJoin(symbols, ",")
		}
	case CryptoPriceSourceChainlink:
		if len(symbols) > 0 {
			filters = map[string]interface{}{"symbols": symbols}
		}
	default:
		return errors.New("unsupported price source")
	}

	sub := RTDSSubscription{
		Action: "subscribe",
		Subscriptions: []RTDSSubscriptionDetail{
			{
				Topic:   string(source),
				Type:    "update",
				Filters: filters,
			},
		},
	}

	c.mu.Lock()
	c.handlers[string(source)] = handler
	c.mu.Unlock()

	return c.send(sub)
}

// Unsubscribe 移除主题订阅。
func (c *RTDSClient) Unsubscribe(topic string) error {
	sub := RTDSSubscription{
		Action: "unsubscribe",
		Subscriptions: []RTDSSubscriptionDetail{
			{Topic: topic},
		},
	}

	c.mu.Lock()
	delete(c.handlers, topic)
	c.mu.Unlock()

	return c.send(sub)
}

// Close 关闭 RTDS 客户端。
func (c *RTDSClient) Close() error {
	c.cancel()
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn == nil {
		return nil
	}
	err := c.conn.Close()
	c.conn = nil
	return err
}

func (c *RTDSClient) send(msg interface{}) error {
	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()
	if conn == nil {
		return errors.New("not connected")
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return conn.WriteMessage(websocket.TextMessage, data)
}

func (c *RTDSClient) readLoop() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
		}

		c.mu.RLock()
		conn := c.conn
		c.mu.RUnlock()
		if conn == nil {
			time.Sleep(time.Second)
			continue
		}

		_, message, err := conn.ReadMessage()
		if err != nil {
			return
		}

		var msg RTDSMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		c.mu.RLock()
		handler := c.handlers[msg.Topic]
		c.mu.RUnlock()
		if handler != nil {
			_ = handler(&msg)
		}
	}
}

func (c *RTDSClient) pingLoop() {
	ticker := time.NewTicker(rtdsPingInterval)
	defer ticker.Stop()
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.mu.RLock()
			conn := c.conn
			c.mu.RUnlock()
			if conn != nil {
				_ = conn.WriteMessage(websocket.TextMessage, []byte("PING"))
			}
		}
	}
}

func stringsJoin(items []string, sep string) string {
	if len(items) == 0 {
		return ""
	}
	out := items[0]
	for i := 1; i < len(items); i++ {
		out += sep + items[i]
	}
	return out
}
