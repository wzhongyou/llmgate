package openaicompat

import "github.com/wzhongyou/llmgate/core"

type providerDef struct {
	name         string
	baseURL      string
	defaultModel string
	models       []string
	envVar       string
	authHeader   string // empty = "Authorization: Bearer {key}"; set to "api-key" etc.
	bodyHook     func(map[string]interface{})
}

var builtins = []providerDef{
	{
		name:         "azure",
		baseURL:      "https://PLACEHOLDER.openai.azure.com/openai/v1",
		defaultModel: "gpt-5.5",
		models:       []string{"gpt-5.5", "gpt-5.4", "gpt-5.4-mini", "gpt-5", "gpt-4o"},
		envVar:       "AZURE_KEY",
		authHeader:   "api-key",
	},
	{
		name:         "baichuan",
		baseURL:      "https://api.baichuan-ai.com/v1",
		defaultModel: "Baichuan4",
		models:       []string{"Baichuan4", "Baichuan4-Turbo", "Baichuan4-Air"},
		envVar:       "BAICHUAN_KEY",
	},
	{
		name:         "cerebras",
		baseURL:      "https://api.cerebras.ai/v1",
		defaultModel: "gpt-oss-120b",
		models:       []string{"gpt-oss-120b", "zai-glm-4.7"},
		envVar:       "CEREBRAS_KEY",
	},
	{
		name:         "cloudflare",
		baseURL:      "https://gateway.ai.cloudflare.com/v1/PLACEHOLDER/PLACEHOLDER",
		defaultModel: "@cf/zai-org/glm-5.2",
		models:       []string{"@cf/zai-org/glm-5.2", "@cf/moonshotai/kimi-k2.6", "@cf/google/gemma-4-26b-a4b-it", "@cf/qwen/qwen3-30b-a3b-fp8"},
		envVar:       "CLOUDFLARE_KEY",
	},
	{
		name:         "deepseek",
		baseURL:      "https://api.deepseek.com/v1",
		defaultModel: "deepseek-v4-pro",
		models:       []string{"deepseek-v4-pro", "deepseek-v4-flash"},
		envVar:       "DEEPSEEK_KEY",
	},
	{
		name:         "doubao",
		baseURL:      "https://ark.cn-beijing.volces.com/api/v3",
		defaultModel: "doubao-seed-2-0-pro-260215",
		models:       []string{"doubao-seed-2-0-pro-260215", "doubao-seed-2-0-lite-260215", "doubao-seed-1-6-251015", "doubao-seed-1-6-flash-250828", "doubao-seed-1-6-lite-251015"},
		envVar:       "DOUBAO_KEY",
	},
	{
		name:         "ernie",
		baseURL:      "https://qianfan.baidubce.com/v2",
		defaultModel: "ernie-5.1",
		models:       []string{"ernie-5.1", "ernie-4.5"},
		envVar:       "ERNIE_KEY",
	},
	{
		name:         "glm",
		baseURL:      "https://open.bigmodel.cn/api/paas/v4",
		defaultModel: "glm-5.2",
		models:       []string{"glm-5.2", "glm-5.1", "glm-5"},
		envVar:       "GLM_KEY",
			bodyHook:     glmBodyHook,
		},
	{
		name:         "grok",
		baseURL:      "https://api.x.ai/v1",
		defaultModel: "grok-4.3",
		models:       []string{"grok-4.3", "grok-4.20-0309-non-reasoning", "grok-4.1", "grok-code-fast-1"},
		envVar:       "GROK_KEY",
	},
	{
		name:         "groq",
		baseURL:      "https://api.groq.com/openai/v1",
		defaultModel: "meta-llama/llama-4-maverick-17b-128e",
		models:       []string{"meta-llama/llama-4-maverick-17b-128e", "llama-3.3-70b-versatile"},
		envVar:       "GROQ_KEY",
	},
	{
		name:         "hunyuan",
		baseURL:      "https://api.hunyuan.cloud.tencent.com/v1",
		defaultModel: "hy3-preview",
		models:       []string{"hy3-preview", "hunyuan-turbos"},
		envVar:       "HUNYUAN_KEY",
	},
	{
		name:         "kimi",
		baseURL:      "https://api.moonshot.cn/v1",
		defaultModel: "kimi-k2.6",
		models:       []string{"kimi-k2.6", "moonshot-v1-8k", "moonshot-v1-32k", "moonshot-v1-128k"},
		envVar:       "KIMI_KEY",
	},
	{
		name:         "llama",
		baseURL:      "https://api.llama.com/v1",
		defaultModel: "llama-4-maverick",
		models:       []string{"llama-4-maverick", "llama-4-scout"},
		envVar:       "LLAMA_KEY",
	},
	{
		name:         "mimo",
		baseURL:      "https://api.xiaomimimo.com/v1",
		defaultModel: "mimo-v2-pro",
		models:       []string{"mimo-v2.5-pro", "mimo-v2-pro"},
		envVar:       "MIMO_KEY",
	},
	{
		name:         "minimax",
		baseURL:      "https://api.minimaxi.com/v1",
		defaultModel: "MiniMax-M3",
		models:       []string{"MiniMax-M3", "MiniMax-M2.7", "MiniMax-M2.5"},
		envVar:       "MINIMAX_KEY",
	},
	{
		name:         "mistral",
		baseURL:      "https://api.mistral.ai/v1",
		defaultModel: "mistral-large-latest",
		models:       []string{"mistral-large-latest", "mistral-small-latest"},
		envVar:       "MISTRAL_KEY",
	},
	{
		name:         "nvidia",
		baseURL:      "https://integrate.api.nvidia.com/v1",
		defaultModel: "nvidia/nemotron-3-super-120b-a12b",
		models:       []string{"nvidia/nemotron-3-super-120b-a12b", "nvidia/llama-3.3-nemotron-super-49b-v1.5", "meta/llama-3.1-70b-instruct", "meta/llama-3.1-8b-instruct", "nvidia/nemotron-3-nano-30b-a3b"},
		envVar:       "NVIDIA_KEY",
	},
	{
		name:         "openai",
		baseURL:      "https://api.openai.com/v1",
		defaultModel: "gpt-5.5",
		models:       []string{"gpt-5.5", "gpt-5.4-mini", "gpt-5-nano", "gpt-4o"},
		envVar:       "OPENAI_KEY",
	},
	{
		name:         "perplexity",
		baseURL:      "https://api.perplexity.ai",
		defaultModel: "sonar-pro",
		models:       []string{"sonar-pro", "sonar", "sonar-reasoning-pro", "sonar-reasoning", "sonar-deep-research"},
		envVar:       "PERPLEXITY_KEY",
	},
	{
		name:         "qwen",
		baseURL:      "https://dashscope.aliyuncs.com/compatible-mode/v1",
		defaultModel: "qwen3.7-max",
		models:       []string{"qwen3.7-max", "qwen3.6-plus", "qwen3.5-max-preview", "qwen-max"},
		envVar:       "QWEN_KEY",
	},
	{
		name:         "sensenova",
		baseURL:      "https://api.sensenova.cn/v1",
		defaultModel: "sensenova-6.7-flash-lite",
		models:       []string{"sensenova-6.7-flash-lite", "sensenova-u1-fast"},
		envVar:       "SENSENOVA_KEY",
	},
	{
		name:         "siliconflow",
		baseURL:      "https://api.siliconflow.cn/v1",
		defaultModel: "Qwen/Qwen3-72B",
		models:       []string{"Qwen/Qwen3-72B", "deepseek-ai/DeepSeek-R1"},
		envVar:       "SILICONFLOW_KEY",
	},
	{
		name:         "spark",
		baseURL:      "https://spark-api-open.xf-yun.com/v1",
		defaultModel: "spark-x2",
		models:       []string{"spark-x2", "spark-x2-flash", "spark-ultra", "spark-pro", "spark-lite"},
		envVar:       "SPARK_KEY",
	},
	{
		name:         "stepfun",
		baseURL:      "https://api.stepfun.com/v1",
		defaultModel: "step-3.7-flash",
		models:       []string{"step-3.7-flash", "step-3.5-flash", "step-3-mini"},
		envVar:       "STEPFUN_KEY",
	},
	{
		name:         "together",
		baseURL:      "https://api.together.xyz/v1",
		defaultModel: "meta-llama/Llama-4-Maverick-17B-128E-Instruct-FP8",
		models:       []string{"meta-llama/Llama-4-Maverick-17B-128E-Instruct-FP8", "meta-llama/Meta-Llama-3.1-70B-Instruct-Turbo"},
		envVar:       "TOGETHER_KEY",
	},
	{
		name:         "yi",
		baseURL:      "https://api.lingyiwanwu.com/v1",
		defaultModel: "yi-lightning",
		models:       []string{"yi-lightning", "yi-large", "yi-large-turbo", "yi-medium", "yi-vision", "yi-spark"},
		envVar:       "YI_KEY",
	},
}

