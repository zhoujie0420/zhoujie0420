# Eino 框架 Quick Start 学习笔记

> 来源：[Eino Quick Start](https://www.cloudwego.io/docs/eino/quick_start/) | [GitHub](https://github.com/cloudwego/eino)

## 快速开始

```bash
git clone https://github.com/cloudwego/eino-examples.git
cd eino-examples/quickstart/chatwitheino

# 配置模型（OpenAI 为例）
export OPENAI_API_KEY="your-key"
export OPENAI_MODEL="gpt-4.1-mini"

# 运行第一章示例
go run ./cmd/ch01 -- "用一句话解释 Eino 的 Component 设计解决了什么问题？"
```

---

## 学习路径（共 10 章）

| 章节 | 主题 | 入口 | 核心能力 |
|------|------|------|----------|
| Ch01 | ChatModel 与 Message | `cmd/ch01` | 基础对话 |
| Ch02 | Agent 与 Runner | `cmd/ch02` | 多轮对话 |
| Ch03 | Memory 与 Session | `cmd/ch03` | 持久化 |
| Ch04 | Tool 与文件系统 | `cmd/ch04` | 工具调用 |
| Ch05 | Middleware | `cmd/ch05` | 中间件扩展 |
| Ch06 | Callback | `cmd/ch06` | 可观测性 |
| Ch07 | Interrupt 与 Resume | `cmd/ch07` | 中断恢复 |
| Ch08 | Graph Tool | `cmd/ch08` | 复杂编排 |
| Ch09 | Skill | `cmd/ch09` | 知识复用 |
| 最终章 | A2UI (Web) | `cmd/final` | 生产级应用 |

---

## 第一章：ChatModel 与 Message（Console）

### 核心概念：Component 接口

Eino 定义了一组 Component 接口（ChatModel、Tool、Retriever、Loader 等），每个接口描述一类可替换的能力：

```go
type BaseChatModel interface {
    Generate(ctx context.Context, input []*schema.Message, opts ...Option) (*schema.Message, error)
    Stream(ctx context.Context, input []*schema.Message, opts ...Option) (
        *schema.StreamReader[*schema.Message], error)
}
```

接口的好处：
- **实现可替换**：OpenAI、Ark、Claude、Ollama 等多种实现，切换模型只需改构造逻辑
- **编排可组合**：Agent、Graph、Chain 等编排层只依赖接口
- **测试可 Mock**：接口天然支持 mock

### schema.Message：对话的基本单位

```go
type Message struct {
    Role      RoleType    // system / user / assistant / tool
    Content   string      // 文本内容
    ToolCalls []ToolCall  // 仅 assistant 消息可能有
}

// 常用构造函数
schema.SystemMessage("You are a helpful assistant.")
schema.UserMessage("What is the weather today?")
schema.AssistantMessage("I don't know.", nil)
schema.ToolMessage("tool result", "call_id")
```

### 关键代码

```go
// 构造输入
messages := []*schema.Message{
    schema.SystemMessage(instruction),
    schema.UserMessage(query),
}

// 调用 Stream
stream, err := cm.Stream(ctx, messages)
if err != nil {
    log.Fatal(err)
}
defer stream.Close()

for {
    chunk, err := stream.Recv()
    if errors.Is(err, io.EOF) {
        break
    }
    if err != nil {
        log.Fatal(err)
    }
    fmt.Print(chunk.Content)
}
```


---

## 第二章：Agent 与 Runner（Console 多轮对话）

### Component vs Agent

| 维度 | ChatModel (Component) | ChatModelAgent (Agent) |
|------|----------------------|----------------------|
| 定位 | 组件 | 智能体 |
| 接口 | `Generate()` / `Stream()` | `Run()` -> `AsyncIterator[*AgentEvent]` |
| 输出 | 直接返回消息 | 事件流（消息 + 控制动作） |
| 能力 | 单纯模型调用 | 可扩展 tools、middleware、interrupt |

类比：ChatModel 像"数据库驱动"，ChatModelAgent 像"业务逻辑层"。

### Agent 接口

```go
type Agent interface {
    Name(ctx context.Context) string
    Description(ctx context.Context) string
    Run(ctx context.Context, input *AgentInput, options ...AgentRunOption) *AsyncIterator[*AgentEvent]
}
```

### ChatModelAgent 创建

```go
agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
    Name:        "Ch02ChatModelAgent",
    Description: "A minimal ChatModelAgent with in-memory multi-turn history.",
    Instruction: instruction,
    Model:       cm,
})
```

### Runner：Agent 的执行入口

```go
runner := adk.NewRunner(ctx, adk.RunnerConfig{
    Agent:          agent,
    EnableStreaming: true,
})

// 方式 1：传入消息列表
events := runner.Run(ctx, history)

// 方式 2：便捷方法
events := runner.Query(ctx, "你好")
```

### AgentEvent 事件

```go
type AgentEvent struct {
    AgentName string
    Output    *AgentOutput  // 输出内容
    Action    *AgentAction  // 控制动作（中断/转移/退出）
    Err       error
}
```

### 多轮对话实现

```go
history := make([]*schema.Message, 0, 16)

for {
    line := readUserInput()
    if line == "" {
        break
    }
    
    // 追加用户消息
    history = append(history, schema.UserMessage(line))
    
    // 调用 Runner
    events := runner.Run(ctx, history)
    
    // 消费事件流，收集回复
    content := collectAssistantFromEvents(events)
    
    // 追加 assistant 消息
    history = append(history, schema.AssistantMessage(content, nil))
}
```

---

## 第三章：Memory 与 Session（持久化对话）

> ⚠️ Memory/Session/Store 是**业务层概念**，不是 Eino 框架核心组件。Eino 只负责"如何处理消息"，"如何存储消息"由业务层决定。

### Session

```go
type Session struct {
    ID        string
    CreatedAt time.Time
    messages  []*schema.Message
}
```

核心方法：`Append(msg)`、`GetMessages()`、`Title()`

### Store

```go
type Store struct {
    dir   string
    cache map[string]*Session
}
```

核心方法：`GetOrCreate(id)`、`List()`、`Delete(id)`

### JSONL 存储格式

```jsonl
{"type":"session","id":"083d16da-...","created_at":"2026-03-11T10:00:00Z"}
{"role":"user","content":"你好，我是张三"}
{"role":"assistant","content":"你好张三！很高兴认识你..."}
```

### 运行

```bash
# 创建新会话
go run ./cmd/ch03

# 恢复已有会话
go run ./cmd/ch03 --session <session-id>
```

### 架构分层

```
┌─────────────────────────────────────────┐
│           业务层（你的代码）               │
│  Session(存储) → GetMessages() → runner.Run()
│       ↑                           │
│    Append() ←──── 助手回复 ────────┘
└─────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────┐
│           框架层（Eino）                  │
│  adk.Runner：接收消息列表，调用 ChatModel  │
└─────────────────────────────────────────┘
```


---

## 第四章～第九章 概要

### Ch04：Tool 与文件系统
- 给 Agent 添加工具调用能力（文件读取、代码搜索等）
- ChatModelAgent 内部实现 ReAct 循环：LLM 决定是否调用 Tool → 执行 Tool → 结果返回 LLM

```go
agent, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
    Model: chatModel,
    ToolsConfig: adk.ToolsConfig{
        ToolsNodeConfig: compose.ToolsNodeConfig{
            Tools: []tool.BaseTool{weatherTool, calculatorTool},
        },
    },
})
```

### Ch05：Middleware
- 中间件机制，统一处理横切关注点（重试、审批、错误处理等）

```go
agent, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
    Middlewares: []adk.AgentMiddleware{fsMiddleware},
})
```

### Ch06：Callback
- 回调机制，监控 Agent 执行过程
- 五种回调：OnStart、OnEnd、OnError、OnStartWithStreamInput、OnEndWithStreamOutput

```go
handler := NewHandlerBuilder().
    OnStartFn(func(ctx context.Context, info *RunInfo, input CallbackInput) context.Context {
        log.Infof("onStart, runInfo: %v, input: %v", info, input)
    }).
    OnEndFn(func(ctx context.Context, info *RunInfo, output CallbackOutput) context.Context {
        log.Infof("onEnd, runInfo: %v, out: %v", info, output)
    }).
    Build()

compiledGraph.Invoke(ctx, input, WithCallbacks(handler))
```

### Ch07：Interrupt 与 Resume
- Agent 可以中断等待人工输入，然后从断点恢复
- 需要注册 CheckPointStore

```go
runner := adk.NewRunner(ctx, adk.RunnerConfig{
    Agent:           complexAgent,
    CheckPointStore: memoryStore,
})

// 启动执行
iter := runner.Query(ctx, "recommend a book", adk.WithCheckPointID("1"))

// 从断点恢复
iter, err = runner.Resume(ctx, "1", adk.WithToolOptions(...))
```

### Ch08：Graph Tool（复杂编排）
- 使用 compose 构建 Graph/Chain/Workflow
- Graph 可以暴露为 Tool 供 Agent 调用

```go
// 构建 Graph
graph := compose.NewGraph[*Input, *Output]()
graph.AddLambdaNode("validate", validateFn)
graph.AddChatModelNode("generate", chatModel)
graph.AddLambdaNode("format", formatFn)
graph.AddEdge(compose.START, "validate")
graph.AddEdge("validate", "generate")
graph.AddEdge("generate", "format")
graph.AddEdge("format", compose.END)
runnable, _ := graph.Compile(ctx)

// 暴露为 Tool
tool, _ := graphtool.NewInvokableGraphTool(graph, "data_pipeline", "Process and validate data")
```

### Ch09：Skill（知识复用）
- 将稳定的知识/指令包注入 Agent，按需加载

---

## 编排 API 对比

| API | 特点 | 适用场景 |
|-----|------|----------|
| Chain | 简单链式有向图，只能往前 | 线性流程 |
| Graph | 有向图（支持环），最灵活 | 复杂逻辑、条件分支 |
| Workflow | 有向无环图，struct 字段级数据映射 | 多输入多输出的数据流 |

---

## ADK 预置 Agent 模式

| 模式 | 说明 | 适用场景 |
|------|------|----------|
| ChatModelAgent | ReAct 循环 | 单 Agent + 工具调用 |
| SequentialAgent | 顺序执行子 Agent | CI/CD、ETL |
| ParallelAgent | 并发执行子 Agent | 多源数据采集 |
| LoopAgent | 循环执行子 Agent | 迭代优化、压测 |
| Supervisor | 一个主 Agent 协调多个子 Agent | 项目管理、客服路由 |
| Plan-Execute | 计划→执行→反思→重新计划 | 复杂研究、多步推理 |
| DeepAgent | 任务分解 + 进度追踪 + 上下文隔离 | 多角色协作、长流程 |

---

## 学习资源

- 官方文档：https://www.cloudwego.io/docs/eino/overview/
- GitHub 仓库：https://github.com/cloudwego/eino
- 示例代码：https://github.com/cloudwego/eino-examples
- 扩展组件：https://github.com/cloudwego/eino-ext
- ADK 详解：https://www.cloudwego.io/docs/eino/overview/eino_adk0_1/
- API 文档：https://pkg.go.dev/github.com/cloudwego/eino
