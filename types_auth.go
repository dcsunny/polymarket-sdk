// types_auth.go 模块
package polymarket

// BanStatus 封禁/限制状态（GET /auth/ban-status/closed-only）。
type BanStatus struct {
	ClosedOnly bool `json:"closed_only"`
}

// ReadonlyAPIKeyResponse 只读 API Key 创建响应。
type ReadonlyAPIKeyResponse struct {
	APIKey string `json:"apiKey"`
}

// BuilderAPIKey Builder API Key（创建返回）。
type BuilderAPIKey struct {
	Key        string `json:"key"`
	Secret     string `json:"secret"`
	Passphrase string `json:"passphrase"`
}

// BuilderAPIKeyResponse Builder API Key 列表项。
type BuilderAPIKeyResponse struct {
	Key       string `json:"key"`
	CreatedAt string `json:"createdAt,omitempty"`
	RevokedAt string `json:"revokedAt,omitempty"`
}
