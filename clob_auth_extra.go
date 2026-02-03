// clob_auth_extra.go 模块
package polymarket

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

// GetClosedOnlyMode 获取 closed-only 模式（L2 认证）。
func (c *CLOBClient) GetClosedOnlyMode(ctx context.Context) (*BanStatus, error) {
	path := "/auth/ban-status/closed-only"
	headers, err := c.l2Headers(http.MethodGet, path, "")
	if err != nil {
		return nil, err
	}

	var resp BanStatus
	if err := c.http.Do(ctx, http.MethodGet, path, nil, nil, headers, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CreateReadonlyAPIKey 创建只读 API Key（L2 认证）。
func (c *CLOBClient) CreateReadonlyAPIKey(ctx context.Context) (*ReadonlyAPIKeyResponse, error) {
	path := "/auth/readonly-api-key"
	headers, err := c.l2Headers(http.MethodPost, path, "")
	if err != nil {
		return nil, err
	}

	var resp ReadonlyAPIKeyResponse
	if err := c.http.Do(ctx, http.MethodPost, path, nil, nil, headers, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetReadonlyAPIKeys 获取当前账号的只读 API Key 列表（L2 认证）。
func (c *CLOBClient) GetReadonlyAPIKeys(ctx context.Context) ([]string, error) {
	path := "/auth/readonly-api-keys"
	headers, err := c.l2Headers(http.MethodGet, path, "")
	if err != nil {
		return nil, err
	}

	var resp []string
	if err := c.http.Do(ctx, http.MethodGet, path, nil, nil, headers, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteReadonlyAPIKey 删除只读 API Key（L2 认证）。
func (c *CLOBClient) DeleteReadonlyAPIKey(ctx context.Context, key string) (bool, error) {
	if key == "" {
		return false, ErrInvalidArgument("key is required")
	}
	path := "/auth/readonly-api-key"
	payload := map[string]string{"key": key}
	body, err := json.Marshal(payload)
	if err != nil {
		return false, err
	}
	headers, err := c.l2Headers(http.MethodDelete, path, string(body))
	if err != nil {
		return false, err
	}

	var resp bool
	if err := c.http.DoRaw(ctx, http.MethodDelete, path, nil, body, headers, &resp); err != nil {
		return false, err
	}
	return resp, nil
}

// ValidateReadonlyAPIKey 校验只读 API Key（无需认证）。
func (c *CLOBClient) ValidateReadonlyAPIKey(ctx context.Context, address, key string) (string, error) {
	if address == "" || key == "" {
		return "", ErrInvalidArgument("address and key are required")
	}
	path := "/auth/validate-readonly-api-key"
	vals := url.Values{}
	vals.Set("address", address)
	vals.Set("key", key)

	var resp string
	if err := c.http.Do(ctx, http.MethodGet, path, vals, nil, nil, &resp); err != nil {
		return "", err
	}
	return resp, nil
}

// CreateBuilderAPIKey 创建 builder API Key（L2 认证）。
func (c *CLOBClient) CreateBuilderAPIKey(ctx context.Context) (*BuilderAPIKey, error) {
	path := "/auth/builder-api-key"
	headers, err := c.l2Headers(http.MethodPost, path, "")
	if err != nil {
		return nil, err
	}

	var resp BuilderAPIKey
	if err := c.http.Do(ctx, http.MethodPost, path, nil, nil, headers, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetBuilderAPIKeys 获取 builder API Key 列表（L2 认证）。
func (c *CLOBClient) GetBuilderAPIKeys(ctx context.Context) ([]BuilderAPIKeyResponse, error) {
	path := "/auth/builder-api-key"
	headers, err := c.l2Headers(http.MethodGet, path, "")
	if err != nil {
		return nil, err
	}

	var resp []BuilderAPIKeyResponse
	if err := c.http.Do(ctx, http.MethodGet, path, nil, nil, headers, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// RevokeBuilderAPIKey 撤销 builder API Key（builder auth）。
func (c *CLOBClient) RevokeBuilderAPIKey(ctx context.Context) error {
	if c.builderAuth == nil {
		return ErrInvalidArgument("builder auth is not configured")
	}
	path := "/auth/builder-api-key"

	headers, err := c.builderAuth.Headers(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	return c.http.Do(ctx, http.MethodDelete, path, nil, nil, headers, nil)
}
