# Changelog

All notable changes to llmgate are documented in this file.

---

## [v0.10.0] — 2026-06-26

### Added
- **8 家新增供应商**：Azure Foundry、Cerebras、Cloudflare AI Gateway、NVIDIA NIM、Perplexity、商汤（SenseNova）、讯飞（星火 Spark）、零一万物（Yi）
- `AuthHeader` 字段支持非标准认证头（Azure `api-key`）
- Provider Coverage 文档章节（设计文档）
- 供应商总数：20 → **28 家**

### Changed
- **刷新全部 28 家供应商模型列表至 2026.06 最新**（DeepSeek、OpenAI、GLM、Grok、豆包、MiniMax、通义千问等）
- 移除不存在的模型（`deepseek-r2`、`gpt-5.5-high`、`grok-4.20-beta1` 等）
- 修正错误模型名（`dola-seed-2.0-pro` → `doubao-seed-2-0-pro-260215` 等）
- README 供应商表格重构为分组 + 完整模型清单

---

## [v0.9.0] — 2026-06-06

### Added
- **Shadow 模型粒度支持**：`shadow_model` 配置项，可为影子流量指定不同于原始请求的模型

### Fixed
- Replay 跳过 Recorder，避免雪球录制
- Replay 请求记录到 Recent，带 `source=replay` 标识
- Server Config 缺少 Harness 字段导致录制配置不生效

---

## [v0.8.0] — 2026-06-06

### Added
- **Harness Engineering 子系统**：请求录制/回放、影子流量、格式合规探针、Hook 扩展机制
  - `core/hook.go` — Hook 接口 + Engine 集成
  - `core/harness/recorder.go` — JSONL 请求录制
  - `core/harness/replay.go` — 请求回放与对比
  - `core/harness/shadow.go` — 异步影子流量
  - `core/harness/probe.go` — 格式合规探针（始终激活）
- Console: Harness 管理界面（录制开关、回放执行、违规查看）
- `gw.Recorder()` / `gw.Shadow()` / `gw.Probe()` SDK API

---

## [v0.7.0] — 2026-05-31

### Added
- `ChatResponse` / `StreamChunk` JSON tag 补全，HTTP 响应符合 OpenAI 协议
- **模型列表升级至 2026.06 最新**
- 补齐通用 LLM SDK 核心能力：结构化输出、多模态、Embeddings、生成参数

### Fixed
- Console token 显示修复
- Recent 记录修复

---

## [v0.6.0] — 2026-05-17

### Added
- GLM thinking 参数适配（`thinking` → `enable_thinking`）
- `reasoning_content` 透传
- **可视化控制台（Console）**：渠道管理、对话调试、Mock 桩、请求记录
  - 嵌入式 Web UI（`console/static/` embed.FS）
  - Admin API（Provider CRUD、Playground、Mock 规则、Recent 日志、Config 持久化）
  - 实时指标 + 自动刷新
- **Provider 自动注册机制**：`init()` 自注册，零配置接入新供应商

---

## [v0.5.0] — 2026-05-16

### Added
- **Tool Use（函数调用）全协议支持**
  - 统一 `Tool` / `ToolCall` 类型适配 OpenAI、Anthropic、Gemini 三套协议
  - `ToolChoice` 支持 `auto` / `none` / `required`
  - 流式 Tool Calls 累积（SSE stream parser）
- **数据驱动供应商架构**
  - 18 个独立 provider 包精简为一张 `builtins` 表
  - 新增 OpenAI 兼容供应商只需一条配置，无需编写 Go 代码
  - `openai-compat` 通用工厂：配置文件中 `protocol = "openai-compat"` 即可接入

---

## [v0.4.0] — 2026-05-16

### Added
- **SSE Streaming 支持**（`ChatStream` interface + OpenAI SSE parser）
- **生产级路由策略**
  - 熔断器（Circuit Breaker）：per-provider，连续失败触发
  - 指数退避 + 抖动重试（Exponential backoff + jitter）
  - 结构化错误类型 `ProviderError`（含 Retryable 标记）

---

## [v0.3.0] — 2026-05-10

### Added
- **14 家供应商 / 3 套协议全覆盖**
  - OpenAI-compatible: openai, deepseek, qwen, doubao, ernie, glm, grok, groq, hunyuan, kimi, minimax, mistral, together
  - Anthropic Messages: anthropic
  - Gemini generateContent: gemini
- 推理 token 追踪（`reasoning_tokens` in Usage）
- 默认模型可配置
- 模块重命名：`llmgateway` → `llmgate`

---

## [v0.2.0] — 2026-05-10

### Added
- 智谱（GLM）+ MiniMax 供应商接入
- 结构化日志（slog）

---

## [v0.1.0] — 2026-05-10

### Added
- Go SDK + DeepSeek 供应商接入
- 基础降级策略（PrimaryFirst、Latency、TimeBased）
- 指标采集（MetricsSnapshot、ProviderStats）
- 配置文件驱动（TOML） + 环境变量注入
- HTTP Server 模式（`POST /v1/chat`、`GET /v1/models`）
