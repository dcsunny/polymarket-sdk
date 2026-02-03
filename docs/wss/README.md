# WSS

`WSSClient` 用于订阅 CLOB 的 websocket（market / user channel）。

## 常用接口

- `ConnectMarketChannel()`：连接市场频道
- `ConnectUserChannel()`：连接用户频道
- `SubscribeMarketChannel(assetIDs, handlers)`：订阅订单簿等市场事件
- `SubscribeUserChannel(markets, handlers)`：订阅用户事件
- `Close()`：关闭连接

## 示例

参考 `examples/wss_orderbook_by_event`。

