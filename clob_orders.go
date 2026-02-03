// clob_orders.go 模块
package polymarket

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"

	order_utils_model "github.com/polymarket/go-order-utils/pkg/model"
)

// APIOrder 是 CLOB API 订单负载。
type APIOrder struct {
	Salt          int64  `json:"salt"`
	Maker         string `json:"maker"`
	Signer        string `json:"signer"`
	Taker         string `json:"taker"`
	TokenID       string `json:"tokenId"`
	MakerAmount   string `json:"makerAmount"`
	TakerAmount   string `json:"takerAmount"`
	Expiration    string `json:"expiration"`
	Nonce         string `json:"nonce"`
	FeeRateBps    string `json:"feeRateBps"`
	Side          string `json:"side"`
	SignatureType int    `json:"signatureType"`
	Signature     string `json:"signature"`
}

// OrderType 表示订单类型。
type OrderType string

const (
	OrderTypeGTC OrderType = "GTC"
	OrderTypeGTD OrderType = "GTD"
	OrderTypeFOK OrderType = "FOK"
	OrderTypeFAK OrderType = "FAK"
)

// PostOrder 封装订单提交。
type PostOrder struct {
	Order     APIOrder  `json:"order"`
	Owner     string    `json:"owner"`
	OrderType OrderType `json:"orderType"`
	// DeferExec 是否延迟执行（与 Node SDK / 官方接口字段对齐）。
	DeferExec bool `json:"deferExec"`
	// PostOnly 是否仅挂单（仅 GTC/GTD 支持；与 Node SDK 字段对齐）。
	PostOnly bool `json:"postOnly"`
}

// OrderResponse 表示订单响应。
type OrderResponse struct {
	Success     bool     `json:"success"`
	ErrorMsg    string   `json:"errorMsg,omitempty"`
	OrderID     string   `json:"orderId,omitempty"`
	OrderHashes []string `json:"orderHashes,omitempty"`
	Status      string   `json:"status,omitempty"`
}

// OrderArgs 包含限价订单参数。
type OrderArgs struct {
	TokenID     string `json:"token_id"`
	MakerAmount string `json:"maker_amount"`
	TakerAmount string `json:"taker_amount"`
	Side        string `json:"side"`
	FeeRateBps  string `json:"fee_rate_bps"`
	Nonce       string `json:"nonce"`
	Expiration  string `json:"expiration"`
	Taker       string `json:"taker"`
}

func signedOrderToAPIOrder(signedOrder *order_utils_model.SignedOrder) APIOrder {
	sideStr := SideSell
	if signedOrder.Order.Side.Int64() == 0 {
		sideStr = SideBuy
	}

	return APIOrder{
		Salt:          signedOrder.Order.Salt.Int64(),
		Maker:         signedOrder.Order.Maker.Hex(),
		Signer:        signedOrder.Order.Signer.Hex(),
		Taker:         signedOrder.Order.Taker.Hex(),
		TokenID:       signedOrder.Order.TokenId.String(),
		MakerAmount:   signedOrder.Order.MakerAmount.String(),
		TakerAmount:   signedOrder.Order.TakerAmount.String(),
		Expiration:    signedOrder.Order.Expiration.String(),
		Nonce:         signedOrder.Order.Nonce.String(),
		FeeRateBps:    signedOrder.Order.FeeRateBps.String(),
		Side:          sideStr,
		SignatureType: int(signedOrder.Order.SignatureType.Int64()),
		Signature:     "0x" + hex.EncodeToString(signedOrder.Signature),
	}
}

// CreateOrder 构建并签名限价订单。
func (c *CLOBClient) CreateOrder(args *OrderArgs) (*order_utils_model.SignedOrder, error) {
	if c.orderBuilder == nil {
		return nil, errors.New("missing private key or address for order builder")
	}
	negRisk, err := c.GetNegRisk(args.TokenID)
	if err != nil {
		return nil, err
	}
	return c.orderBuilder.BuildAndSignOrder(args, negRisk)
}

