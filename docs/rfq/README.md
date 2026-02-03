# RFQ

RFQ（Request For Quote）相关接口，位于 `CLOBClient`：

- `CreateRfqRequest` / `CancelRfqRequest` / `GetRfqRequests`
- `CreateRfqQuote` / `CancelRfqQuote`
- `GetRfqRequesterQuotes` / `GetRfqQuoterQuotes` / `GetRfqBestQuote`
- `GetRfqConfig`
- `AcceptRfqQuote` / `ApproveRfqOrder`（payload 透传）

说明：

- RFQ 接口需要 L2 认证
- 请求体里 `userType` 会自动使用 `Config.SignatureType`

