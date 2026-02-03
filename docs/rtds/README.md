# RTDS

`RTDSClient` 用于订阅实时数据流（ws-live-data.polymarket.com）。

## 常用接口

- `Connect()`：建立连接
- `SubscribeCryptoPrices(source, symbols, handler)`：订阅行情
- `Unsubscribe(topic)`：取消订阅
- `Close()`：关闭连接

