# 异步编排多个下游服务调用 - 面试准备

## 背景

交易查询接口需要同时调用清算服务、风控服务、卡组织查询等多个下游，之前串行调用每个下游 300-500ms，总耗时 1.8s。改为 errgroup 并行编排后，总耗时降至 400ms。

---

## Q：用什么方式实现的异步编排？

有两种方案，最终选的 errgroup。

### 方案一：WaitGroup + channel

最基础的做法，用 WaitGroup 等待所有 goroutine 完成，用 channel 收集错误：

```go
var wg sync.WaitGroup
errCh := make(chan error, 3)

var settleResult *SettleResult
var riskResult *RiskResult
var cardResult *CardResult

wg.Add(3)

go func() {
    defer wg.Done()
    var err error
    settleResult, err = settleService.Query(ctx, txnId)
    if err != nil {
        errCh <- fmt.Errorf("settle query failed: %w", err)
    }
}()

go func() {
    defer wg.Done()
    var err error
    riskResult, err = riskService.Check(ctx, txnId)
    if err != nil {
        errCh <- fmt.Errorf("risk check failed: %w", err)
    }
}()

go func() {
    defer wg.Done()
    var err error
    cardResult, err = cardService.GetInfo(ctx, cardBin)
    if err != nil {
        errCh <- fmt.Errorf("card query failed: %w", err)
    }
}()

wg.Wait()
close(errCh)

// 检查是否有错误
for err := range errCh {
    return nil, err
}
// 聚合结果
```

缺点：
- 需要手动管理 WaitGroup 的 Add/Done，容易漏写导致死锁
- 错误收集需要自己用 channel 实现，代码量多
- 没有 context 取消联动，一个失败了其他的还在跑，浪费资源

### 方案二：errgroup（最终采用）

errgroup 内部封装了 WaitGroup 和错误传播机制，自带 context 取消，更简洁安全：

```go
g, ctx := errgroup.WithContext(ctx)

var settleResult *SettleResult
var riskResult *RiskResult
var cardResult *CardResult

g.Go(func() error {
    var err error
    settleResult, err = settleService.Query(ctx, txnId)
    return err
})

g.Go(func() error {
    var err error
    riskResult, err = riskService.Check(ctx, txnId)
    return err
})

g.Go(func() error {
    var err error
    cardResult, err = cardService.GetInfo(ctx, cardBin)
    return err
})

if err := g.Wait(); err != nil {
    return nil, err
}
// 聚合结果
```

选 errgroup 的原因：
- 不需要手动管理 Add/Done，不会忘写
- 自动收集第一个错误，不需要额外的 channel
- `WithContext` 版本自带取消传播，一个失败其他的立刻收到取消信号，节省资源

---

## Q：多个下游调用之间有没有依赖关系？

大部分没有依赖，清算查询、风控查询、卡信息查询可以完全并行。但有一种情况有依赖，比如需要先拿到卡组织信息再查对应的路由规则。

### 什么是 DAG 编排

DAG（有向无环图）编排的核心思想：把每个下游调用看作图中的一个节点，节点之间的依赖关系就是有向边。没有入边的节点可以立即并行执行，有入边的节点必须等所有前置节点完成后才能执行。

以我们的场景为例，调用依赖关系如下：

```
清算查询 ──────────────────────┐
                                ├──→ 聚合结果
风控查询 ──────────────────────┤
                                │
卡组织查询 ──→ 路由规则查询 ───┤
                                │
商户信息查询 ─────────────────┘
```

- 清算查询、风控查询、卡组织查询、商户信息查询：无依赖，第一阶段并行
- 路由规则查询：依赖卡组织查询的结果，第二阶段执行

### 简单实现：分阶段编排

最直接的做法，按依赖层级分阶段，每个阶段内部用 errgroup 并行：

```go
// 第一阶段：无依赖的调用并行
g1, ctx := errgroup.WithContext(ctx)

var settleResult *SettleResult
var riskResult *RiskResult
var cardResult *CardResult
var merchantResult *MerchantResult

g1.Go(func() error {
    var err error
    settleResult, err = settleService.Query(ctx, txnId)
    return err
})
g1.Go(func() error {
    var err error
    riskResult, err = riskService.Check(ctx, txnId)
    return err
})
g1.Go(func() error {
    var err error
    cardResult, err = cardService.GetInfo(ctx, cardBin)
    return err
})
g1.Go(func() error {
    var err error
    merchantResult, err = merchantService.Query(ctx, merchantId)
    return err
})

if err := g1.Wait(); err != nil {
    return nil, err
}

// 第二阶段：依赖第一阶段结果的调用
g2, ctx := errgroup.WithContext(ctx)

var routeResult *RouteResult

g2.Go(func() error {
    var err error
    // 用第一阶段拿到的 cardResult 去查路由
    routeResult, err = routeService.GetRule(ctx, cardResult.OrgCode)
    return err
})

if err := g2.Wait(); err != nil {
    return nil, err
}

// 聚合所有结果
```