// PostOrder submits a signed order（提交单个订单，POST /order）。
func (c *CLOBClient) PostOrder(ctx context.Context, signedOrder *order_utils_model.SignedOrder, orderType OrderType) (*OrderResponse, error) {
	return c.PostOrderWithOptions(ctx, signedOrder, orderType, PostOrderOptions{})
}

// PostOrderOptions 下单可选参数（与 Node SDK 对齐）。
type PostOrderOptions struct {
	// DeferExec 是否延迟执行
	DeferExec bool
	// PostOnly 是否仅挂单（仅 GTC/GTD 支持）
	PostOnly bool
}

// PostOrderWithOptions 提交带有额外选项的已签名订单。
func (c *CLOBClient) PostOrderWithOptions(ctx context.Context, signedOrder *order_utils_model.SignedOrder, orderType OrderType, opts PostOrderOptions) (*OrderResponse, error) {
	if signedOrder == nil {
		return nil, ErrInvalidArgument("signedOrder is required")
	}
	path := "/order"

	if opts.PostOnly && orderType != OrderTypeGTC && orderType != OrderTypeGTD {
		return nil, ErrInvalidArgument("postOnly is only supported for GTC and GTD orders")
	}

	apiOrder := signedOrderToAPIOrder(signedOrder)

	postOrder := PostOrder{
		Order:     apiOrder,
		Owner:     c.APIKey(),
		OrderType: orderType,
		DeferExec: opts.DeferExec,
		PostOnly:  opts.PostOnly,
	}

	body, err := json.Marshal(postOrder)
	if err != nil {
		return nil, err
	}
	headers, err := c.l2Headers(http.MethodPost, path, string(body))
	if err != nil {
		return nil, err
	}

	// builder flow：如果配置了 builder auth，则注入 builder headers
	if c.builderAuth != nil {
		if bh, berr := c.builderAuth.Headers(http.MethodPost, path, body); berr == nil {
			headers = mergeHeaders(headers, bh)
		}
	}

	var resp OrderResponse
	if err := c.http.DoRaw(ctx, http.MethodPost, path, nil, body, headers, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// PostOrdersArgs 批量下单参数（与 Node SDK PostOrdersArgs 对齐）。
type PostOrdersArgs struct {
	Order     *order_utils_model.SignedOrder
	OrderType OrderType
	PostOnly  *bool
}

// PostOrdersOptions 批量下单选项。
type PostOrdersOptions struct {
	// DeferExec 是否延迟执行（对所有订单生效）
	DeferExec bool
	// DefaultPostOnly 默认 postOnly（当 PostOrdersArgs.PostOnly 为 nil 时生效）
	DefaultPostOnly bool
}

// PostOrdersSigned 批量提交 SignedOrder（POST /orders）。
func (c *CLOBClient) PostOrdersSigned(ctx context.Context, args []PostOrdersArgs, opts PostOrdersOptions) ([]*OrderResponse, error) {
	if len(args) == 0 {
		return nil, ErrInvalidArgument("args is required")
	}
	if len(args) > 15 {
		return nil, ErrInvalidArgument("max 15 orders per batch")
	}

	orders := make([]*PostOrder, 0, len(args))
	for _, a := range args {
		if a.Order == nil {
			return nil, ErrInvalidArgument("order is required")
		}
		postOnly := opts.DefaultPostOnly
		if a.PostOnly != nil {
			postOnly = *a.PostOnly
		}
		if postOnly && a.OrderType != OrderTypeGTC && a.OrderType != OrderTypeGTD {
			return nil, ErrInvalidArgument("postOnly is only supported for GTC and GTD orders")
		}

		orders = append(orders, &PostOrder{
			Order:     signedOrderToAPIOrder(a.Order),
			Owner:     c.APIKey(),
			OrderType: a.OrderType,
			DeferExec: opts.DeferExec,
			PostOnly:  postOnly,
		})
	}

	return c.PostOrders(ctx, orders)
}

// PostOrders 提交多个订单。
func (c *CLOBClient) PostOrders(ctx context.Context, orders []*PostOrder) ([]*OrderResponse, error) {
	if len(orders) == 0 {
		return nil, ErrInvalidArgument("orders is required")
	}
	if len(orders) > 15 {
		return nil, ErrInvalidArgument("max 15 orders per batch")
	}
	path := "/orders"
	body, err := json.Marshal(orders)
	if err != nil {
		return nil, err
	}
	headers, err := c.l2Headers(http.MethodPost, path, string(body))
	if err != nil {
		return nil, err
	}

	// builder flow：如果配置了 builder auth，则注入 builder headers
	if c.builderAuth != nil {
		if bh, berr := c.builderAuth.Headers(http.MethodPost, path, body); berr == nil {
			headers = mergeHeaders(headers, bh)
		}
	}

	var resp []*OrderResponse
	if err := c.http.DoRaw(ctx, http.MethodPost, path, nil, body, headers, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// CancelOrders 取消多个订单。
func (c *CLOBClient) CancelOrders(ctx context.Context, orderIDs []string) (*CancelOrdersResponse, error) {
	if len(orderIDs) == 0 {
		return nil, ErrInvalidArgument("orderIDs is required")
	}
	path := "/orders"
	body, err := json.Marshal(orderIDs)
	if err != nil {
		return nil, err
	}
	headers, err := c.l2Headers(http.MethodDelete, path, string(body))
	if err != nil {
		return nil, err
	}
	var resp CancelOrdersResponse
	if err := c.http.DoRaw(ctx, http.MethodDelete, path, nil, body, headers, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CancelOrdersResponse 表示取消响应。
type CancelOrdersResponse struct {
	Canceled    []string          `json:"canceled"`
	NotCanceled map[string]string `json:"not_canceled"`
}

// CancelOrder 取消单个订单。
func (c *CLOBClient) CancelOrder(ctx context.Context, orderID string) (*CancelOrdersResponse, error) {
	if orderID == "" {
		return nil, ErrInvalidArgument("orderID is required")
	}
	path := "/order"
	payload := map[string]string{"orderID": orderID}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	headers, err := c.l2Headers(http.MethodDelete, path, string(body))
	if err != nil {
		return nil, err
	}
	var resp CancelOrdersResponse
	if err := c.http.DoRaw(ctx, http.MethodDelete, path, nil, body, headers, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CancelAllOrders 取消所有订单。
func (c *CLOBClient) CancelAllOrders(ctx context.Context) (*CancelOrdersResponse, error) {
	path := "/cancel-all"
	headers, err := c.l2Headers(http.MethodDelete, path, "")
	if err != nil {
		return nil, err
	}
	var resp CancelOrdersResponse
	if err := c.http.Do(ctx, http.MethodDelete, path, nil, nil, headers, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CancelAll 是 CancelAllOrders 的别名（与 Node SDK 命名对齐）。
func (c *CLOBClient) CancelAll(ctx context.Context) (*CancelOrdersResponse, error) {
	return c.CancelAllOrders(ctx)
}

// CancelMarketOrders 取消市场或资产的订单。
func (c *CLOBClient) CancelMarketOrders(ctx context.Context, req CancelMarketOrdersRequest) (*CancelOrdersResponse, error) {
	path := "/cancel-market-orders"
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	headers, err := c.l2Headers(http.MethodDelete, path, string(body))
	if err != nil {
		return nil, err
	}
	var resp CancelOrdersResponse
	if err := c.http.DoRaw(ctx, http.MethodDelete, path, nil, body, headers, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func mergeHeaders(base, extra map[string]string) map[string]string {
	if len(extra) == 0 {
		return base
	}
	out := make(map[string]string, len(base)+len(extra))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range extra {
		out[k] = v
	}
	return out
}
