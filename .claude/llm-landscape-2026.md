# LLM 行业竞品与排行榜调研报告（2026年6月）

---

## 一、行业排行榜总结

### 1.1 LMArena Chatbot Arena 排行榜（2026年中）

| Rank | Model | Elo | Organization |
|------|-------|-----|-------------|
| 1 | Claude Fable 5 | 1,510 | Anthropic |
| 2 | Claude Opus 4.8 Thinking | 1,506 | Anthropic |
| 3 | GPT-5.5 (high) | 1,506 | OpenAI |
| 4 | Claude Opus 4.7 Thinking | 1,505 | Anthropic |
| 5 | Gemini 3.1 Pro | 1,505 | Google DeepMind |
| 6 | Claude Opus 4.8 | 1,504 | Anthropic |
| 7 | Gemini 3.5 Flash | 1,504 | Google DeepMind |
| 8 | Claude Opus 4.7 | 1,503 | Anthropic |
| 9 | Claude Opus 4.6 Thinking | 1,503 | Anthropic |
| 10 | Grok-4.20 | 1,496 | xAI |
| 11 | GPT-5.4 (high) | 1,495 | OpenAI |
| 12 | Gemini 3 Pro | 1,492 | Google DeepMind |
| 13 | Claude Opus 4.6 | 1,490 | Anthropic |
| 14 | GLM-5.2 | 1,488 | Z.ai (Zhipu) |
| 15 | Qwen3.7-Max | 1,486 | Alibaba |

**关键趋势**：头部模型 Elo 极度压缩（前4名在25分内），竞争焦点从纯智能转向成本、延迟、可靠性、领域专项能力。中国模型（GLM-5.2、Qwen3.7-Max、DeepSeek V4 Pro、Kimi K2.6）正在快速追赶。

### 1.2 专项 Benchmark 领先者

| Benchmark | 最佳模型 | 得分 |
|-----------|---------|------|
| **GPQA Diamond** (推理) | Claude Opus 4.7 | 94.2% |
| **SWE-Bench Pro** (编程) | Claude Opus 4.7 | 64.3% |
| **GSM8K** (数学) | GPT-5.5 | 94.2% |
| **Mensa IQ** (视觉推理) | Grok 4.20 / GPT 5.4 Pro | 145 |
| **Legal Drafting** | Claude Opus 4.8 | 67.6% |
| **Multimodal** | Gemini 3.1 Pro | — |

### 1.3 成本效率排行

中国模型在成本效率上优势巨大：
- **DeepSeek V4 Pro**：约为 Claude Opus 的 ~1/22 成本，达到 ~93% 能力
- **Qwen3.6-Plus**：约 $1.25/M tokens
- **GLM-5.2**：约为美国前沿模型的 1/6 成本

---

## 二、当前项目已支持的 Provider（21个）

| # | Provider | 类型 | API协议 |
|---|----------|------|---------|
| 1 | openai | OpenAI-compat | OpenAI |
| 2 | anthropic | Native | Anthropic |
| 3 | gemini | Native | Gemini |
| 4 | deepseek | OpenAI-compat | OpenAI |
| 5 | qwen (阿里云) | OpenAI-compat | OpenAI |
| 6 | doubao (字节/火山) | OpenAI-compat | OpenAI |
| 7 | ernie (百度) | OpenAI-compat | OpenAI |
| 8 | hunyuan (腾讯) | OpenAI-compat | OpenAI |
| 9 | glm (智谱) | OpenAI-compat | OpenAI |
| 10 | kimi (月之暗面) | OpenAI-compat | OpenAI |
| 11 | minimax | OpenAI-compat | OpenAI |
| 12 | stepfun (阶跃星辰) | OpenAI-compat | OpenAI |
| 13 | baichuan (百川) | OpenAI-compat | OpenAI |
| 14 | mistral | OpenAI-compat | OpenAI |
| 15 | grok (xAI) | OpenAI-compat | OpenAI |
| 16 | groq | OpenAI-compat | OpenAI |
| 17 | together | OpenAI-compat | OpenAI |
| 18 | llama (Meta) | OpenAI-compat | OpenAI |
| 19 | mimo (小米) | OpenAI-compat | OpenAI |
| 20 | siliconflow (硅基流动) | OpenAI-compat | OpenAI |
| 21 | openai-compat (generic) | OpenAI-compat | OpenAI |

---

## 三、竞品平台覆盖分析

### 3.1 国内云厂商

| 云厂商 | AI服务 | 核心模型 | 已在 llmgate？ | API兼容性 |
|--------|--------|---------|:---:|-----------|
| **阿里云** | DashScope/百炼 | Qwen3.7-Max, Qwen3.6-Plus, Qwen3-Coder, Qwen2.5-VL | ✅ qwen | OpenAI |
| **火山引擎** | ARK | Doubao-Seed-1.6, Doubao-Pro/Lite | ✅ doubao | OpenAI |
| **百度云** | 千帆 | Ernie-5.1, Ernie-4.5 | ✅ ernie | OpenAI |
| **腾讯云** | 混元 | HY3-Preview, Hunyuan-TurboS | ✅ hunyuan | OpenAI |
| **华为云** | ModelArts/盘古 | Pangu-NLP-N4-MoE-Reasoner, Pangu-NLP-N4-128K | ❌ 待添加 | OpenAI? |
| **讯飞** | 星火/星辰MaaS | Spark X2, Spark X2 Flash | ❌ 待添加 | OpenAI |
| **京东云** | 言犀 | ChatRhino | ❌ 待调研 | 需确认 |