### 进阶实现：通用 DAG 编排器

如果依赖关系更复杂（多层级、多分支），可以抽象成一个通用的 DAG 编排器：

```go
// 定义一个任务节点
type Task struct {
    Name    string
    Deps    []string          // 依赖的任务名
    Run     func(ctx context.Context, deps map[string]interface{}) (interface{}, error)
}

// DAG 编排器
type DAGScheduler struct {
    tasks   map[string]*Task
    results sync.Map
}

func (s *DAGScheduler) AddTask(t *Task) {
    s.tasks[t.Name] = t
}

func (s *DAGScheduler) Execute(ctx context.Context) error {
    // 计算每个任务的入度
    inDegree := make(map[string]int)
    downstream := make(map[string][]string) // 前置任务 -> 后续任务列表

    for name, task := range s.tasks {
        inDegree[name] = len(task.Deps)
        for _, dep := range task.Deps {
            downstream[dep] = append(downstream[dep], name)
        }
    }

    // 用 channel 通知可执行的任务
    ready := make(chan string, len(s.tasks))
    var mu sync.Mutex

    // 入度为 0 的任务立即就绪
    for name, deg := range inDegree {
        if deg == 0 {
            ready <- name
        }
    }

    g, ctx := errgroup.WithContext(ctx)
    remaining := len(s.tasks)

    for remaining > 0 {
        name := <-ready
        remaining--

        taskName := name
        g.Go(func() error {
            task := s.tasks[taskName]

            // 收集依赖结果
            deps := make(map[string]interface{})
            for _, dep := range task.Deps {
                val, _ := s.results.Load(dep)
                deps[dep] = val
            }

            // 执行任务
            result, err := task.Run(ctx, deps)
            if err != nil {
                return fmt.Errorf("task %s failed: %w", taskName, err)
            }
            s.results.Store(taskName, result)

            // 更新下游任务的入度
            mu.Lock()
            for _, next := range downstream[taskName] {
                inDegree[next]--
                if inDegree[next] == 0 {
                    ready <- next
                }
            }
            mu.Unlock()

            return nil
        })
    }

    return g.Wait()
}
```

使用方式：

```go
scheduler := &DAGScheduler{tasks: make(map[string]*Task)}

scheduler.AddTask(&Task{
    Name: "card",
    Deps: nil, // 无依赖
    Run: func(ctx context.Context, deps map[string]interface{}) (interface{}, error) {
        return cardService.GetInfo(ctx, cardBin)
    },
})

scheduler.AddTask(&Task{
    Name: "route",
    Deps: []string{"card"}, // 依赖 card 任务
    Run: func(ctx context.Context, deps map[string]interface{}) (interface{}, error) {
        cardResult := deps["card"].(*CardResult)
        return routeService.GetRule(ctx, cardResult.OrgCode)
    },
})

scheduler.AddTask(&Task{
    Name: "settle",
    Deps: nil,
    Run: func(ctx context.Context, deps map[string]interface{}) (interface{}, error) {
        return settleService.Query(ctx, txnId)
    },
})

err := scheduler.Execute(ctx)
```

### 实际选型

我们业务场景依赖层级最多两层，用的是分阶段编排（简单实现），够用且好维护。通用 DAG 编排器适合依赖关系复杂、经常变化的场景，但引入了额外复杂度，需要权衡。

---

## Q：怎么收集多个 goroutine 的返回结果？

用闭包捕获外部变量。每个 goroutine 写自己的变量，不存在竞争，因为每个变量只有一个 goroutine 在写。`g.Wait()` 返回之后才读取，有 happens-before 保证，不需要加锁。

---

## Q：如果其中一个下游调用失败了怎么办？

分两种策略：

- **核心链路**（如清算结果）：失败就整体失败，errgroup 通过 context 取消其他还在进行的调用，避免资源浪费
- **非核心链路**（如营销信息）：失败走降级，返回默认值，不影响主流程

非核心调用单独起 goroutine，不放进 errgroup，用 recover 兜底：

