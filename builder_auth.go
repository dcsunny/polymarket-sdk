// builder_auth.go 模块
package polymarket

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

// BuilderAuth 签名 Polymarket builder API 请求。
type BuilderAuth struct {
	APIKey     string
	Secret     string
	Passphrase string
}

func NewBuilderAuth(apiKey, secret, passphrase string) *BuilderAuth {
	return &BuilderAuth{
		APIKey:     apiKey,
		Secret:     secret,
		Passphrase: passphrase,
	}
}

// Headers 返回认证头。
func (b *BuilderAuth) Headers(method, path string, body []byte) (map[string]string, error) {
	timestamp := time.Now().Unix()
	signature := b.buildHmacSignature(b.Secret, fmt.Sprintf("%d", timestamp), method, path, string(body))

	return map[string]string{
		"POLY_BUILDER_API_KEY":    b.APIKey,
		"POLY_BUILDER_SIGNATURE":  signature,
		"POLY_BUILDER_TIMESTAMP":  fmt.Sprintf("%d", timestamp),
		"POLY_BUILDER_PASSPHRASE": b.Passphrase,
		"Content-Type":            "application/json",
	}, nil
}

func (b *BuilderAuth) buildHmacSignature(secret, timestamp, method, requestPath, body string) string {
	secretBytes, err := base64.URLEncoding.DecodeString(secret)
	if err != nil {
		secretBytes = []byte(secret)
	}

	message := timestamp + method + requestPath
	if body != "" {
		body = strings.ReplaceAll(body, "'", "\"")
		message += body
	}

	h := hmac.New(sha256.New, secretBytes)
	h.Write([]byte(message))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}
