// rest.go 模块
package polymarket

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dcsunny/polymarket-sdk/internal/httpx"
)

// RESTClient 处理 Polymarket REST API。
type RESTClient struct {
	http *httpx.Client
}

func NewRESTClient(http *httpx.Client) *RESTClient {
	return &RESTClient{http: http}
}

// EventsQuery 过滤事件列表。
type EventsQuery struct {
	Limit  int
	Offset int

	Order     string
	Ascending bool

	IDs   []int
	Slugs []string

	TagID         *int
	ExcludeTagIDs []int
	RelatedTags   *bool

	Featured *bool
	Closed   *bool
	CYOM     *bool

	IncludeChat     *bool
	IncludeTemplate *bool

	Recurrence string

	StartDateMin *time.Time
	StartDateMax *time.Time
	EndDateMin   *time.Time
	EndDateMax   *time.Time
}

// Events 返回事件列表。
func (c *RESTClient) Events(ctx context.Context, q EventsQuery) ([]*Event, error) {
	if q.Order == "" {
		q.Order = "id"
	}
	if q.Limit <= 0 {
		q.Limit = 100
	}

	vals := url.Values{}
	vals.Set("limit", strconv.Itoa(q.Limit))
	vals.Set("offset", strconv.Itoa(q.Offset))
	vals.Set("order", q.Order)
	vals.Set("ascending", strconv.FormatBool(q.Ascending))

	if len(q.IDs) > 0 {
		ids := make([]string, len(q.IDs))
		for i, id := range q.IDs {
			ids[i] = strconv.Itoa(id)
		}
		vals.Set("id", strings.Join(ids, ","))
	}
	if len(q.Slugs) > 0 {
		vals.Set("slug", strings.Join(q.Slugs, ","))
	}
	if q.TagID != nil {
		vals.Set("tag_id", strconv.Itoa(*q.TagID))
	}
	if q.RelatedTags != nil {
		vals.Set("related_tags", strconv.FormatBool(*q.RelatedTags))
	}
	if q.Featured != nil {
		vals.Set("featured", strconv.FormatBool(*q.Featured))
	}
	if q.Closed != nil {
		vals.Set("closed", strconv.FormatBool(*q.Closed))
	}
	if q.CYOM != nil {
		vals.Set("cyom", strconv.FormatBool(*q.CYOM))
	}
	if q.IncludeChat != nil {
		vals.Set("include_chat", strconv.FormatBool(*q.IncludeChat))
	}
	if q.IncludeTemplate != nil {
		vals.Set("include_template", strconv.FormatBool(*q.IncludeTemplate))
	}
	if q.Recurrence != "" {
		vals.Set("recurrence", q.Recurrence)
	}
	if q.StartDateMin != nil {
		vals.Set("start_date_min", q.StartDateMin.Format(time.RFC3339))
	}
	if q.StartDateMax != nil {
		vals.Set("start_date_max", q.StartDateMax.Format(time.RFC3339))
	}
	if q.EndDateMin != nil {
		vals.Set("end_date_min", q.EndDateMin.Format(time.RFC3339))
	}
	if q.EndDateMax != nil {
		vals.Set("end_date_max", q.EndDateMax.Format(time.RFC3339))
	}

	for _, id := range q.ExcludeTagIDs {
		vals.Add("exclude_tag_id", strconv.Itoa(id))
	}

	var events []*Event
	if err := c.http.Do(ctx, http.MethodGet, "/events", vals, nil, nil, &events); err != nil {
		return nil, err
	}
	return events, nil
}

// EventBySlugQuery 控制 slug 获取选项。
type EventBySlugQuery struct {
	IncludeChat     *bool
	IncludeTemplate *bool
}

// EventBySlug 根据 slug 获取单个事件。
func (c *RESTClient) EventBySlug(ctx context.Context, slug string, q EventBySlugQuery) (*Event, error) {
	if slug == "" {
		return nil, ErrInvalidArgument("slug is required")
	}

	vals := url.Values{}
	if q.IncludeChat != nil {
		vals.Set("include_chat", strconv.FormatBool(*q.IncludeChat))
	}
	if q.IncludeTemplate != nil {
		vals.Set("include_template", strconv.FormatBool(*q.IncludeTemplate))
	}

	var event Event
	path := "/events/slug/" + url.PathEscape(slug)
	if err := c.http.Do(ctx, http.MethodGet, path, vals, nil, nil, &event); err != nil {
		return nil, err
	}
	return &event, nil
}

