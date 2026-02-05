// clob_rfq.go 模块
package polymarket

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

type createRfqRequestBody struct {
	AssetIn   string `json:"assetIn"`
	AssetOut  string `json:"assetOut"`
	AmountIn  string `json:"amountIn"`
	AmountOut string `json:"amountOut"`
	UserType  int    `json:"userType"`
}

type createRfqQuoteBody struct {
	RequestID string `json:"requestId"`
	AssetIn   string `json:"assetIn"`
	AssetOut  string `json:"assetOut"`
	AmountIn  string `json:"amountIn"`
	AmountOut string `json:"amountOut"`
	UserType  int    `json:"userType"`
}

// CreateRfqRequest 创建 RFQ request（POST /rfq/request，L2 认证）。
func (c *CLOBClient) CreateRfqRequest(ctx context.Context, payload CreateRfqRequestPayload) (*RfqRequestResponse, error) {
	if payload.AssetIn == "" || payload.AssetOut == "" || payload.AmountIn == "" || payload.AmountOut == "" {
		return nil, ErrInvalidArgument("assetIn/assetOut/amountIn/amountOut are required")
	}
	path := EndpointCreateRfqRequest

	bodyObj := createRfqRequestBody{
		AssetIn:   payload.AssetIn,
		AssetOut:  payload.AssetOut,
		AmountIn:  payload.AmountIn,
		AmountOut: payload.AmountOut,
		UserType:  c.sigType,
	}
	body, err := json.Marshal(bodyObj)
	if err != nil {
		return nil, err
	}

	headers, err := c.l2Headers(http.MethodPost, path, string(body))
	if err != nil {
		return nil, err
	}

	var resp RfqRequestResponse
	if err := c.http.DoRaw(ctx, http.MethodPost, path, nil, body, headers, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CancelRfqRequest 取消 RFQ request（DELETE /rfq/request，L2 认证）。
func (c *CLOBClient) CancelRfqRequest(ctx context.Context, requestID string) (string, error) {
	if requestID == "" {
		return "", ErrInvalidArgument("requestID is required")
	}
	path := EndpointCancelRfqRequest

	bodyObj := CancelRfqRequestParams{RequestID: requestID}
	body, err := json.Marshal(bodyObj)
	if err != nil {
		return "", err
	}

	headers, err := c.l2Headers(http.MethodDelete, path, string(body))
	if err != nil {
		return "", err
	}

	var raw []byte
	if err := c.http.DoRaw(ctx, http.MethodDelete, path, nil, body, headers, &raw); err != nil {
		return "", err
	}
	return decodeMaybeJSONString(raw), nil
}

// GetRfqRequests 获取 RFQ requests（GET /rfq/data/requests，L2 认证）。
func (c *CLOBClient) GetRfqRequests(ctx context.Context, params *GetRfqRequestsParams) (*PaginatedResponse[RfqRequest], error) {
	path := EndpointGetRfqRequests
	vals := url.Values{}
	if params != nil {
		if params.Offset != "" {
			vals.Set("offset", params.Offset)
		}
		if params.Limit > 0 {
			vals.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.State != "" {
			vals.Set("state", string(params.State))
		}
		addRepeated(vals, "requestIds", params.RequestIDs)
		addRepeated(vals, "markets", params.Markets)

		addFloat(vals, "sizeMin", params.SizeMin)
		addFloat(vals, "sizeMax", params.SizeMax)
		addFloat(vals, "sizeUsdcMin", params.SizeUsdcMin)
		addFloat(vals, "sizeUsdcMax", params.SizeUsdcMax)
		addFloat(vals, "priceMin", params.PriceMin)
		addFloat(vals, "priceMax", params.PriceMax)

		if params.SortBy != "" {
			vals.Set("sortBy", string(params.SortBy))
		}
		if params.SortDir != "" {
			vals.Set("sortDir", string(params.SortDir))
		}
	}

	headers, err := c.l2Headers(http.MethodGet, path, "")
	if err != nil {
		return nil, err
	}

	var resp PaginatedResponse[RfqRequest]
	if err := c.http.Do(ctx, http.MethodGet, path, vals, nil, headers, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateRfqQuote 创建 RFQ quote（POST /rfq/quote，L2 认证）。
func (c *CLOBClient) CreateRfqQuote(ctx context.Context, payload CreateRfqQuotePayload) (*RfqQuoteResponse, error) {
	if payload.RequestID == "" || payload.AssetIn == "" || payload.AssetOut == "" || payload.AmountIn == "" || payload.AmountOut == "" {
		return nil, ErrInvalidArgument("requestId/assetIn/assetOut/amountIn/amountOut are required")
	}
	path := EndpointCreateRfqQuote

	bodyObj := createRfqQuoteBody{
		RequestID: payload.RequestID,
		AssetIn:   payload.AssetIn,
		AssetOut:  payload.AssetOut,
		AmountIn:  payload.AmountIn,
		AmountOut: payload.AmountOut,
		UserType:  c.sigType,
	}
	body, err := json.Marshal(bodyObj)
	if err != nil {
		return nil, err
	}

	headers, err := c.l2Headers(http.MethodPost, path, string(body))
	if err != nil {
		return nil, err
	}

	var resp RfqQuoteResponse
	if err := c.http.DoRaw(ctx, http.MethodPost, path, nil, body, headers, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CancelRfqQuote 取消 RFQ quote（DELETE /rfq/quote，L2 认证）。
func (c *CLOBClient) CancelRfqQuote(ctx context.Context, quoteID string) (string, error) {
	if quoteID == "" {
		return "", ErrInvalidArgument("quoteID is required")
	}
	path := EndpointCancelRfqQuote

	bodyObj := CancelRfqQuoteParams{QuoteID: quoteID}
	body, err := json.Marshal(bodyObj)
	if err != nil {
		return "", err
	}

	headers, err := c.l2Headers(http.MethodDelete, path, string(body))
	if err != nil {
		return "", err
	}

	var raw []byte
	if err := c.http.DoRaw(ctx, http.MethodDelete, path, nil, body, headers, &raw); err != nil {
		return "", err
	}
	return decodeMaybeJSONString(raw), nil
}

// GetRfqRequesterQuotes 获取 requester 视角的 quotes（GET /rfq/data/requester/quotes，L2 认证）。
func (c *CLOBClient) GetRfqRequesterQuotes(ctx context.Context, params *GetRfqQuotesParams) (*PaginatedResponse[RfqQuote], error) {
	return c.getRfqQuotes(ctx, EndpointGetRfqRequesterQuotes, params)
}

// GetRfqQuoterQuotes 获取 quoter 视角的 quotes（GET /rfq/data/quoter/quotes，L2 认证）。
func (c *CLOBClient) GetRfqQuoterQuotes(ctx context.Context, params *GetRfqQuotesParams) (*PaginatedResponse[RfqQuote], error) {
	return c.getRfqQuotes(ctx, EndpointGetRfqQuoterQuotes, params)
}

func (c *CLOBClient) getRfqQuotes(ctx context.Context, path string, params *GetRfqQuotesParams) (*PaginatedResponse[RfqQuote], error) {
	vals := url.Values{}
	if params != nil {
		if params.Offset != "" {
			vals.Set("offset", params.Offset)
		}
		if params.Limit > 0 {
			vals.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.State != "" {
			vals.Set("state", string(params.State))
		}

		addRepeated(vals, "quoteIds", params.QuoteIDs)
		addRepeated(vals, "requestIds", params.RequestIDs)
		addRepeated(vals, "markets", params.Markets)

		addFloat(vals, "sizeMin", params.SizeMin)
		addFloat(vals, "sizeMax", params.SizeMax)
		addFloat(vals, "sizeUsdcMin", params.SizeUsdcMin)
		addFloat(vals, "sizeUsdcMax", params.SizeUsdcMax)
		addFloat(vals, "priceMin", params.PriceMin)
		addFloat(vals, "priceMax", params.PriceMax)

		if params.SortBy != "" {
			vals.Set("sortBy", string(params.SortBy))
		}
		if params.SortDir != "" {
			vals.Set("sortDir", string(params.SortDir))
		}
	}

	headers, err := c.l2Headers(http.MethodGet, path, "")
	if err != nil {
		return nil, err
	}

	var resp PaginatedResponse[RfqQuote]
	if err := c.http.Do(ctx, http.MethodGet, path, vals, nil, headers, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetRfqBestQuote 获取某个 request 的最佳 quote（GET /rfq/data/best-quote，L2 认证）。
func (c *CLOBClient) GetRfqBestQuote(ctx context.Context, params *GetRfqBestQuoteParams) (*RfqQuote, error) {
	path := EndpointGetRfqBestQuote
	vals := url.Values{}
	if params != nil && params.RequestID != "" {
		vals.Set("requestId", params.RequestID)
	}

	headers, err := c.l2Headers(http.MethodGet, path, "")
	if err != nil {
		return nil, err
	}

	var resp RfqQuote
	if err := c.http.Do(ctx, http.MethodGet, path, vals, nil, headers, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetRfqConfig 获取 RFQ 配置（GET /rfq/config，L2 认证）。
func (c *CLOBClient) GetRfqConfig(ctx context.Context) (json.RawMessage, error) {
	path := EndpointRfqConfig
	headers, err := c.l2Headers(http.MethodGet, path, "")
	if err != nil {
		return nil, err
	}

	var resp json.RawMessage
	if err := c.http.Do(ctx, http.MethodGet, path, nil, nil, headers, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// AcceptRfqQuote 接受 RFQ quote（POST /rfq/request/accept，L2 认证）。
// payload 的具体字段随官方接口变动较频繁，这里直接透传 JSON（对齐 Node SDK 的行为）。
func (c *CLOBClient) AcceptRfqQuote(ctx context.Context, payload any) (string, error) {
	path := EndpointRfqRequestsAccept
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	headers, err := c.l2Headers(http.MethodPost, path, string(body))
	if err != nil {
		return "", err
	}
	var raw []byte
	if err := c.http.DoRaw(ctx, http.MethodPost, path, nil, body, headers, &raw); err != nil {
		return "", err
	}
	return decodeMaybeJSONString(raw), nil
}

// ApproveRfqOrder 审批 RFQ quote 并创建订单（POST /rfq/quote/approve，L2 认证）。
// payload 直接透传 JSON。
func (c *CLOBClient) ApproveRfqOrder(ctx context.Context, payload any) (string, error) {
	path := EndpointRfqQuoteApprove
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	headers, err := c.l2Headers(http.MethodPost, path, string(body))
	if err != nil {
		return "", err
	}
	var raw []byte
	if err := c.http.DoRaw(ctx, http.MethodPost, path, nil, body, headers, &raw); err != nil {
		return "", err
	}
	return decodeMaybeJSONString(raw), nil
}

func addRepeated(vals url.Values, key string, items []string) {
	for _, it := range items {
		if it != "" {
			vals.Add(key, it)
		}
	}
}

func addFloat(vals url.Values, key string, v *float64) {
	if v == nil {
		return
	}
	vals.Set(key, strconv.FormatFloat(*v, 'f', -1, 64))
}

func decodeMaybeJSONString(raw []byte) string {
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}
	return string(raw)
}
