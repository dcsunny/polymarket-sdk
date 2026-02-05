// client.go 模块
package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	ierr "github.com/dcsunny/polymarket-sdk/internal/errors"
)

// Client 是一个带有基础 URL 和默认值的轻量级 HTTP 封装。
type Client struct {
	baseURL string
	http    *http.Client
	headers map[string]string
	debug   bool
}

// New 创建新的 HTTP 客户端。
func New(baseURL string, timeout time.Duration, proxy string, userAgent string, debug bool) (*Client, error) {
	if baseURL == "" {
		return nil, errors.New("baseURL is required")
	}

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	if proxy != "" {
		if proxyURL, err := url.Parse(proxy); err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	}

	h := map[string]string{}
	if userAgent != "" {
		h["User-Agent"] = userAgent
	}

	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		http: &http.Client{
			Timeout:   timeout,
			Transport: transport,
		},
		headers: h,
		debug:   debug,
	}, nil
}

// Do 发送 JSON 请求。
func (c *Client) Do(ctx context.Context, method, path string, query url.Values, body any, headers map[string]string, out any) error {
	var payload []byte
	var err error
	if body != nil {
		switch v := body.(type) {
		case []byte:
			payload = v
		case string:
			payload = []byte(v)
		default:
			payload, err = json.Marshal(body)
			if err != nil {
				return err
			}
		}
	}
	return c.DoRaw(ctx, method, path, query, payload, headers, out)
}

// DoRaw 发送带有原始字节主体的请求。
func (c *Client) DoRaw(ctx context.Context, method, path string, query url.Values, body []byte, headers map[string]string, out any) error {
	reqURL, err := c.resolveURL(path, query)
	if err != nil {
		return err
	}

	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}
	fmt.Println(reqURL)
	req, err := http.NewRequestWithContext(ctx, method, reqURL, reader)
	if err != nil {
		return err
	}

	for k, v := range c.headers {
		if v != "" {
			req.Header.Set(k, v)
		}
	}
	for k, v := range headers {
		if v != "" {
			req.Header.Set(k, v)
		}
	}
	if body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return parseAPIError(resp, respBytes)
	}

	if out == nil {
		return nil
	}

	switch v := out.(type) {
	case *[]byte:
		*v = respBytes
		return nil
	default:
		if len(respBytes) == 0 {
			return nil
		}
		return json.Unmarshal(respBytes, out)
	}
}

func (c *Client) resolveURL(path string, query url.Values) (string, error) {
	base, err := url.Parse(c.baseURL)
	if err != nil {
		return "", err
	}

	if path != "" {
		ref := &url.URL{Path: path}
		base = base.ResolveReference(ref)
	}

	q := base.Query()
	for k, vals := range query {
		for _, v := range vals {
			if v != "" {
				q.Add(k, v)
			}
		}
	}
	base.RawQuery = q.Encode()
	return base.String(), nil
}

func parseAPIError(resp *http.Response, body []byte) error {
	errObj := struct {
		Error     string `json:"error"`
		Message   string `json:"message"`
		Code      string `json:"code"`
		RequestID string `json:"request_id"`
	}{
		RequestID: resp.Header.Get("X-Request-Id"),
	}
	_ = json.Unmarshal(body, &errObj)

	msg := errObj.Message
	if msg == "" {
		msg = errObj.Error
	}

	return &ierr.APIError{
		Status:    resp.StatusCode,
		Code:      errObj.Code,
		Message:   msg,
		RequestID: errObj.RequestID,
		Body:      string(body),
	}
}