### 3.2 国外云厂商

| 云厂商 | AI服务 | 核心自有模型 | 已在 llmgate？ | API兼容性 |
|--------|--------|------------|:---:|-----------|
| **Google Cloud** | Vertex AI | Gemini 3.1 Pro/Flash | ✅ gemini | Gemini/OpenAI |
| **AWS** | Bedrock | Nova 2 Sonic/Lite/Premier, Titan | ❌ 待添加 | AWS SDK (非OpenAI) |
| **Microsoft Azure** | Foundry (原Azure OpenAI) | MAI系列, 代理GPT/Claude等 | ❌ 待添加 | OpenAI |
| **Oracle Cloud** | OCI Generative AI | (代理Meta/Cohere/xAI等) | ❌ 待添加 | OpenAI-compat |
| **IBM Cloud** | watsonx.ai | Granite 4.0 系列 | ❌ 待添加 | 需确认 |
| **Cloudflare** | Workers AI | (代理GPT-OSS/Kimi/GLM/Qwen) | ❌ 待添加 | OpenAI-compat |

### 3.3 专业AI平台

| 平台 | 核心模型/特色 | 已在 llmgate？ | 备注 |
|------|-------------|:---:|------|
| **Perplexity** | Sonar (搜索增强), Sonar-Pro | ❌ 待添加 | 以实时搜索为核心差异 |
| **NVIDIA NIM** | Nemotron-3-Ultra/Super/Nano | ❌ 待添加 | 自有芯片+代理多厂商 |
| **Cerebras** | 最快推理 (Llama/GPT-OSS/Qwen) | ❌ 待添加 | 1K-3K tok/s |
| **Replicate** | 代理30+模型社区 | ❌ 待添加 | 已被Cloudflare收购 |
| **Together AI** | Llama, 开源模型托管 | ✅ together | — |

### 3.4 国内AI公司（非云厂）

| 公司 | 模型 | 已在 llmgate？ | API兼容性 |
|------|------|:---:|-----------|
| **智谱(Z.ai)** | GLM-5.2, GLM-5.1, GLM-4.7 | ✅ glm | OpenAI |
| **月之暗面(Moonshot)** | Kimi-K2.6, Kimi-K2.5 | ✅ kimi | OpenAI |
| **MiniMax** | M2.7, M3 | ✅ minimax | OpenAI |
| **阶跃星辰(StepFun)** | Step-3.7-Flash, Step-3.5 | ✅ stepfun | OpenAI |
| **百川智能(Baichuan)** | Baichuan4, Baichuan4-Air | ✅ baichuan | OpenAI |
| **DeepSeek** | V4-Pro, V4-Flash, R2 | ✅ deepseek | OpenAI |
| **商汤(SenseTime)** | SenseNova 6.7 Flash-Lite | ❌ 待添加 | OpenAI |
| **零一万物(01.AI)** | Yi-Lightning, Yi-Large | ❌ 待添加 | OpenAI |
| **昆仑万维** | SkyClaw-v1.0, 天工4.0 | ❌ 待添加 | OpenAI |
| **科大讯飞(iFlytek)** | Spark X2, X2 Flash | ❌ 待添加 | OpenAI (星辰MaaS) |
| **小米(Mimo)** | Mimo-V2.5-Pro | ✅ mimo | OpenAI |

---

## 四、建议新增 Provider 优先级

### P0 — 国内云厂商（补齐全覆盖）

| Provider | 关键理由 | API协议 | 工作量估计 |
|----------|---------|---------|-----------|
| **hunyuan** | ✅ 已有 | — | — |
| **qwen** | ✅ 已有 | — | — |
| **doubao** | ✅ 已有 | — | — |
| **ernie** | ✅ 已有 | — | — |
| **huawei** (华为云盘古) | 国内TOP5云厂商，唯一缺失 | 需确认协议（类似/modelarts） | 中 |
| **spark** (讯飞星火) | 国内重要AI平台，星辰MaaS OpenAI兼容 | OpenAI | 低 |
| **sensenova** (商汤) | 多模态领先，OpenAI兼容 | OpenAI | 低 |
| **yi** (零一万物) | 李开复团队，性价比突出，OpenAI兼容 | OpenAI | 低 |
| **skywork** (昆仑万维天工) | SkyClaw Agent模型，OpenAI兼容 | OpenAI | 低 |

### P1 — 国外云厂商（战略全覆盖）

