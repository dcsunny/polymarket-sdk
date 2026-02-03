# Polymarket SDK Implementation Plan

目标
提供一个简单、易用、性能可靠的 Go SDK，入口为 `pm.New()`，覆盖 REST、CLOB、WSS、RTDS、Safe/Proxy/Relayer 钱包能力。使用标准库为主，外部依赖最少且可解释。

设计原则
- 一个入口对象 `SDK` 聚合各子客户端
- 统一配置 `Config`，默认值合理，显式传参可覆盖
- 业务逻辑与传输层分离，内部包隐藏复杂度
- 类型清晰、错误明确、函数简洁
- 合理复用连接与缓存，保证性能与资源可控

公开 API 形态
- `pm.New(cfg)` 返回 `*SDK`
- `SDK.REST` 访问 REST API
- `SDK.CLOB` 访问 CLOB API
- `SDK.WSS` 访问 WebSocket 订阅
- `SDK.RTDS` 访问 RTDS 数据流
- `SDK.Wallet` 访问 Safe/Proxy/Relayer

建议目录结构
```
polymarket-sdk/
  go.mod
  README.md
  docs/
    IMPLEMENTATION_PLAN.md
    ARCHITECTURE.md
  examples/
    rest_markets.go
    clob_orders.go
    wss_market.go
    rtds_prices.go
    wallet_safe.go
    wallet_proxy.go
    relayer_redeem.go
  internal/
    auth/
      l1_signer.go
      l2_hmac.go
      eip712.go
    httpx/
      client.go
      request.go
      retry.go
    errors/
      api_error.go
    cache/
      ttl.go
  config.go
  client.go
  rest.go
  clob.go
  wss.go
  rtds.go
  wallet_safe.go
  wallet_proxy.go
  wallet_relayer.go
  types_rest.go
  types_clob.go
  types_ws.go
  types_wallet.go
```

依赖策略
- 首选标准库 `net/http`, `crypto/*`, `encoding/json`, `context`
- 必要依赖
1. 以太坊签名与 ABI 解析建议使用 `github.com/ethereum/go-ethereum`
2. WebSocket 建议使用 `nhooyr.io/websocket` 或 `golang.org/x/net/websocket`
- 除上述以外不新增依赖，避免冗余

核心配置设计
- `Config` 内含 `BaseURL`, `CLOBBaseURL`, `Timeout`, `Proxy`, `Debug`, `ChainID`
- 认证字段 `Address`, `PrivateKey`, `APIKey`, `APISecret`, `Passphrase`
- `New()` 内部校验和派生

实现逻辑步骤
1. 定义公共类型与错误模型
2. 实现 HTTP 传输层封装
3. 实现认证模块 L1/L2
4. 实现 REST 客户端
5. 实现 CLOB 客户端
6. 实现 WSS 客户端
7. 实现 RTDS 客户端
8. 实现 Safe/Proxy/Relayer 钱包
9. 补充示例与文档
10. 添加测试与性能验证

实现步骤细节

1. 定义公共类型与错误模型
- `types_rest.go` 定义 Event、Market、Pricing 等 REST 类型
- `types_clob.go` 定义 Order、Trade、Book 等 CLOB 类型
- `types_ws.go` 定义 WSS/RTDS 事件类型
- `errors.go` 定义可读错误类型，包含 `Code`, `Message`, `HTTPStatus`, `RequestID`
- `internal/errors/api_error.go` 解析 API 错误响应为统一结构

2. 实现 HTTP 传输层封装
- 使用 `http.Client` + 自定义 `Transport` 保持连接复用
- 支持 `Timeout`、`Proxy`、`UserAgent`
- `internal/httpx` 提供 `Do(ctx, method, url, body, headers)` 封装
- 仅在必要处增加重试，默认对幂等 GET 进行有限次重试

3. 实现认证模块 L1/L2
- L2 HMAC 使用标准库 `crypto/hmac` + `sha256`
- L1 EIP-712 依赖 go-ethereum
- 提供 `Signer` 接口，保证 CLOB 与 Wallet 可以复用
- 内部缓存地址与签名相关中间值，减少重复计算

4. 实现 REST 客户端
- `REST` 结构持有 `httpx.Client`
- 提供 `Events(ctx, EventsQuery)`、`Markets(ctx, MarketsQuery)`、`Pricing(ctx, PricingQuery)`
- 参数尽量使用结构体，序列化时自动忽略空字段
- 解析时间字段使用 `time.Parse` 统一处理

5. 实现 CLOB 客户端
- `CLOB` 结构持有 `httpx.Client` 与 `auth.Signer`
- 支持 L1 创建/派生 API Key
- 支持 L2 请求签名与自动注入头
- 内置必要缓存，如 tickSize、negRisk、feeRate
- 订单相关结构体简洁化，必要时支持扩展字段

6. 实现 WSS 客户端
- 统一连接管理与重连策略
- 提供 `ConnectMarketChannel()` 与 `ConnectUserChannel()`
- 订阅调用使用 `SubscribeMarket(ids, handlers)`、`SubscribeUser(ids, handlers)`
- 事件分发走 map[string]Handler
- 支持心跳与断线重连

7. 实现 RTDS 客户端
- 连接 `ws-live-data.polymarket.com`
- 提供 `SubscribeCryptoPrices(source, symbols, handler)`
- 支持自动重连与订阅重放

8. 实现 Safe/Proxy/Relayer 钱包
- Safe 与 Proxy 逻辑复用签名与交易构造
- `Wallet` 子模块提供简洁函数 `Split`, `Merge`, `Redeem`, `Convert`
- Relayer 支持 `/nonce` `/deployed` `/submit` 流程
- 支持 `BuilderAuth` 作为可选 header

9. 补充示例与文档
- 以 `examples/` 形式覆盖 REST / CLOB / WSS / RTDS / Wallet
- README 重点展示 3 个最短路径示例
- `docs/ARCHITECTURE.md` 说明模块分工

10. 测试与性能验证
- 单测覆盖 L2 签名、URL 序列化、API 错误解析
- 使用 `httptest` 模拟 REST/CLOB
- WSS/RTDS 使用本地 WS 服务器模拟
- 提供基准测试验证签名与序列化性能

性能与简洁性建议
- HTTP 连接复用，禁用不必要的 gzip 解压时避免开销
- JSON 解析时使用 `Decoder` 流式解析大列表
- 保持函数短小，避免多个嵌套层级
- 默认不做全量缓存，只在 CLOB 必要字段做小型缓存

里程碑建议
1. M1: Config + HTTPX + Errors + REST
2. M2: L1/L2 Auth + CLOB
3. M3: WSS + RTDS
4. M4: Wallet + Relayer
5. M5: 文档 + 示例 + 测试

