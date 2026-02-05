// clob_balance.go 模块
package polymarket

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// GetBalanceAllowance 获取余额与授权（L2 认证）。
func (c *CLOBClient) GetBalanceAllowance(ctx context.Context, params *BalanceAllowanceParams) (*BalanceAllowanceResponse, error) {
	path := EndpointGetBalanceAllowance
	headers, err := c.l2Headers(http.MethodGet, path, "")
	if err != nil {
		return nil, err
	}

	vals := url.Values{}
	vals.Set("signature_type", strconv.Itoa(c.sigType))
	if params != nil {
		if params.AssetType != "" {
			vals.Set("asset_type", string(params.AssetType))
		}
		if params.TokenID != "" {
			vals.Set("token_id", params.TokenID)
		}
	}

	var resp BalanceAllowanceResponse
	if err := c.http.Do(ctx, http.MethodGet, path, vals, nil, headers, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateBalanceAllowance 触发余额与授权刷新（L2 认证）。
func (c *CLOBClient) UpdateBalanceAllowance(ctx context.Context, params *BalanceAllowanceParams) error {
	path := EndpointUpdateBalanceAllowance
	headers, err := c.l2Headers(http.MethodGet, path, "")
	if err != nil {
		return err
	}

	vals := url.Values{}
	vals.Set("signature_type", strconv.Itoa(c.sigType))
	if params != nil {
		if params.AssetType != "" {
			vals.Set("asset_type", string(params.AssetType))
		}
		if params.TokenID != "" {
			vals.Set("token_id", params.TokenID)
		}
	}

	return c.http.Do(ctx, http.MethodGet, path, vals, nil, headers, nil)
}
