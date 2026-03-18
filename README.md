# ETH Sentinel SDK

以太坊链上事件监控 SDK，提供**嵌入式解析**和 **gRPC 远程服务**两种使用模式。

## 整体架构

```
┌─────────────────────────────────────────────────────────┐
│                      调用方进程                           │
│                                                         │
│  sdk.NewScanBuilder(client)                             │
│      .WithDEX(UniswapV2, Swap)                         │
│      .FromBlock(22000000).ToBlock(22001000)             │
│      .Stream(ctx)                   ← 统一 API          │
│          │                                              │
│    ┌─────┴──────┐                                       │
│    │            │                                       │
│ Embedded    Remote gRPC                                 │
│ (本地解析)   (连接 Server)                               │
└────┼────────────┼────────────────────────────────────── ┘
     │            │
     │     ┌──────▼──────────────────────────────────┐
     │     │         Sentinel gRPC Server             │
     │     │  (:50051)                                │
     │     │                                          │
     │     │  SentinelService                         │
     │     │   ├─ ScanBlock (unary)                   │
     │     │   ├─ ScanBlocks (server stream)          │
     │     │   ├─ SubscribeBlocks (server stream)     │
     │     │   └─ SubscribeWhaleEvents (server stream)│
     │     └───────────────┬──────────────────────────┘
     │                     │
     └──────────┬──────────┘
                │
        ┌───────▼────────────────────────────┐
        │          Scanner                    │
        │  ScanBlock / ScanBlocks             │
        │       │ goroutine pool (ants)        │
        └───────┼────────────────────────────┘
                │  per-block
        ┌───────▼────────────────────────────┐
        │       Chain of Responsibility       │
        │  ┌──────────┐  ┌──────────┐        │
        │  │UniswapV2 │→ │  ERC20   │→ …     │
        │  │  Node    │  │  Node    │        │
        │  └──────────┘  └──────────┘        │
        └───────┬────────────────────────────┘
                │ UnifiedEvent
        ┌───────▼────────────────────────────┐
        │     Business Services               │
        │  WhaleService / ArbitrageService    │
        └────────────────────────────────────┘
```

## 快速开始

### 1. 生成 Proto 代码

```bash
make install-tools   # 首次：安装 protoc 插件
make proto           # 生成 gen/sentinel/v1/*.go
```

### 2. 嵌入模式（直连节点，进程内解析）

```go
import "github.com/ETHSentinel/client"

client, err := sdk.New(
    sdk.WithRPCURL("https://mainnet.infura.io/v3/YOUR_KEY"),
    sdk.WithWSURL("wss://mainnet.infura.io/ws/v3/YOUR_KEY"),
)
defer client.Close()

// 扫描单块
result, err := sdk.NewScanBuilder(client).
    WithDEX(sdk.ProtocolImplUniswapV2, sdk.EventMethodSwap).
    WithToken(sdk.ProtocolImplERC20, sdk.EventMethodTransfer).
    ScanOne(ctx, 22000000)

// 流式扫描区间
ch, err := sdk.NewScanBuilder(client).
    FromBlock(22000000).ToBlock(22001000).
    WithDEX(sdk.ProtocolImplUniswapV2).
    Stream(ctx)

for res := range ch {
    for _, ev := range res.Events {
        if swap, ok := ev.GetDetail().(*sdk.SwapData); ok {
            fmt.Println("Swap:", swap.FromToken.Hex(), "→", swap.ToToken.Hex())
        }
    }
}
```

### 3. 启动独立 gRPC Server

```bash
ETH_RPC_URL=https://mainnet.infura.io/v3/YOUR_KEY \
ETH_WS_URL=wss://mainnet.infura.io/ws/v3/YOUR_KEY \
make run-server
```

### 4. 远程 gRPC 模式

```go
client, err := sdk.NewRemote("localhost:50051")
defer client.Close()

ch, err := sdk.NewScanBuilder(client).
    FromBlock(22000000).ToBlock(22001000).
    WithDEX(sdk.ProtocolImplUniswapV2, sdk.EventMethodSwap).
    Stream(ctx)
```