// MarketsQuery 过滤市场列表。
type MarketsQuery struct {
	Limit  int
	Offset int

	Order     string
	Ascending bool

	IDs                []string
	Slug               string
	ClobTokenIDs       []string
	ConditionIDs       []string
	MarketMakerAddress []string

	LiquidityNumMin *float64
	LiquidityNumMax *float64
	VolumeNumMin    *float64
	VolumeNumMax    *float64

	StartDateMin *time.Time
	StartDateMax *time.Time
	EndDateMin   *time.Time
	EndDateMax   *time.Time

	TagID          *int
	RelatedTags    *bool
	IncludeTag     *bool
	Closed         *bool
	QuestionIDs    []string
	RewardsMinSize *float64

	GameID              string
	SportsMarketTypes   []string
	UMAResolutionStatus string
}

// Markets 返回市场列表。
func (c *RESTClient) Markets(ctx context.Context, q MarketsQuery) ([]*Market, error) {
	if q.Limit <= 0 {
		q.Limit = 100
	}

	vals := url.Values{}
	vals.Set("limit", strconv.Itoa(q.Limit))
	vals.Set("offset", strconv.Itoa(q.Offset))
	if q.Order != "" {
		vals.Set("order", q.Order)
	}
	vals.Set("ascending", strconv.FormatBool(q.Ascending))

	if len(q.IDs) > 0 {
		vals.Set("id", strings.Join(q.IDs, ","))
	}
	if q.Slug != "" {
		vals.Set("slug", q.Slug)
	}
	if len(q.ClobTokenIDs) > 0 {
		vals.Set("clob_token_ids", strings.Join(q.ClobTokenIDs, ","))
	}
	if len(q.ConditionIDs) > 0 {
		vals.Set("condition_ids", strings.Join(q.ConditionIDs, ","))
	}
	if len(q.MarketMakerAddress) > 0 {
		vals.Set("market_maker_address", strings.Join(q.MarketMakerAddress, ","))
	}
	if q.LiquidityNumMin != nil {
		vals.Set("liquidity_num_min", strconv.FormatFloat(*q.LiquidityNumMin, 'f', -1, 64))
	}
	if q.LiquidityNumMax != nil {
		vals.Set("liquidity_num_max", strconv.FormatFloat(*q.LiquidityNumMax, 'f', -1, 64))
	}
	if q.VolumeNumMin != nil {
		vals.Set("volume_num_min", strconv.FormatFloat(*q.VolumeNumMin, 'f', -1, 64))
	}
	if q.VolumeNumMax != nil {
		vals.Set("volume_num_max", strconv.FormatFloat(*q.VolumeNumMax, 'f', -1, 64))
	}
	if q.StartDateMin != nil {
		vals.Set("start_date_min", q.StartDateMin.Format(time.RFC3339))
	}
	if q.StartDateMax != nil {
		vals.Set("start_date_max", q.StartDateMax.Format(time.RFC3339))
	}
	if q.EndDateMin != nil {
		vals.Set("end_date_min", q.EndDateMin.Format(time.RFC3339))
	}
	if q.EndDateMax != nil {
		vals.Set("end_date_max", q.EndDateMax.Format(time.RFC3339))
	}
	if q.TagID != nil {
		vals.Set("tag_id", strconv.Itoa(*q.TagID))
	}
	if q.RelatedTags != nil {
		vals.Set("related_tags", strconv.FormatBool(*q.RelatedTags))
	}
	if q.IncludeTag != nil {
		vals.Set("include_tag", strconv.FormatBool(*q.IncludeTag))
	}
	if q.Closed != nil {
		vals.Set("closed", strconv.FormatBool(*q.Closed))
	}
	if len(q.QuestionIDs) > 0 {
		vals.Set("question_ids", strings.Join(q.QuestionIDs, ","))
	}
	if q.RewardsMinSize != nil {
		vals.Set("rewards_min_size", strconv.FormatFloat(*q.RewardsMinSize, 'f', -1, 64))
	}
	if q.GameID != "" {
		vals.Set("game_id", q.GameID)
	}
	if len(q.SportsMarketTypes) > 0 {
		vals.Set("sports_market_types", strings.Join(q.SportsMarketTypes, ","))
	}
	if q.UMAResolutionStatus != "" {
		vals.Set("uma_resolution_status", q.UMAResolutionStatus)
	}

	var markets []*Market
	if err := c.http.Do(ctx, http.MethodGet, "/markets", vals, nil, nil, &markets); err != nil {
		return nil, err
	}
	return markets, nil
}
