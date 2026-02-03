// clob_trades.go 模块
package polymarket

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// Trade represents a trade record.
type Trade struct {
	ID              string       `json:"id"`
	TakerOrderID    string       `json:"taker_order_id"`
	Market          string       `json:"market"`
	AssetID         string       `json:"asset_id"`
	Side            string       `json:"side"`
	Size            string       `json:"size"`
	FeeRateBps      string       `json:"fee_rate_bps"`
	Price           string       `json:"price"`
	Status          string       `json:"status"`
	MatchTime       string       `json:"match_time"`
	LastUpdate      string       `json:"last_update"`
	Outcome         string       `json:"outcome"`
	BucketIndex     int          `json:"bucket_index"`
	Owner           string       `json:"owner"`
	MakerAddress    string       `json:"maker_address"`
	TransactionHash string       `json:"transaction_hash"`
	Type            string       `json:"type"`
	MakerOrders     []MakerOrder `json:"maker_orders"`
}

// GetTradesRequest filters trades.
type GetTradesRequest struct {
	ID     string
	Taker  string
	Maker  string
	Market string
	Before int64
	After  int64
}

// GetTradesResponse represents paginated trades.
type GetTradesResponse struct {
	Count      int      `json:"count"`
	Data       []*Trade `json:"data"`
	Limit      int      `json:"limit"`
	NextCursor string   `json:"next_cursor"`
}

// GetTrades returns trades and auto-pages until end.
func (c *CLOBClient) GetTrades(ctx context.Context, req *GetTradesRequest) ([]*Trade, error) {
	nextCursor := InitialCursor
	var all []*Trade

	for nextCursor != EndCursor {
		resp, err := c.GetTradesPage(ctx, req, nextCursor)
		if err != nil {
			return nil, err
		}
		all = append(all, resp.Data...)
		nextCursor = resp.NextCursor
	}

	return all, nil
}

// GetTradesPage 获取交易单页（GET /data/trades）。
// nextCursor 为空时默认使用 InitialCursor。
func (c *CLOBClient) GetTradesPage(ctx context.Context, req *GetTradesRequest, nextCursor string) (*GetTradesResponse, error) {
	path := "/data/trades"
	vals := url.Values{}
	if req != nil {
		if req.ID != "" {
			vals.Set("id", req.ID)
		}
		if req.Taker != "" {
			vals.Set("taker", req.Taker)
		}
		if req.Maker != "" {
			vals.Set("maker", req.Maker)
		}
		if req.Market != "" {
			vals.Set("market", req.Market)
		}
		if req.Before > 0 {
			vals.Set("before", strconv.FormatInt(req.Before, 10))
		}
		if req.After > 0 {
			vals.Set("after", strconv.FormatInt(req.After, 10))
		}
	}
	if nextCursor == "" {
		nextCursor = InitialCursor
	}
	vals.Set("next_cursor", nextCursor)

	headers, err := c.l2Headers(http.MethodGet, path, "")
	if err != nil {
		return nil, err
	}

	var resp GetTradesResponse
	if err := c.http.Do(ctx, http.MethodGet, path, vals, nil, headers, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetTradeByID returns a trade by id.
func (c *CLOBClient) GetTradeByID(ctx context.Context, tradeID string) (*Trade, error) {
	trades, err := c.GetTrades(ctx, &GetTradesRequest{ID: tradeID})
	if err != nil {
		return nil, err
	}
	if len(trades) == 0 {
		return nil, ErrInvalidArgument("trade not found")
	}
	return trades[0], nil
}

// GetTradesByMarket returns trades for market.
func (c *CLOBClient) GetTradesByMarket(ctx context.Context, marketID string) ([]*Trade, error) {
	return c.GetTrades(ctx, &GetTradesRequest{Market: marketID})
}

// GetTradesByAddress returns trades for address (as taker or maker).
func (c *CLOBClient) GetTradesByAddress(ctx context.Context, address string, asTaker bool) ([]*Trade, error) {
	req := &GetTradesRequest{}
	if asTaker {
		req.Taker = address
	} else {
		req.Maker = address
	}
	return c.GetTrades(ctx, req)
}
