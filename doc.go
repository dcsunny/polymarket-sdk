// Package polymarket 提供 Polymarket 的 Go SDK。
//
// 入口为 `New(Config)`，返回 `*SDK`，聚合了：
// - REST（gamma-api）
// - CLOB（下单/撤单/订单簿/成交等）
// - WSS（market / user websocket）
// - RTDS（实时行情 websocket）
// - Wallet / Relayer（链上操作与 relayer 提交）
package polymarket
