# Rewards

奖励相关接口，位于 `CLOBClient`：

公开接口：

- `GetCurrentRewards`：当前奖励市场列表（自动分页）
- `GetRawRewardsForMarket`：指定市场奖励配置（自动分页）

需要 L2 认证：

- `GetEarningsForUserForDay`：单日收益明细（自动分页）
- `GetTotalEarningsForUserForDay`：单日总收益
- `GetUserEarningsAndMarketsConfig`：单日各市场收益与配置（自动分页）
- `GetRewardPercentages`：各市场奖励占比

