package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	pm "github.com/dcsunny/polymarket-sdk"
	"github.com/joho/godotenv"
)

// 监听某个市场的订单簿

func main() {
	_ = godotenv.Load()

	eventSlug := os.Getenv("EVENT_SLUG")
	assetID := os.Getenv("ASSET_ID")
	if eventSlug == "" && assetID == "" {
		fmt.Println("set EVENT_SLUG or ASSET_ID")
		return
	}

	sdk, err := pm.New(pm.Config{})
	if err != nil {
		fmt.Printf("init sdk failed: %v\n", err)
		return
	}

	if assetID == "" {
		ctx := context.Background()
		event, err := sdk.REST.EventBySlug(ctx, eventSlug, pm.EventBySlugQuery{})
		if err != nil {
			fmt.Printf("load event failed: %v\n", err)
			return
		}
		assetID = extractAssetID(event)
		if assetID == "" {
			fmt.Println("cannot resolve asset id from event; set ASSET_ID manually")
			return
		}
	}

	if err := sdk.WSS.ConnectMarketChannel(); err != nil {
		fmt.Printf("connect wss failed: %v\n", err)
		return
	}

	handlers := map[string]pm.WSSMessageHandler{
		"book": func(data json.RawMessage) error {
			var msg pm.WSSBookMessage
			//fmt.Println(string(data))
			if err := json.Unmarshal(data, &msg); err != nil {
				return err
			}
			bestBid := ""
			bestBidSize := ""
			bestAsk := ""
			bestAskSize := ""
			if len(msg.Bids) > 0 {
				bidIndex := len(msg.Bids) - 1
				bestBid = msg.Bids[bidIndex].Price
				bestBidSize = msg.Bids[bidIndex].Size
			}
			if len(msg.Asks) > 0 {
				askIndex := len(msg.Asks) - 1
				bestAsk = msg.Asks[askIndex].Price
				bestAskSize = msg.Asks[askIndex].Size
			}
			fmt.Printf("book asset=%s bid=%s bidSize=%s ask=%s askSize=%s\n", msg.AssetID, bestBid, bestBidSize, bestAsk, bestAskSize)
			return nil
		},
	}

	if err := sdk.WSS.SubscribeMarketChannel([]string{assetID}, handlers); err != nil {
		fmt.Printf("subscribe failed: %v\n", err)
		return
	}

	fmt.Printf("listening orderbook for asset_id=%s\n", assetID)
	select {}
}

func extractAssetID(event *pm.Event) string {
	for _, raw := range event.Markets {
		m, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		if v, ok := m["clobTokenIds"]; ok {
			if id := pickToken(v); id != "" {
				return id
			}
		}
		if v, ok := m["clob_token_ids"]; ok {
			if id := pickToken(v); id != "" {
				return id
			}
		}
	}
	return ""
}

func pickToken(v interface{}) string {
	switch t := v.(type) {
	case []interface{}:
		for _, item := range t {
			if s, ok := item.(string); ok && s != "" {
				return s
			}
		}
	case []string:
		if len(t) > 0 {
			return t[0]
		}
	case string:
		if t == "" {
			return ""
		}
		if strings.HasPrefix(t, "[") {
			var arr []string
			if err := json.Unmarshal([]byte(t), &arr); err == nil && len(arr) > 0 {
				return arr[0]
			}
		}
		parts := strings.Split(t, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	return ""
}