```go
go func() {
    defer func() {
        if r := recover(); r != nil {
            log.Errorf("marketing query panic: %v", r)
        }
    }()
    marketResult, _ = marketService.Query(ctx, txnId)
}()
```

---

## Q：有没有用 context 做超时控制？超时时间怎么定的？

有。外层接口整体超时 2s，每个下游调用单独设置超时 800ms。超时时间根据下游服务的 P99 耗时来定，从监控平台拉历史数据，取 P99 再加 buffer。

```go
ctx, cancel := context.WithTimeout(ctx, 800*time.Millisecond)
defer cancel()
```

---

## Q：有没有做重试？重试策略是什么？

有，用的是循环队列 + 节点轮转的方式做重试。

### 为什么不用普通重试

普通重试（比如简单 for 循环 + 指数退避）有个问题：如果失败原因是某个下游实例网络超时或者宕机了，重试还是打到同一个实例上，大概率还是超时，等于白重试。

### 循环队列重试的思路

把下游服务的多个实例地址维护在一个循环队列（Ring Buffer）里，每次调用从队列里取当前节点，如果失败了就移动指针到下一个节点重试，这样每次重试打到的是不同的实例：

```go
type RingRetrier struct {
    nodes   []string  // 下游实例地址列表
    current int       // 当前指针位置
    mu      sync.Mutex
}

func (r *RingRetrier) Next() string {
    r.mu.Lock()
    defer r.mu.Unlock()
    node := r.nodes[r.current]
    r.current = (r.current + 1) % len(r.nodes)
    return node
}

func (r *RingRetrier) DoWithRetry(ctx context.Context, maxRetry int, fn func(ctx context.Context, addr string) error) error {
    var lastErr error
    for i := 0; i <= maxRetry; i++ {
        addr := r.Next()
        err := fn(ctx, addr)
        if err == nil {
            return nil
        }
        lastErr = err

        // 只有网络超时和 5xx 才重试，4xx 直接返回
        if !isRetryable(err) {
            return err
        }

        log.Warnf("retry %d/%d, switch to next node, addr=%s, err=%v", i+1, maxRetry, addr, err)
    }
    return fmt.Errorf("all retries exhausted: %w", lastErr)
}

func isRetryable(err error) bool {
    // 网络超时、连接拒绝、5xx 可重试
    if errors.Is(err, context.DeadlineExceeded) {
        return true
    }
    var httpErr *HTTPError
    if errors.As(err, &httpErr) {
        return httpErr.StatusCode >= 500
    }
    return false
}
```

### 整体策略

- 幂等的查询接口：最多重试 1 次，通过循环队列切换到不同实例
- 非幂等的写接口：不重试，直接返回错误让上游决策
- 只有网络超时和 5xx 才重试，4xx 直接返回
- 配合指数退避，首次重试间隔 100ms，避免瞬间打满下游

### 和普通重试的对比

| | 普通重试 | 循环队列重试 |
|---|---|---|
| 重试目标 | 同一个实例 | 自动切换到下一个实例 |
| 单实例故障 | 大概率连续失败 | 绕过故障实例，成功率高 |
| 实现复杂度 | 低 | 稍高，需要维护节点列表和指针 |
| 适用场景 | 偶发抖动 | 实例级故障、网络分区 |

---

## Q：goroutine 泄漏怎么防？

三个手段：

1. 所有下游调用都带 context 超时，不会无限等待
2. errgroup 的 context 取消机制，一个失败其他的收到取消信号
3. 线上用 pprof 监控 goroutine 数量，配了告警阈值

---

## Q：1.8s 是怎么分析出来的？

两个手段配合定位的。

### 1. Prometheus + Grafana 监控大盘

我们用 Prometheus 采集接口级别的耗时指标，通过 Grafana 看板发现这个接口的 P50 在 1.5s，P99 到了 2.3s，明显偏高。

具体做法是在中间件里埋点，用 Histogram 记录每个接口的耗时分布：

```go
var httpDuration = prometheus.NewHistogramVec(
    prometheus.HistogramOpts{
        Name:    "http_request_duration_seconds",
        Help:    "HTTP request duration in seconds",
        Buckets: []float64{0.1, 0.3, 0.5, 0.8, 1.0, 1.5, 2.0, 3.0},
    },
    []string{"method", "path", "status"},
)

func MetricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        duration := time.Since(start).Seconds()
        httpDuration.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(statusCode)).Observe(duration)
    })
}
```

在 Grafana 里用 PromQL 查询分位数：

```promql
# P50
histogram_quantile(0.5, rate(http_request_duration_seconds_bucket{path="/api/txn/query"}[5m]))

# P99
histogram_quantile(0.99, rate(http_request_duration_seconds_bucket{path="/api/txn/query"}[5m]))
```