func init() {
	for _, def := range builtins {
		d := def
		core.RegisterProvider(d.name, func(cfg core.ProviderConfig) (core.Provider, error) {
			baseURL := cfg.BaseURL
			if baseURL == "" {
				baseURL = d.baseURL
			}
			defaultModel := cfg.DefaultModel
			if defaultModel == "" {
				defaultModel = d.defaultModel
			}
			return &Provider{
				name:         d.name,
				key:          cfg.Key,
				baseURL:      baseURL,
				defaultModel: defaultModel,
				models:       d.models,
				AuthHeader:   d.authHeader,
				BodyHook:     d.bodyHook,
			}, nil
		})
		core.RegisterProviderEnv(d.envVar, d.name)
	}

	// Generic factory for user-defined OpenAI-compatible providers via config protocol = "openai-compat"
	core.RegisterProvider("openai-compat", func(cfg core.ProviderConfig) (core.Provider, error) {
		return &Provider{
			name:         cfg.Name,
			key:          cfg.Key,
			baseURL:      cfg.BaseURL,
			defaultModel: cfg.DefaultModel,
		}, nil
	})
}

// glmBodyHook adapts request body for GLM API.
// GLM does not support OpenAI's "thinking" param; it uses "enable_thinking" instead.
func glmBodyHook(body map[string]interface{}) {
	if _, ok := body["thinking"]; ok {
		delete(body, "thinking")
		body["enable_thinking"] = false
	}
}