### 5. 鲸鱼监控

```go
import "github.com/ETHSentinel/internal/service/whale"

whaleSvc := whale.NewService(whale.Config{
    MinUSDValue: decimal.NewFromInt(100_000), // 10万U
}, logger)

// 接入 scanner 事件流
eventCh := make(chan []sdk.Event, 10)
whaleBehaviors := whaleSvc.Run(ctx, eventCh)

for be := range whaleBehaviors {
    fmt.Printf("[鲸鱼] %s 行为=%s 价值=$%s\n",
        be.Address.Hex(), be.EventType, be.ValueUSD)
}
```

## 目录结构

```
eth-sentinel-sdk/
├── proto/sentinel/v1/          # Protobuf 定义
│   ├── events.proto            #   统一事件类型
│   └── sentinel.proto          #   gRPC 服务定义
├── gen/sentinel/v1/            # 生成的 Go 代码（make proto）
├── client/                        # 对外 SDK（调用方 import 这里）
│   ├── client.go               #   Client 接口 + 嵌入/远程两种实现
│   ├── builder.go              #   链式 Fluent Builder
│   ├── option.go               #   Functional Options
│   └── types.go                #   类型别名透出
├── internal/
│   ├── conn/conn.go            #   连接管理（多节点、代理）
│   ├── parser/
│   │   ├── engine.go           #   解析引擎入口
│   │   ├── comm/               #   核心类型 + 接口
│   │   │   ├── types.go
│   │   │   ├── unified.go      #   UnifiedEvent 接口
│   │   │   ├── invoker.go      #   命令模式 Invoker
│   │   │   ├── impl_mgr.go     #   策略模式管理器
│   │   │   └── protocol_mgr.go #   顶层注册表
│   │   ├── dex/uniswap_v2/     #   UniswapV2 解析（filterer 缓存）
│   │   └── token/{erc20,erc721}/
│   ├── scanner/
│   │   ├── scanner.go          #   扫描引擎（goroutine pool）
│   │   ├── chain.go            #   责任链构建
│   │   └── types.go
│   └── service/
│       ├── whale/whale.go      #   鲸鱼监控服务
│       └── arbitrage/          #   套利监控骨架
├── server/grpc/
│   ├── server.go               #   gRPC Server 启动
│   ├── sentinel_handler.go     #   SentinelService 实现
│   └── interceptor.go          #   日志/恢复拦截器
├── cmd/
│   ├── server/main.go          #   启动 gRPC Server
│   └── example/main.go         #   SDK 使用示例
├── Makefile
├── Dockerfile
└── go.mod
```

## 原代码 vs 重构后对比

| 问题 | 原代码 | 重构后 |
|------|--------|--------|
| `ListEvents()` 类型不匹配 | 返回 `map`，`parse_chain.go` for-range 报错 | 返回 `[]EventSig`，`sigIndex` 哈希表 O(1) 查找 |
| `SetEvents` 空实现 | 过滤器完全失效 | `SetFilter` 真实过滤，按 `EventMethod` 语义名筛选 |
| filterer 每次重建 | `ParseSwapEvent` 每调用一次就 `NewUniswappairFilterer` | `ensureCache` 懒加载 + `sync.RWMutex` 缓存复用 |
| 全局 panic | `log.Fatalf` / `panic` 散布各处 | 错误明确 `return error`，仅初始化阶段 ABI 解析 panic |
| 依赖注入缺失 | `conn.go` 硬编码路径，`parser.RegisterAllParser()` 有副作用 | 构造函数注入，`Engine.NewEngine()` 显式初始化 |
| 无 gRPC 层 | REST + 内部直调 | 完整 gRPC 服务（unary + server stream + 拦截器） |
| SDK 入口不清晰 | 无 | `sdk.Client` 接口 + `sdk.NewScanBuilder` 链式 API |



## 扩展新协议

