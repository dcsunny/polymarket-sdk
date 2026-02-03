// l2_hmac.go 模块
package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

// L2Signature 为 CLOB L2 认证生成 HMAC-SHA256 签名。
// message = timestamp + method + requestPath + body
func L2Signature(secret, timestamp, method, requestPath, body string) string {
	decoded, err := base64.URLEncoding.DecodeString(secret)
	if err != nil {
		decoded, _ = base64.StdEncoding.DecodeString(secret)
	}

	msg := timestamp + method + requestPath + body
	mac := hmac.New(sha256.New, decoded)
	mac.Write([]byte(msg))
	return base64.URLEncoding.EncodeToString(mac.Sum(nil))
}
