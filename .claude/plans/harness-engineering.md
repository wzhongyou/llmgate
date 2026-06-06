# Harness Engineering 实现计划

## 先提交 ChatResponse json tag 修复

将刚才补全 `ChatResponse` json tag 的改动提交（独立 commit）。

---

## Harness Engineering — P0: 请求录制与回放

### 设计思路

在 Engine 层加入可插拔的 `Hook` 接口（观察者模式），不侵入 provider 调用逻辑。录制/回放/影子流量都通过 Hook 实现。

### 新增文件

1. **`core/hook.go`** — Hook 接口定义 + Engine 注册方法

```go
type Hook interface {
    AfterChat(ctx context.Context, req *ChatRequest, resp *ChatResponse, err error)
}
```

2. **`core/harness/recorder.go`** — 请求录制器，实现 Hook 接口

- 将请求/响应对以 JSONL 格式追加写入文件
- 每条记录包含：timestamp、provider、model、request、response、error、latency_ms
- 文件路径可配置，默认 `harness_records.jsonl`

3. **`core/harness/replay.go`** — 回放器

- 读取 JSONL 文件，逐条回放请求到指定 provider
- 返回对比结果：原始响应 vs 新响应
- 对比维度：finish_reason 一致性、tool_calls 结构一致性、token 用量变化比例

4. **`core/harness/shadow.go`** — 影子流量

- 实现 Hook 接口，AfterChat 时异步将同一请求发给 shadow provider
- shadow 结果写入独立的 JSONL 文件供对比
- 不阻塞主请求

### Engine 改动

在 `core/engine.go` 中：
- Engine 增加 `hooks []Hook` 字段
- 增加 `AddHook(h Hook)` 方法
- `callProvider` 成功后调用所有 hooks 的 `AfterChat`

### 配置

在 `core/config.go` 的 `GatewayConfig` 中增加：

```toml
[harness]
record = true                    # 启用录制
record_path = "harness_records.jsonl"
shadow_provider = ""             # 空则不启用影子流量
shadow_path = "harness_shadow.jsonl"
```

### SDK 层

`sdk/gateway.go` 中 `InitFromConfig` 根据 harness 配置自动注册对应 hook。

### 控制台集成

- 新增 **Harness** tab（第 5 个 tab）
- 功能：
  - 查看录制状态（开/关、已录制条数、文件大小）
  - 触发回放：选择目标 provider，对已录制数据跑回放
  - 查看回放对比结果
  - 配置影子 provider

### API 路由

```
GET    /admin/api/harness/status     # 录制状态
POST   /admin/api/harness/toggle     # 开关录制
POST   /admin/api/harness/replay     # 触发回放（body: {provider, limit})
GET    /admin/api/harness/replay/{id} # 回放结果
```

---

## P1: 格式合规探针

作为另一个 Hook 实现，在 AfterChat 中检测：

- tool_calls 时 arguments 是否为合法 JSON
- arguments 是否符合请求中 tools 定义的 JSON Schema
- finish_reason 是否在已知集合内

违规记录写入 ring buffer，在控制台 Harness tab 展示告警。

---

## 实现顺序

1. 提交 ChatResponse json tag 修复
2. 实现 `core/hook.go` — Hook 接口 + Engine 集成
3. 实现 `core/harness/recorder.go` — 请求录制
4. 实现 `core/harness/replay.go` — 回放对比
5. 实现 `core/harness/shadow.go` — 影子流量
6. 配置层集成（config 解析 + SDK 自动注册）
7. 控制台 Harness tab 前后端
8. 格式合规探针

每一步独立可用，逐步叠加。
