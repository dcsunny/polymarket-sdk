# Examples

`examples/` 是一个独立的 Go module（有自己的 `go.mod`），用于演示 `polymarket-sdk` 的常见用法。

## 依赖与环境变量

- 使用 `github.com/joho/godotenv` 自动加载当前目录下的 `.env`
- 每个示例目录内都有一份 `.env` 示例文件（可按需补齐）

## 运行

推荐在每个示例目录内执行（这样会自动加载该目录下的 `.env`）：

```bash
cd wss_orderbook_by_event && go run .
cd ../clob_trades && go run .
cd ../clob_place_order && go run .
cd ../clob_balance_allowance && go run .
cd ../clob_open_orders && go run .
cd ../clob_order_history && go run .
cd ../relayer_deploy_safe && go run .
cd ../relayer_redeem_tokens && go run .
```

### wss_orderbook_by_event

用于实时监听某个市场的订单簿（WSS market channel）。

需要设置其一：

- `EVENT_SLUG`：事件 slug（会先走 REST `EventBySlug` 自动解析一个 `asset_id`）
- `ASSET_ID`：直接指定 CLOB 的 `asset_id`

### clob_trades

查询某个市场的成交记录（CLOB trades，L2 认证）。

- `POLYMARKET_MARKET_ID`
- `POLYMARKET_ADDRESS`
- `POLYMARKET_API_KEY`
- `POLYMARKET_API_SECRET`
- `POLYMARKET_PASSPHRASE`

### clob_place_order

下单示例（CLOB 下单，L2 + 签名相关配置）。

必需/常用：

- `POLYMARKET_PRIVATE_KEY`
- `POLYMARKET_ADDRESS`
- `POLYMARKET_FUNDER`
- `POLYMARKET_SIG_TYPE`（默认 `2`，Safe）
- `POLYMARKET_CHAIN_ID`（默认 `137`）
- `POLYMARKET_API_KEY`
- `POLYMARKET_API_SECRET`
- `POLYMARKET_PASSPHRASE`
- `POLYMARKET_TOKEN_ID`

### clob_balance_allowance

获取账户余额与授权（`/balance-allowance`）。

常用：

- `POLYMARKET_ADDRESS`
- `POLYMARKET_API_KEY`
- `POLYMARKET_API_SECRET`
- `POLYMARKET_PASSPHRASE`

可选：

- `POLYMARKET_SIG_TYPE`（`0`=EOA, `1`=Proxy, `2`=Safe；需与账户类型一致）
- `POLYMARKET_TOKEN_ID`（查询某个 conditional token 的 balance/allowance）

### clob_open_orders

获取账户在某个市场/资产上的活跃订单（open orders）。

可选：

- `POLYMARKET_MARKET_ID`
- `POLYMARKET_ASSET_ID`

### clob_order_history

通过历史成交（trades）反推订单历史（只覆盖“发生过成交”的订单）。

可选：

- `POLYMARKET_MARKET_ID`
- `MAX_PAGES`（默认 2）
- `MAX_ORDERS`（默认 20）

### relayer_redeem_tokens

使用 relayer 提交赎回（redeem positions）。

必填：

- `POLYMARKET_RPC_URL`
- `POLYMARKET_PRIVATE_KEY`
- `POLYMARKET_CONDITION_ID`（或 `POLYMARKET_MARKET_ID`）

可选：

- `POLYMARKET_RELAYER_URL`（默认 `https://relayer-v2.polymarket.com`）
- `AUTO_DEPLOY_SAFE`（默认 false）
- NegRisk：`IS_NEGRISK=true` + `YES_AMOUNT`/`NO_AMOUNT`
- CTF：`IS_NEGRISK=false` + `INDEX_SETS` + `COLLATERAL_TOKEN`
- Builder：`POLYMARKET_BUILDER_API_KEY`/`POLYMARKET_BUILDER_API_SECRET`/`POLYMARKET_BUILDER_PASSPHRASE`

### relayer_deploy_safe

部署 Safe（relayer `/deploy-safe`）。

必填：

- `POLYMARKET_RPC_URL`
- `POLYMARKET_PRIVATE_KEY`
- `POLYMARKET_BUILDER_API_KEY` / `POLYMARKET_BUILDER_API_SECRET` / `POLYMARKET_BUILDER_PASSPHRASE`

可选：

- `POLYMARKET_RELAYER_URL`（默认 `https://relayer-v2.polymarket.com`）
- `POLYMARKET_CHAIN_ID`（默认 137）
