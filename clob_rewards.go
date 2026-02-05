// clob_rewards.go 模块
package polymarket

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
)

// GetEarningsForUserForDay 获取用户某一天的收益明细（GET /rewards/user，L2 认证，自动分页）。
func (c *CLOBClient) GetEarningsForUserForDay(ctx context.Context, date string) ([]UserEarning, error) {
	if date == "" {
		return nil, ErrInvalidArgument("date is required")
	}
	path := EndpointGetEarningsForUserForDay

	nextCursor := InitialCursor
	var all []UserEarning
	for nextCursor != EndCursor {
		vals := url.Values{}
		vals.Set("date", date)
		vals.Set("signature_type", strconv.Itoa(c.sigType))
		vals.Set("next_cursor", nextCursor)

		headers, err := c.l2Headers(http.MethodGet, path, "")
		if err != nil {
			return nil, err
		}

		var page PaginatedResponse[UserEarning]
		if err := c.http.Do(ctx, http.MethodGet, path, vals, nil, headers, &page); err != nil {
			return nil, err
		}
		all = append(all, page.Data...)
		nextCursor = page.NextCursor
	}
	return all, nil
}

// GetTotalEarningsForUserForDay 获取用户某一天的总收益（GET /rewards/user/total，L2 认证）。
func (c *CLOBClient) GetTotalEarningsForUserForDay(ctx context.Context, date string) ([]TotalUserEarning, error) {
	if date == "" {
		return nil, ErrInvalidArgument("date is required")
	}
	path := EndpointGetTotalEarningsForUserForDay
	vals := url.Values{}
	vals.Set("date", date)
	vals.Set("signature_type", strconv.Itoa(c.sigType))

	headers, err := c.l2Headers(http.MethodGet, path, "")
	if err != nil {
		return nil, err
	}

	var resp []TotalUserEarning
	if err := c.http.Do(ctx, http.MethodGet, path, vals, nil, headers, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetUserEarningsAndMarketsConfig 获取用户某天在各市场的收益与奖励配置（GET /rewards/user/markets，L2 认证，自动分页）。
func (c *CLOBClient) GetUserEarningsAndMarketsConfig(ctx context.Context, date string, orderBy string, position string, noCompetition bool) ([]UserRewardsEarning, error) {
	if date == "" {
		return nil, ErrInvalidArgument("date is required")
	}
	path := EndpointGetRewardsEarningsPercentages

	nextCursor := InitialCursor
	var all []UserRewardsEarning
	for nextCursor != EndCursor {
		vals := url.Values{}
		vals.Set("date", date)
		vals.Set("signature_type", strconv.Itoa(c.sigType))
		vals.Set("next_cursor", nextCursor)
		if orderBy != "" {
			vals.Set("order_by", orderBy)
		}
		if position != "" {
			vals.Set("position", position)
		}
		if noCompetition {
			vals.Set("no_competition", "true")
		}

		headers, err := c.l2Headers(http.MethodGet, path, "")
		if err != nil {
			return nil, err
		}

		var page PaginatedResponse[UserRewardsEarning]
		if err := c.http.Do(ctx, http.MethodGet, path, vals, nil, headers, &page); err != nil {
			return nil, err
		}
		all = append(all, page.Data...)
		nextCursor = page.NextCursor
	}
	return all, nil
}

// GetRewardPercentages 获取用户各市场奖励占比（GET /rewards/user/percentages，L2 认证）。
func (c *CLOBClient) GetRewardPercentages(ctx context.Context) (RewardsPercentages, error) {
	path := EndpointGetLiquidityRewardPercentages
	vals := url.Values{}
	vals.Set("signature_type", strconv.Itoa(c.sigType))

	headers, err := c.l2Headers(http.MethodGet, path, "")
	if err != nil {
		return nil, err
	}

	var resp RewardsPercentages
	if err := c.http.Do(ctx, http.MethodGet, path, vals, nil, headers, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetCurrentRewards 获取当前所有奖励市场（GET /rewards/markets/current，公开接口，自动分页）。
func (c *CLOBClient) GetCurrentRewards(ctx context.Context) ([]MarketReward, error) {
	path := EndpointGetRewardsMarketsCurrent
	nextCursor := InitialCursor
	var all []MarketReward

	for nextCursor != EndCursor {
		vals := url.Values{}
		vals.Set("next_cursor", nextCursor)

		var page PaginatedResponse[MarketReward]
		if err := c.http.Do(ctx, http.MethodGet, path, vals, nil, nil, &page); err != nil {
			return nil, err
		}
		all = append(all, page.Data...)
		nextCursor = page.NextCursor
	}
	return all, nil
}

// GetRawRewardsForMarket 获取某个市场（condition_id）的奖励配置（GET /rewards/markets/{condition_id}，公开接口，自动分页）。
func (c *CLOBClient) GetRawRewardsForMarket(ctx context.Context, conditionID string) ([]MarketReward, error) {
	if conditionID == "" {
		return nil, ErrInvalidArgument("conditionID is required")
	}
	path := EndpointGetRewardsMarketsPrefix + url.PathEscape(conditionID)
	nextCursor := InitialCursor
	var all []MarketReward

	for nextCursor != EndCursor {
		vals := url.Values{}
		vals.Set("next_cursor", nextCursor)

		var page PaginatedResponse[MarketReward]
		if err := c.http.Do(ctx, http.MethodGet, path, vals, nil, nil, &page); err != nil {
			return nil, err
		}
		all = append(all, page.Data...)
		nextCursor = page.NextCursor
	}
	return all, nil
}
