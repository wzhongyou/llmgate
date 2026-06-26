# llmgate

> **高性能 LLM 网关，专为 Go 智能体打造** — 统一 API 接入 28 家模型供应商，函数调用、流量控制、自动降级、智能体调试一站解决。

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

[快速开始](#快速开始) · [SDK API](#api) · [网关部署](#server-模式) · [控制台](#控制台) · [设计文档](docs/llmgate-design.md)

```go
// ❌ 没有 llmgate：每家 SDK 不同，切换模型 = 到处改集成代码
resp, _ := deepseek.Chat(ctx, "sk-xxx", deepseekReq)
// → 换成 Claude？得换成 anthropic.Messages(...)，请求格式全变

// ✅ 有 llmgate：一套 API，任意模型，自动降级，内置指标
gw, _ := llmgate.NewFromFile("llmgate.toml")
reply, _ := gw.With("claude").Chat(ctx, req)         // 切换只需改一个字符串
reply, _ := gw.Fallback("claude", "deepseek").Chat(ctx, req) // 自动降级
snap := gw.Snapshot()                                  // 各 provider 延迟及 token 统计
```

---

## 这是什么

你在构建智能体应用，需要接入多个大模型。三个真实问题：

- **换模型要改代码** — 每家 API 格式不同，业务逻辑和模型耦合
- **出了问题不知道为什么** — 是模型慢？被限速？还是 token 超出预算？
- **Token 用量不透明** — 多个 provider 混用，输入/输出/推理 token 各自消耗对不上

**llmgate** 的做法：统一接口屏蔽差异，每次调用白盒记录 provider / 模型 / token 明细 / 延迟，内置降级和延迟限制策略，为可视化监控铺路。

**三种使用形态，按需选择：**

| 形态 | 适用场景 |
|------|----------|
| **SDK** | Go 项目直接引入，一行代码接入多个模型 |
| **Gateway** | 独立部署 HTTP 服务，非 Go 项目也能接入 |
| **Studio** | 可视化控制台，渠道管理 / 对话调试 / Mock 桩 / 请求记录一屏看清 |

### 核心亮点

- **函数调用（Tool Use）全协议支持** — 统一的 `Tool`/`ToolCall` 类型同时适配 OpenAI、Anthropic、Gemini 三套协议，为 ReAct 智能体和 ToolNode 工作流做好准备。
- **数据驱动的供应商架构** — 26 个 OpenAI 兼容 provider 统一为一张 `builtins` 表。新增供应商只需一条配置，无需编写 Go 代码。
- **零代码自定义接入** — 配置文件中 `protocol = "openai-compat"` 即可接入任意 OpenAI 兼容 API，完全不用写代码。

---

## 快速开始

```bash
go get github.com/wzhongyou/llmgate
```

**三种配置方式任选其一：**

**方式一 — 配置文件（推荐）**

```bash
cp llmgate.toml.example llmgate.toml
# 编辑 key 字段
```

```go
gw, err := llmgate.NewFromFile("llmgate.toml")
```

**方式二 — 环境变量**

```bash
export DEEPSEEK_KEY="sk-xxx"
export GLM_KEY="your-glm-key"
```

```go
gw := llmgate.New()  // 自动从环境变量加载
```

**方式三 — 代码**

```go
gw := llmgate.New()
gw.Use("deepseek", "sk-xxx")
```

**开始使用：**

```go
package main

import (
    "context"
    "fmt"

    "github.com/wzhongyou/llmgate"

)

func main() {
    gw, err := llmgate.NewFromFile("llmgate.toml")
    if err != nil {
        panic(err)
    }

    ctx := context.Background()
    reply, err := gw.Chat(ctx, &llmgate.ChatRequest{
        Messages: []llmgate.Message{
            {Role: "user", Content: "帮我写一个 Go HTTP server"},
        },
    })
    if err != nil {
        panic(err)
    }
    fmt.Printf("[%s] %s\n", reply.Provider, reply.Content)
}
```

---

## API

```go
gw := llmgate.New()

// 注册 provider
gw.Use("deepseek", "sk-xxx")
gw.Use("anthropic", "sk-xxx")

// 指定 provider
reply, _ := gw.With("anthropic").Chat(ctx, req)

// 降级链
reply, _ := gw.Fallback("anthropic", "deepseek").Chat(ctx, req)

// 流式输出（SSE）
ch, err := gw.ChatStream(ctx, &llmgate.ChatRequest{
    Messages: []llmgate.Message{{Role: "user", Content: "你好"}},
})
for chunk := range ch {
    if chunk.Error != nil { break }
    fmt.Print(chunk.Content)
}

// 函数调用（Function Calling / Tool Use）
reply, _ := gw.Chat(ctx, &llmgate.ChatRequest{
    Messages: []llmgate.Message{{Role: "user", Content: "北京今天天气怎么样？"}},
    Tools: []llmgate.Tool{{
        Type: "function",
        Function: llmgate.ToolFunction{
            Name:        "get_weather",
            Description: "获取指定城市的当前天气",
            Parameters: map[string]interface{}{
                "type": "object",
                "properties": map[string]interface{}{
                    "city": map[string]interface{}{"type": "string"},
                },
                "required": []string{"city"},
            },
        },
    }},
    ToolChoice: "auto",
})
// reply.ToolCalls 包含模型决定调用的工具
// reply.FinishReason == "tool_calls"

// 指标
snap := gw.Snapshot()
fmt.Printf("DeepSeek 延迟: %.2f ms\n", snap.Providers["deepseek"].AvgLatencyMs)
```

## 控制台

网关二进制内嵌了开发者控制台。配置 `admin_token` 即可启用：

```toml
[server]
listen_addr = ":8080"
admin_token = "your-secret"
```

```bash
go run ./examples/server
# 浏览器打开 http://localhost:8080/admin/
```

**五大页面覆盖开发全流程：**

| 页面 | 用途 |
|------|------|
| **Channels（渠道）** | 可视化管理 provider — 增/删/改/测，查看状态和延迟 |
| **Playground（调试）** | 交互式测试 — 选模型、看流式输出、Token 用量、延迟 |
| **Mock（桩）** | 模拟返回 — 模拟 429/500/超时/空响应，测试异常处理 |
| **Recent（记录）** | 最近 200 条请求 — 完整输入输出查看，每 5s 自动刷新 |
| **Harness（治具）** | 请求录制/回放、影子流量对比、格式合规告警 |

Mock provider 以 `"mock"` 名称注册在引擎上，在 Playground 中选择即可验证规则。

详见 [llmgate-design.md](docs/llmgate-design.md#console)。

---

**优先级（从高到低）：**
1. `.Fallback(...)` — 代码中显式指定降级链
2. `.With(...)` — 固定使用某个 provider
3. `UseStrategy(...)` — 自定义策略
4. 自动检测（llmgate.toml → 环境变量 → 代码）

---

## Server 模式

独立部署 HTTP 服务，多语言接入：

```bash
cp llmgate.toml.example llmgate.toml
go run examples/server/main.go
```

```toml
# llmgate.toml
[[providers]]
name = "glm"
key = "${GLM_KEY}"

[[providers]]
name = "deepseek"
key = "${DEEPSEEK_KEY}"

[strategy]
primary = "glm"
fallback = ["deepseek"]
latency_threshold_ms = 5000

[server]
listen_addr = ":8080"
admin_token = "your-secret"   # 启用 /admin/ 控制台
```

接口：
- `POST /v1/chat` — 对话，支持函数调用（可选 `?provider=` / `?fallback=` 查询参数）
- `GET /v1/models` — 可用模型列表
- `GET /health` — 健康检查
- `GET /metrics` — Prometheus 指标
- `GET /admin/` — 可视化控制台（需配置 `admin_token`）

---

## Harness（测试治具）

内置的 LLM 测试工程工具，解决模型升级验证和质量监控问题。

```toml
[harness]
record = true                         # 录制所有请求/响应到 JSONL
shadow_provider = "deepseek"          # 影子流量：异步发给另一个 provider 对比
```

| 能力 | 说明 |
|------|------|
| **请求录制** | 每次 Chat 调用自动记录到 JSONL 文件，用于回归测试基准 |
| **回放对比** | 历史请求重新发给新模型，对比 finish_reason、tool_calls、token 变化 |
| **影子流量** | 异步复制真实请求到 shadow provider，零延迟影响评估新模型 |
| **合规探针** | 始终检测 tool_calls JSON 合法性和 finish_reason 规范性 |
| **Hook 扩展** | 实现 `core.Hook` 接口，注入自定义观测逻辑 |

自定义 Hook：

```go
type MetricsHook struct{}

func (h *MetricsHook) AfterChat(ctx context.Context, req *core.ChatRequest, resp *core.ChatResponse, err error) {
    // 上报到你的监控系统
}

gw.Engine().AddHook(&MetricsHook{})
```

详见 [llmgate-design.md](docs/llmgate-design.md#harness-engineering)。

---

## 可观测性

每次 `/v1/chat` 请求输出一条结构化 JSON 日志：

```json
{"time":"...","level":"INFO","msg":"request",
 "request_id":"1747123456789","method":"POST","path":"/v1/chat",
 "status":200,"latency_ms":312.5,"remote_addr":"127.0.0.1:54321",
 "provider":"glm","model":"glm-5.1",
 "input_tokens":15,"output_tokens":42,"reasoning_tokens":0}
```

注入自定义 logger：

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
srv, _ := server.New(cfg, server.WithLogger(logger))
```

---

## 支持的供应商

### 全球云平台

| 供应商 | 配置名 | 协议 | 默认模型 | 可用模型 |
|--------|--------|------|----------|----------|
| OpenAI | `openai` | OpenAI 兼容 | `gpt-5.5` | `gpt-5.5`, `gpt-5.4-mini`, `gpt-5-nano`, `gpt-4o` |
| Microsoft Azure Foundry | `azure` | OpenAI 兼容 | `gpt-5.5` | `gpt-5.5`, `gpt-5.4`, `gpt-5.4-mini`, `gpt-5`, `gpt-4o` |
| Google (Gemini) | `gemini` | Gemini 原生 | `gemini-3.1-flash` | 用户配置 |
| Anthropic (Claude) | `anthropic` | Anthropic 原生 | `claude-sonnet-4-6` | 用户配置 |
| Cloudflare AI Gateway | `cloudflare` | OpenAI 兼容 | `@cf/zai-org/glm-5.2` | `glm-5.2`, `kimi-k2.6`, `gemma-4-26b`, `qwen3-30b` |
| Cerebras | `cerebras` | OpenAI 兼容 | `gpt-oss-120b` | `gpt-oss-120b`, `zai-glm-4.7` |
| NVIDIA (NIM) | `nvidia` | OpenAI 兼容 | `nemotron-3-super-120b` | `nemotron-3-super-120b`, `llama-3.3-nemotron-super-49b`, `llama-3.1-70b`, `llama-3.1-8b`, `nemotron-3-nano-30b` |
| Groq | `groq` | OpenAI 兼容 | `llama-4-maverick-17b` | `llama-4-maverick-17b`, `llama-3.3-70b-versatile` |
| Together AI | `together` | OpenAI 兼容 | `Llama-4-Maverick-17B-FP8` | `Llama-4-Maverick-17B-FP8`, `Llama-3.1-70B-Turbo` |

### 中国云平台

| 供应商 | 配置名 | 协议 | 默认模型 | 可用模型 |
|--------|--------|------|----------|----------|
| 阿里云（通义千问） | `qwen` | OpenAI 兼容 | `qwen3.7-max` | `qwen3.7-max`, `qwen3.6-plus`, `qwen3.5-max-preview`, `qwen-max` |
| 字节跳动（豆包） | `doubao` | OpenAI 兼容 | `doubao-seed-2-0-pro` | `seed-2-0-pro`, `seed-2-0-lite`, `seed-1-6`, `seed-1-6-flash`, `seed-1-6-lite` |
| 百度（文心 ERNIE） | `ernie` | OpenAI 兼容 | `ernie-5.1` | `ernie-5.1`, `ernie-4.5` |
| 腾讯（混元） | `hunyuan` | OpenAI 兼容 | `hy3-preview` | `hy3-preview`, `hunyuan-turbos` |
| 讯飞（星火 Spark） | `spark` | OpenAI 兼容 | `spark-x2` | `spark-x2`, `spark-x2-flash`, `spark-ultra`, `spark-pro`, `spark-lite` |
| 商汤（日日新） | `sensenova` | OpenAI 兼容 | `sensenova-6.7-flash-lite` | `sensenova-6.7-flash-lite`, `sensenova-u1-fast` |

### 中国 AI 公司

| 供应商 | 配置名 | 协议 | 默认模型 | 可用模型 |
|--------|--------|------|----------|----------|
| DeepSeek | `deepseek` | OpenAI 兼容 | `deepseek-v4-pro` | `deepseek-v4-pro`, `deepseek-v4-flash` |
| 智谱（GLM） | `glm` | OpenAI 兼容 | `glm-5.2` | `glm-5.2`, `glm-5.1`, `glm-5` |
| 月之暗面（Kimi） | `kimi` | OpenAI 兼容 | `kimi-k2.6` | `kimi-k2.6`, `moonshot-v1-8k`, `moonshot-v1-32k`, `moonshot-v1-128k` |
| MiniMax | `minimax` | OpenAI 兼容 | `MiniMax-M3` | `MiniMax-M3`, `MiniMax-M2.7`, `MiniMax-M2.5` |
| 阶跃星辰（StepFun） | `stepfun` | OpenAI 兼容 | `step-3.7-flash` | `step-3.7-flash`, `step-3.5-flash`, `step-3-mini` |
| 百川智能（Baichuan） | `baichuan` | OpenAI 兼容 | `Baichuan4` | `Baichuan4`, `Baichuan4-Turbo`, `Baichuan4-Air` |
| 零一万物（Yi） | `yi` | OpenAI 兼容 | `yi-lightning` | `yi-lightning`, `yi-large`, `yi-large-turbo`, `yi-medium`, `yi-vision`, `yi-spark` |

### 全球 AI 公司

| 供应商 | 配置名 | 协议 | 默认模型 | 可用模型 |
|--------|--------|------|----------|----------|
| xAI (Grok) | `grok` | OpenAI 兼容 | `grok-4.3` | `grok-4.3`, `grok-4.20`, `grok-4.1`, `grok-code-fast-1` |
| Perplexity | `perplexity` | OpenAI 兼容 | `sonar-pro` | `sonar-pro`, `sonar`, `sonar-reasoning-pro`, `sonar-reasoning`, `sonar-deep-research` |
| Mistral | `mistral` | OpenAI 兼容 | `mistral-large-latest` | `mistral-large-latest`, `mistral-small-latest` |
| Meta (Llama) | `llama` | OpenAI 兼容 | `llama-4-maverick` | `llama-4-maverick`, `llama-4-scout` |

### 推理平台

| 供应商 | 配置名 | 协议 | 默认模型 | 可用模型 |
|--------|--------|------|----------|----------|
| SiliconFlow (硅基流动) | `siliconflow` | OpenAI 兼容 | `Qwen/Qwen3-72B` | `Qwen/Qwen3-72B`, `deepseek-ai/DeepSeek-R1` |
| 小米（MiMo） | `mimo` | OpenAI 兼容 | `mimo-v2-pro` | `mimo-v2.5-pro`, `mimo-v2-pro` |

### 自定义接入

| 配置名 | 协议 | 说明 |
|--------|------|------|
| `openai-compat` | OpenAI 兼容 | 任意 OpenAI 兼容 API，通过配置文件 `protocol = "openai-compat"` 接入 |

> **28 家供应商**（26 家 OpenAI 兼容 + 2 家原生协议），3 套协议全覆盖。所有供应商均支持 `base_url` 覆盖。

```toml
[[providers]]
name = "my-provider"
key = "${MY_PROVIDER_KEY}"
base_url = "https://api.my-provider.com/v1"
protocol = "openai-compat"
```

---

## 项目结构

```
llmgate/
├── core/                 # Provider 接口、引擎、策略、指标
│   └── providers/
│       ├── openaicompat/ # 全部 26 个 OpenAI 兼容 provider（数据驱动）
│       ├── anthropic/    # Anthropic Messages API
│       └── gemini/       # Gemini generateContent API
├── sdk/                  # Go SDK
├── server/               # HTTP 服务
├── docs/                 # llmgate 设计（类型、协议、SDK、网关、扩展）
└── examples/             # 使用示例
├── console/              # 可视化控制台（内嵌）
```

---

## 运行测试

```bash
# 1. 配置 key
cp llmgate.toml.example llmgate.toml
# 填入真实 key，或直接设置环境变量：
# export GLM_KEY=xxx  MINIMAX_KEY=xxx  DEEPSEEK_KEY=xxx

# 2. 运行集成测试
go test ./sdk/ ./server/ -v -count=1
```

未配置 key 时测试自动跳过。

---

## 路线图

- [x] **v0.1** — Go SDK + DeepSeek + 基础降级策略 + 指标采集
- [x] **v0.2** — 智谱（GLM）+ MiniMax + 结构化日志（slog）
- [x] **v0.3** — 14 家供应商 / 3 套协议全覆盖，推理 token 追踪，默认模型可配置
- [x] **v1.0** — Streaming（SSE）+ 生产级路由策略（熔断、限流、重试）
- [x] **v1.1** — 函数调用（Tool Use）全协议支持；28 家供应商；数据驱动架构
- [x] **v1.5** — 可视化控制台：渠道管理、对话调试、Mock 桩、请求记录
- [x] **v2.0** — Harness Engineering：请求录制/回放、影子流量、格式合规探针、Hook 扩展机制

---

## 新增 Provider

**OpenAI 兼容协议** — 在 [core/providers/openaicompat/builtins.go](core/providers/openaicompat/builtins.go) 的 `builtins` 表里加一行：

```go
{
    name:         "myprovider",
    baseURL:      "https://api.myprovider.com/v1",
    defaultModel: "my-model",
    models:       []string{"my-model"},
    envVar:       "MYPROVIDER_KEY",
},
```

**自定义协议** — 实现 `Provider` 接口，参考 [llmgate-design.md](docs/llmgate-design.md#adding-a-provider)。

---

## License

[MIT](LICENSE)