这一步能发现"哪个接口慢"，但看不出"慢在哪"。

### 2. 链路追踪定位瓶颈

知道接口慢之后，用链路追踪（类似 Jaeger）看具体的调用链。从 trace 上能清楚看到每个下游调用的耗时和串行关系：

```
[txn/query] ─── 1.8s ──────────────────────────────────────────┐
  ├── [settleService.Query]   ── 450ms ──┐                     │
  │                                       │                     │
  ├── [riskService.Check]     ── 320ms ──┐│ (串行，逐个等待)    │
  │                                       ││                    │
  ├── [cardService.GetInfo]   ── 380ms ──┐││                    │
  │                                       │││                   │
  ├── [routeService.GetRule]  ── 350ms ──┐│││                   │
  │                                       ││││                  │
  └── [merchantService.Query] ── 300ms ──┘│││                   │
                                          └┘┘                   │
                                    总计 ≈ 1.8s ────────────────┘
```

一看就明白了：5 个下游串行调用，每个 300-500ms，叠加起来 1.8s。其中大部分没有依赖关系，完全可以并行。

### 优化后的监控数据

并行化上线后，同样从 Prometheus 看板确认效果：

- P50：350ms
- P99：600ms
- 接口超时率从 3.2% 降到 0.1% 以下

---

## Q：并发数怎么控制的？

这个场景下游就 4-5 个，goroutine 数量固定，不需要额外限制。

但在批量处理场景（如批量对账），任务量可能上万，不能无限制开 goroutine，我们手写了一个协程池来控制并发。

### 为什么要控制

- goroutine 虽然轻量（初始栈 2KB），但上万个同时跑会导致调度开销大、内存飙升
- 下游服务有连接数和限流限制，并发太高会被拒绝
- 不控制的话 goroutine 泄漏很难排查

### 手写协程池实现

核心思路：预先创建固定数量的 worker goroutine，通过 channel 分发任务，worker 从 channel 取任务执行：

```go
type Task func()

type Pool struct {
    taskCh   chan Task
    workerN  int
    wg       sync.WaitGroup
}

func NewPool(workerN int, queueSize int) *Pool {
    p := &Pool{
        taskCh:  make(chan Task, queueSize),
        workerN: workerN,
    }
    p.start()
    return p
}

func (p *Pool) start() {
    for i := 0; i < p.workerN; i++ {
        go func() {
            for task := range p.taskCh {
                task()
                p.wg.Done()
            }
        }()
    }
}

// Submit 提交任务，队列满时会阻塞
func (p *Pool) Submit(task Task) {
    p.wg.Add(1)
    p.taskCh <- task
}

// Wait 等待所有任务完成
func (p *Pool) Wait() {
    p.wg.Wait()
}

// Shutdown 关闭协程池
func (p *Pool) Shutdown() {
    close(p.taskCh)
}
```

### 使用示例（批量对账场景）

```go
pool := NewPool(10, 100) // 10个worker，队列容量100

var mu sync.Mutex
var results []*ReconcileResult
var errs []error

for _, item := range reconcileItems { // 可能上万条
    item := item
    pool.Submit(func() {
        result, err := reconcile(ctx, item)
        mu.Lock()
        if err != nil {
            errs = append(errs, err)
        } else {
            results = append(results, result)
        }
        mu.Unlock()
    })
}

pool.Wait()
pool.Shutdown()
```

### 和其他方案的对比

| | 无限制 goroutine | channel 信号量 | 协程池 |
|---|---|---|---|
| 实现复杂度 | 低 | 中 | 稍高 |
| goroutine 数量 | 不可控 | 可控，但每个任务还是新建 goroutine | 固定，复用 worker |
| 内存开销 | 高 | 中 | 低，worker 复用 |
| 适用场景 | 少量并发 | 中等规模 | 大批量任务 |

协程池的优势在于 worker 是复用的，不会反复创建销毁 goroutine，在批量对账这种上万任务的场景下内存更稳定。

---

## Q：如果下游有限流，并行调用会不会触发限流？

会。上线初期风控服务被打触发了限流。后来做了两件事：

1. 跟下游确认限流阈值，确保并发量在安全范围内
2. 调用层加了客户端侧的限流（令牌桶），控制对单个下游的并发请求数

---

## 总结关键词

errgroup 并行编排 → context 超时取消 → 核心/非核心降级策略 → 链路追踪定位瓶颈 → pprof 监控兜底
