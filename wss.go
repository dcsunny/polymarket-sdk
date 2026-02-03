// wss.go 模块
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
	wssPingInterval = 10 * time.Second
)

// WSSClient 处理 Polymarket WebSocket 订阅。
type WSSClient struct {
	cfg Config

	mu       sync.RWMutex
	conn     *websocket.Conn
	handlers map[string]WSSMessageHandler

	ctx    context.Context
	cancel context.CancelFunc
}

func NewWSSClient(cfg Config) *WSSClient {
	ctx, cancel := context.WithCancel(context.Background())
	return &WSSClient{
		cfg:      cfg,
		handlers: make(map[string]WSSMessageHandler),
		ctx:      ctx,
		cancel:   cancel,
	}
}

// ConnectUserChannel 连接到用户频道。
func (c *WSSClient) ConnectUserChannel() error {
	return c.connect(c.cfg.WSSUserURL)
}

// ConnectMarketChannel 连接到市场频道。
func (c *WSSClient) ConnectMarketChannel() error {
	return c.connect(c.cfg.WSSMarketURL)
}

func (c *WSSClient) connect(endpoint string) error {
	if endpoint == "" {
		return errors.New("endpoint is required")
	}
	conn, _, err := websocket.DefaultDialer.Dial(endpoint, nil)
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

// SubscribeUserChannel 订阅用户事件。
func (c *WSSClient) SubscribeUserChannel(markets []string, handlers map[string]WSSMessageHandler) error {
	if c.cfg.APIKey == "" || c.cfg.APISecret == "" || c.cfg.Passphrase == "" {
		return errors.New("missing API credentials for user channel")
	}
	sub := WSSSubscription{
		Auth: &WSSAuth{
			APIKey:     c.cfg.APIKey,
			Secret:     c.cfg.APISecret,
			Passphrase: c.cfg.Passphrase,
		},
		Type:    WSSChannelTypeUser,
		Markets: markets,
	}
	c.registerHandlers(handlers)
	return c.send(sub)
}

// SubscribeMarketChannel 订阅市场事件。
func (c *WSSClient) SubscribeMarketChannel(assetIDs []string, handlers map[string]WSSMessageHandler) error {
	sub := WSSSubscription{
		Type:      WSSChannelTypeMarket,
		AssetsIDs: assetIDs,
	}
	c.registerHandlers(handlers)
	return c.send(sub)
}

// RegisterHandler 为事件类型注册处理器。
func (c *WSSClient) RegisterHandler(eventType string, handler WSSMessageHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers[eventType] = handler
}

// Close 关闭连接。
func (c *WSSClient) Close() error {
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

func (c *WSSClient) registerHandlers(handlers map[string]WSSMessageHandler) {
	if handlers == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, h := range handlers {
		c.handlers[k] = h
	}
}

func (c *WSSClient) send(msg interface{}) error {
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

func (c *WSSClient) readLoop() {
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

		if len(message) == 0 {
			continue
		}
		if string(message) == "PING" || string(message) == "PONG" {
			continue
		}

		if message[0] == '[' {
			var items []json.RawMessage
			if err := json.Unmarshal(message, &items); err != nil {
				continue
			}
			for _, item := range items {
				c.handleMessage(item)
			}
			continue
		}

		c.handleMessage(json.RawMessage(message))
	}
}

func (c *WSSClient) handleMessage(msg json.RawMessage) {
	var base struct {
		EventType string `json:"event_type"`
	}
	if err := json.Unmarshal(msg, &base); err != nil {
		return
	}

	c.mu.RLock()
	handler := c.handlers[base.EventType]
	c.mu.RUnlock()

	if handler != nil {
		_ = handler(msg)
	}
}

func (c *WSSClient) pingLoop() {
	ticker := time.NewTicker(wssPingInterval)
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
