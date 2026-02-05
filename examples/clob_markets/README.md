# CLOB Markets 示例

这个示例演示如何使用 Polymarket SDK 获取 CLOB 服务下的市场数据。

## 功能说明

1. **获取市场列表**：使用 `GetMarkets()` 方法获取 CLOB 市场列表（分页）
2. **处理分页**：演示如何通过 `next_cursor` 遍历所有市场数据
3. **获取单个市场**：使用 `GetMarket()` 方法通过 condition_id 获取单个市场详情

## CLOB Markets vs Gamma API Markets

| API | 服务地址 | 方法 | 用途 |
|-----|---------|------|------|
| **CLOB Markets** | `https://clob.polymarket.com` | `sdk.CLOB.GetMarkets()` | 获取可在链上订单簿交易的市场 |
| **Gamma API Markets** | `https://gamma-api.polymarket.com` | `sdk.REST.Markets()` | 获取市场的基础数据和元信息 |

## 运行方式

### 1. 准备环境变量

创建 `.env` 文件：

```bash
POLYMARKET_API_KEY=your_api_key
POLYMARKET_API_SECRET=your_api_secret
POLYMARKET_PASSPHRASE=your_passphrase
POLYMARKET_ADDRESS=your_wallet_address
PROXY=socks5://127.0.0.1:1080  # 可选
```

### 2. 运行示例

```bash
cd examples/clob_markets
go run main.go
```

### 3. 获取单个市场详情

```bash
go run main.go <condition_id>

# 例如:
go run main.go 0x357dB11734418424652a99459eE2b610E3716997
```

## 输出示例

```
========== 第 1 页 ==========
Limit: 100
Count: 100
NextCursor: MjA=
Data length: 100

  [1] 市场信息:
      Condition ID: 0x357dB11734418424652a99459eE2b610E3716997
      Question: Will Bitcoin reach $100k by end of 2025?
      Ticker: BTC-100K-2025
      Active: true
      Closed: false
      Orders: https://clob.polymarket.com/orders/?...

  [2] 市场信息:
      Condition ID: 0x8a7c56f1b07d4d3e9b5a3c2d1e0f9a8b7c6d5e4f
      Question: Will Ethereum reach $10k by end of 2025?
      Ticker: ETH-10K-2025
      Active: true
      Closed: false

... (还有 97 个市场)

========== 获取单个市场示例 ==========
提示: 可以传入 condition_id 参数获取单个市场详情
```

## 代码说明

### 分页处理

```go
// 使用初始游标开始
nextCursor := pm.InitialCursor

for {
    // 获取数据
    resp, err := sdk.CLOB.GetMarkets(ctx, nextCursor)

    // 处理数据...
    for _, rawMarket := range resp.Data {
        // 解析市场数据
    }

    // 检查是否还有下一页
    if resp.NextCursor == "" || resp.NextCursor == pm.EndCursor || resp.NextCursor == nextCursor {
        break
    }

    nextCursor = resp.NextCursor
}
```

### 获取单个市场

```go
// 通过 condition_id 获取单个市场详情
rawMarket, err := sdk.CLOB.GetMarket(ctx, conditionID)

// 解析 JSON 数据
var market map[string]interface{}
json.Unmarshal(rawMarket, &market)
```

## 数据结构

`GetMarkets()` 返回的 `PaginationPayload` 结构：

```go
type PaginationPayload struct {
    Limit      int               `json:"limit"`      // 每页数量限制
    Count      int               `json:"count"`      // 当前页数据量
    NextCursor string            `json:"next_cursor"` // 下一页游标
    Data       []json.RawMessage `json:"data"`       // 市场数据（JSON）
}
```

## 注意事项

1. **CLOB 市场** 是实际可以在订单簿交易的市场，包含交易相关的字段（如 `orders`）
2. **分页获取** 时建议添加延迟，避免请求过快被限流
3. **初始游标** 建议使用 `pm.InitialCursor`（值为 `"MA=="`）
4. 当 `next_cursor` 为空、等于 `pm.EndCursor`，或与当前相同时，表示没有更多数据