1. 在 `internal/parser/dex/` 下新建目录（如 `curve/`）
2. 实现 `comm.ProtocolImplParser` 接口（参考 `uniswap_v2/parser.go`）
3. 在 `internal/parser/dex/register.go` 的 `RegisterAll` 中添加注册
4. 在 `comm/types.go` 的 `ProtocolImpl` 常量中添加新名称
5. 在 `proto/sentinel/v1/events.proto` 的 `ProtocolImpl` 枚举中添加新值，重新 `make proto`

无需修改 Scanner、责任链、gRPC Handler 的任何代码。

---

### 具体步骤（以添加 Curve 为例）

#### 第一步：新建解析器目录和文件

```
sdk/internal/parser/dex/curve/parser.go   ← 新建这个文件
```

内容模板如下，照着 `uniswap_v2/parser.go` 的结构写：

```go
package curve

import (
    "github.com/ETHSentinel/sdk/internal/parser/comm"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/crypto"
)

// Curve 的 TokenExchange 事件签名
var SigCurveTokenExchange = comm.EventSig(
    crypto.Keccak256Hash([]byte("TokenExchange(address,int128,uint256,int128,uint256)")),
)

type Parser struct {
    invoker *comm.EventParseInvoker
}

func NewParser() *Parser {
    p := &Parser{
        invoker: comm.NewEventParseInvoker(comm.ProtocolImpl("Curve")),
    }
    p.invoker.RegisterOne(SigCurveTokenExchange, comm.EventMethodSwap, p.parseTokenExchange)
    return p
}

func (p *Parser) HandleEvent(sig comm.EventSig, log types.Log, meta comm.EventMetadata) (comm.UnifiedEvent, error) {
    return p.invoker.HandleEvent(sig, log, meta)
}

func (p *Parser) ListEventSigs() []comm.EventSig { return p.invoker.ListEventSigs() }
func (p *Parser) SetFilter(methods []comm.EventMethod) { p.invoker.SetFilter(methods) }

func (p *Parser) parseTokenExchange(log types.Log, meta comm.EventMetadata) (comm.UnifiedEvent, error) {
    // 解析 Curve 的 TokenExchange 事件
    // ...
    return &comm.UnifiedEventData{
        Metadata: comm.EventMetadata{
            ProtocolTypeVal: comm.ProtocolTypeDEX,
            ProtocolImplVal: comm.ProtocolImpl("Curve"),
            // ...
        },
    }, nil
}
```



#### 第二步：在 register.go 加一行注册

打开 sdk/internal/parser/dex/register.go，加一行：

```go
gopackage dex

import (
"github.com/ETHSentinel/sdk/internal/parser/comm"
uniswapv2 "github.com/ETHSentinel/sdk/internal/parser/dex/uniswap_v2"
curve "github.com/ETHSentinel/sdk/internal/parser/dex/curve"  // ← 加这行 import
)

func RegisterAll(mgr *comm.ProtocolImplManager) error {
if err := mgr.RegisterStrategy(comm.ProtocolImplUniswapV2, uniswapv2.NewParser()); err != nil {
return err
}
// ↓ 加这行
if err := mgr.RegisterStrategy(comm.ProtocolImpl("Curve"), curve.NewParser()); err != nil {
return err
}
return nil
}
```

#### 

#### 第三步：使用时指定新协议

在调用 SDK 的地方直接用字符串指定：

```go
ch, err := sdk.NewScanBuilder(client).
    FromBlock(22000000).
    ToBlock(22001000).
    WithDEX(comm.ProtocolImpl("Curve"), comm.EventMethodSwap).  // ← 直接用
    Stream(ctx)
```

---

## 总结

每次加新 DEX 只需动两个地方

```
eth-sentinel-sdk/
└── internal/
        └── parser/
            └── dex/
                ├── register.go           ← 第二步：加一行注册
                ├── uniswap_v2/
                │   └── parser.go         （参考模板）
                └── curve/                ← 第一步：新建这个目录
                    └── parser.go         ← 第一步：新建这个文件
```

**Scanner、责任链、gRPC Handler、SDK Client 全部不需要改**，新协议自动进入责任链参与事件匹配。