| Provider | 关键理由 | API协议 | 工作量估计 |
|----------|---------|---------|-----------|
| **azure** (微软Foundry) | 全球第二大云，企业首选 | OpenAI (兼容) | 低 |
| **bedrock** (AWS) | 全球第一大云，110+模型托管 | AWS SDK (Converse API) 或 OpenAI-compat proxy | 高 |
| **oci** (Oracle) | 企业级，Cohere/Meta/xAI等代理 | OpenAI-compat (部分) | 中 |
| **watsonx** (IBM) | 企业级，Granite自研+第三方 | 需确认 | 中 |
| **cloudflare** (Workers AI) | 边缘推理，多模型代理 | OpenAI-compat | 低 |

### P2 — 专业AI平台

| Provider | 关键理由 | API协议 | 工作量估计 |
|----------|---------|---------|-----------|
| **perplexity** (搜索增强) | 独特搜索+推理能力 | OpenAI | 低 |
| **nvidia** (NIM) | GPU生态整合，80+模型 | OpenAI | 低 |
| **cerebras** (最快推理) | 速度差异化(1-3K tok/s) | OpenAI | 低 |

---

## 五、各 Provider 核心模型速查

### 华为云盘古
- `Pangu-NLP-N4-MoE-Reasoner-32K-5.0.0.1` — 旗舰MoE推理模型
- `Pangu-NLP-N4-Reasoner-128K-3.0.1.2` — 长上下文推理
- `Pangu-NLP-N4-32K-2.5.35` — 通用NLP

### 讯飞星火
- `Spark X2` — 旗舰深度推理模型
- `Spark X2 Flash` — 轻量推理
- 星辰MaaS平台支持 GLM-5.1, DeepSeek V4, Kimi K2.6, Qwen3.5 等第三方

### 商汤日日新
- `SenseNova 6.7 Flash-Lite` — 原生多模态智能体模型
- 兼容 OpenAI API，`https://api.sensenova.cn/v1`

### 零一万物 Yi
- `Yi-Lightning` — 高性价比MoE (~200B)，128K上下文
- `Yi-Large-Preview` — 大规模预览版，256K上下文
- 兼容 OpenAI API，`https://api.lingyiwanwu.com/v1`

### 昆仑万维天工
- `SkyClaw-v1.0` — 100万Token上下文Agent原生模型
- `SkyClaw-v1.0-lite` — 轻量版
- 兼容 OpenAI API

### AWS Bedrock (最复杂)
- 18个Provider，110+模型
- 自有模型: Nova 2 Sonic/Lite/Premier, Nova Pro, Titan Embeddings
- 第三方: Claude 4.x, Llama 4, Mistral Large 3, DeepSeek V3.2, Qwen3, Cohere, AI21
- API: AWS SDK (非原生OpenAI)，可通过Converse API或bedrock-access-gateway转OpenAI格式

### Azure Foundry
- GPT-5.5/5.4/5.3系列 + Claude Opus/Sonnet/Haiku 4.x + 第三方(Grok, DeepSeek, Kimi, Mistral, Cohere等)
- MAI自有模型: MAI-Image-2.5, MAI-Voice-2, MAI-Transcribe-1.5
- 兼容 OpenAI API

### Oracle OCI
- 代理: Meta Llama 4/3.3, Cohere Command A, OpenAI gpt-oss, xAI Grok
- 部分模型兼容OpenAI协议

### IBM watsonx
- 自有: Granite 4.0 系列 (h-small/tiny/micro), Granite 3.x, Granite代码/视觉/语音/安全
- 第三方: Llama 4, Mistral Large 3, gpt-oss-120b, DeepSeek R1

### Cloudflare Workers AI
- 代理: Llama 4/3.3, Kimi K2.6/K2.7, GLM-5.2, Gemma 4, Qwen3, gpt-oss, Nemotron
- 兼容 OpenAI API

### Perplexity
- Sonar-Pro, Sonar, Sonar-Reasoning-Pro, Sonar-Deep-Research
- 核心差异: 实时联网搜索+来源引用
- 兼容 OpenAI API

### NVIDIA NIM
- 自有: Nemotron-3-Ultra 550B, Nemotron-3-Super 120B, Nemotron-3-Nano, Llama-Nemotron
- 代理: Llama 4, Mistral Large 3, DeepSeek V4, Qwen3.5, Kimi K2.6, gpt-oss, Gemma 4
- 兼容 OpenAI API

### Cerebras
- 超快推理: GPT-OSS 120B (~3,000 tok/s), Qwen3 32B (~2,600 tok/s)
- 兼容 OpenAI API

---

## 六、实施建议

1. **短期（本次迭代）**：通过 `openaicompat/builtins.go` 添加 6 个国内 Provider：
   - `huawei`（需先确认盘古API是否原生兼容OpenAI协议，可能需要自定义适配器）
   - `spark`
   - `sensenova`
   - `yi`
   - `skywork`
   - `perplexity`

2. **中期（下一迭代）**：添加 Azure Foundry + Cloudflare Workers AI（均为OpenAI兼容，低工作量）

3. **长期**：AWS Bedrock（需专门的AWS SDK适配器）和 OCI/IBM watsonx
