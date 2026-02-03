# REST

`RESTClient` 封装 gamma-api（默认 `https://gamma-api.polymarket.com`），用于事件/市场查询。

## 常用接口

- `Events(ctx, q)`：事件列表
- `EventBySlug(ctx, slug, q)`：通过 slug 获取事件
- `Markets(ctx, q)`：市场列表

## 示例

```go
sdk, _ := pm.New(pm.Config{})
event, _ := sdk.REST.EventBySlug(ctx, "some-event-slug", pm.EventBySlugQuery{})
_ = event
```

