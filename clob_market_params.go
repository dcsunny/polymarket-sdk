// clob_market_params.go 模块
package polymarket

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// TickSizeResponse represents tick size response.
type TickSizeResponse struct {
	MinimumTickSize float64 `json:"minimum_tick_size"`
}

// NegRiskResponse represents neg risk response.
type NegRiskResponse struct {
	NegRisk bool `json:"neg_risk"`
}

// FeeRateResponse represents fee rate response.
type FeeRateResponse struct {
	BaseFee int `json:"base_fee"`
}

// GetTickSize returns tick size for a token.
func (c *CLOBClient) GetTickSize(tokenID string) (string, error) {
	if tokenID == "" {
		return "", ErrInvalidArgument("tokenID is required")
	}
	if tick, ok := c.tickSizeCache[tokenID]; ok {
		return tick, nil
	}

	vals := url.Values{}
	vals.Set("token_id", tokenID)
	var resp TickSizeResponse
	if err := c.http.Do(context.Background(), http.MethodGet, EndpointGetTickSize, vals, nil, nil, &resp); err != nil {
		return "", err
	}
	tick := fmt.Sprintf("%g", resp.MinimumTickSize)
	c.tickSizeCache[tokenID] = tick
	return tick, nil
}

// GetNegRisk returns neg risk flag for a token.
func (c *CLOBClient) GetNegRisk(tokenID string) (bool, error) {
	if tokenID == "" {
		return false, ErrInvalidArgument("tokenID is required")
	}
	if v, ok := c.negRiskCache[tokenID]; ok {
		return v, nil
	}

	vals := url.Values{}
	vals.Set("token_id", tokenID)
	var resp NegRiskResponse
	if err := c.http.Do(context.Background(), http.MethodGet, EndpointGetNegRisk, vals, nil, nil, &resp); err != nil {
		return false, err
	}
	c.negRiskCache[tokenID] = resp.NegRisk
	return resp.NegRisk, nil
}

// GetFeeRateBps returns fee rate for a token.
func (c *CLOBClient) GetFeeRateBps(tokenID string) (int, error) {
	if tokenID == "" {
		return 0, ErrInvalidArgument("tokenID is required")
	}
	if v, ok := c.feeRateCache[tokenID]; ok {
		return v, nil
	}

	vals := url.Values{}
	vals.Set("token_id", tokenID)
	var resp FeeRateResponse
	if err := c.http.Do(context.Background(), http.MethodGet, EndpointGetFeeRate, vals, nil, nil, &resp); err != nil {
		return 0, err
	}
	c.feeRateCache[tokenID] = resp.BaseFee
	return resp.BaseFee, nil
}

func priceValid(price float64, tickSize string) bool {
	tick, err := strconv.ParseFloat(tickSize, 64)
	if err != nil {
		return false
	}
	min := tick
	max := 1.0 - tick
	return price >= min && price <= max
}

func (c *CLOBClient) resolveTickSize(tokenID string, userTickSize string) (string, error) {
	minTick, err := c.GetTickSize(tokenID)
	if err != nil {
		return "", err
	}
	if userTickSize == "" {
		return minTick, nil
	}

	user, err1 := strconv.ParseFloat(userTickSize, 64)
	min, err2 := strconv.ParseFloat(minTick, 64)
	if err1 != nil || err2 != nil {
		return "", fmt.Errorf("invalid tick size format")
	}
	if user < min {
		return "", fmt.Errorf("invalid tick size (%s), minimum for market is %s", userTickSize, minTick)
	}
	return userTickSize, nil
}

func (c *CLOBClient) resolveFeeRate(tokenID string, userFeeRate int) (int, error) {
	marketFee, err := c.GetFeeRateBps(tokenID)
	if err != nil {
		return 0, err
	}
	if marketFee > 0 && userFeeRate > 0 && userFeeRate != marketFee {
		return 0, fmt.Errorf("invalid user fee rate: (%d), market fee rate must be %d", userFeeRate, marketFee)
	}
	return marketFee, nil
}
