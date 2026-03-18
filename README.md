# ETH Sentinel SDK

An Ethereum on-chain event monitoring SDK that supports both **embedded parsing** and **gRPC remote service** modes, with multi-chain, Mempool monitoring, and HD wallet capabilities.

---

## Architecture

```
┌──────────────────────────────────────────────────────────┐
│                      Caller Process                       │
│                                                          │
│  sentinel.NewScanBuilder(client)                         │
│      .WithDEX(UniswapV2, Swap)                          │
│      .FromBlock(22000000).ToBlock(22001000)              │
│      .Stream(ctx)                  ← Unified API         │
│              │                                           │
│       ┌──────┴──────┐                                    │
│       │             │                                    │
│   Embedded      Remote gRPC                              │
│  (in-process)  (connect to Server)                       │
└───────┼─────────────┼────────────────────────────────────┘
        │             │
        │      ┌──────▼─────────────────────────────────┐
        │      │        Sentinel gRPC Server             │
        │      │  (:50051)                               │
        │      │                                         │
        │      │  SentinelService                        │
        │      │   ├─ ScanBlock        (unary)           │
        │      │   ├─ ScanBlocks       (server stream)   │
        │      │   ├─ SubscribeBlocks  (server stream)   │
        │      │  MempoolService                         │
        │      │   └─ SubscribePending (server stream)   │
        │      │  WalletService                          │
        │      │   ├─ CreateWallet                       │
        │      │   ├─ GetBalance                         │
        │      │   └─ SendTransaction                    │
        │      └──────────────┬──────────────────────────┘
        │                     │
        └──────────┬──────────┘
                   │
         ┌─────────▼──────────────────────┐
         │           Scanner               │
         │   ScanBlock / ScanBlocks        │
         │      goroutine pool (ants)      │
         └─────────┬──────────────────────┘
                   │  per-block
         ┌─────────▼──────────────────────┐
         │    Chain of Responsibility      │
         │  ┌──────────┐  ┌──────────┐    │
         │  │UniswapV2 │→ │  ERC20   │→ … │
         │  └──────────┘  └──────────┘    │
         └─────────┬──────────────────────┘
                   │  UnifiedEvent
         ┌─────────▼──────────────────────┐
         │       Business Services         │
         │  WhaleService / ArbitrageService│
         └────────────────────────────────┘
```

---

## Project Structure

```
eth-sentinel-final/
├── internal/
│   ├── chain/                   # Multi-chain abstraction (ETH/BSC/Polygon/Arbitrum)
│   │   ├── chain.go             #   Chain interface + BaseChain
│   │   ├── chains.go            #   Built-in chain configs
│   │   └── registry.go          #   Global chain registry
│   ├── conn/
│   │   ├── conn.go              #   Node connection manager
│   │   └── pool.go              #   Multi-chain connection pool
│   ├── parser/
│   │   ├── engine.go            #   Parser engine entry point
│   │   ├── comm/                #   Core types + interfaces
│   │   │   ├── types.go         #     Protocol/Event type definitions
│   │   │   ├── unified.go       #     UnifiedEvent interface
│   │   │   ├── invoker.go       #     Command pattern invoker
│   │   │   ├── impl_mgr.go      #     Strategy pattern manager
│   │   │   └── protocol_mgr.go  #     Top-level registry
│   │   ├── dex/uniswap_v2/      #   UniswapV2 parser (filterer cache)
│   │   └── token/{erc20,erc721}/
│   ├── scanner/
│   │   ├── scanner.go           #   Block scanning engine (goroutine pool)
│   │   ├── chain.go             #   Chain of responsibility builder
│   │   └── types.go
│   ├── mempool/
│   │   ├── monitor.go           #   Mempool monitor (pending tx subscription)
│   │   └── filter.go            #   Built-in filters
│   ├── wallet/
│   │   ├── wallet.go            #   Unified wallet manager
│   │   ├── hd/hdwallet.go       #   BIP44 HD wallet
│   │   ├── keystore/            #   Encrypted local key storage
│   │   └── transaction/         #   Transaction builder + sender
│   └── service/
│       ├── whale/whale.go       #   Whale monitoring service
│       └── arbitrage/           #   Arbitrage detection (skeleton)
├── client/                      # Public SDK (callers import this)
│   ├── client.go                #   Client interface + embedded/remote implementations
│   ├── builder.go               #   Fluent ScanBuilder
│   ├── mempool.go               #   MempoolClient
│   ├── wallet.go                #   WalletClient
│   └── option.go                #   Functional options
├── server/
│   ├── grpc/
│   │   ├── server.go            #   gRPC server startup
│   │   ├── sentinel_handler.go  #   SentinelService implementation
│   │   ├── mempool_handler.go   #   MempoolService implementation
│   │   ├── wallet_handler.go    #   WalletService implementation
│   │   └── interceptor.go       #   Logging + recovery interceptors
│   └── proto/sentinel/v1/
│       ├── events.proto         #   Unified event types
│       ├── sentinel.proto       #   SentinelService definition
│       ├── mempool.proto        #   MempoolService definition
│       └── wallet.proto         #   WalletService definition
├── gen/sentinel/v1/             # Generated protobuf code (run make proto)
├── test/
│   ├── example/main.go          #   Full SDK usage example
│   └── server/main.go           #   gRPC server startup entry
├── Makefile
├── Dockerfile
└── go.mod
```

---

## Quick Start

### Step 1 — Generate Proto Code

```bash
# First time: install protoc plugins
make install-tools

# Generate gen/sentinel/v1/*.go
make proto
```

On Windows (without make):
```powershell
protoc --proto_path=server/proto --go_out=gen --go_opt=paths=source_relative `
       --go-grpc_out=gen --go-grpc_opt=paths=source_relative `
       server/proto/sentinel/v1/events.proto `
       server/proto/sentinel/v1/sentinel.proto `
       server/proto/sentinel/v1/mempool.proto `
       server/proto/sentinel/v1/wallet.proto
```

### Step 2 — Install Dependencies

```bash
go mod tidy
```

### Step 3 — Run the Server

```bash
# Linux / macOS
ETH_RPC_URL=https://mainnet.infura.io/v3/YOUR_KEY \
ETH_WS_URL=wss://mainnet.infura.io/ws/v3/YOUR_KEY \
go run ./test/server
```

```powershell
# Windows PowerShell
$env:ETH_RPC_URL="https://mainnet.infura.io/v3/YOUR_KEY"
$env:ETH_WS_URL="wss://mainnet.infura.io/ws/v3/YOUR_KEY"
go run ./test/server
```

### Step 4 — Run the Example

```powershell
go run ./test/example
```

---

## SDK Usage

### Embedded Mode — Connect directly to an Ethereum node

```go
import sentinel "github.com/ETHSentinel/client"

client, err := sentinel.New(
    sentinel.WithChainID(uint64(sentinel.ChainETH)),
    sentinel.WithRPCURL("https://mainnet.infura.io/v3/YOUR_KEY"),
    sentinel.WithWSURL("wss://mainnet.infura.io/ws/v3/YOUR_KEY"),
    sentinel.WithWorkerPoolSize(10),
)
defer client.Close()
```

### Scan a Single Block

```go
result, err := sentinel.NewScanBuilder(client).
    WithDEX(sentinel.ProtocolImplUniswapV2, sentinel.EventMethodSwap).
    WithToken(sentinel.ProtocolImplERC20, sentinel.EventMethodTransfer).
    ScanOne(ctx, 22000000)

for _, ev := range result.Events {
    if swap, ok := ev.GetDetail().(*sentinel.SwapData); ok {
        fmt.Println("Swap:", swap.FromToken.Hex(), "→", swap.ToToken.Hex())
    }
}
```

### Stream a Block Range

```go
ch, err := sentinel.NewScanBuilder(client).
    FromBlock(22000000).
    ToBlock(22001000).
    WithDEX(sentinel.ProtocolImplUniswapV2).
    Stream(ctx)

for res := range ch {
    fmt.Printf("Block #%s | txs=%d | events=%d\n",
        res.BlockNumber, res.TxCount, len(res.Events))
}
```

### Remote gRPC Mode — Connect to a deployed Sentinel Server

```go
client, err := sentinel.NewRemote("localhost:50051")
defer client.Close()

ch, err := sentinel.NewScanBuilder(client).
    FromBlock(22000000).ToBlock(22001000).
    WithDEX(sentinel.ProtocolImplUniswapV2, sentinel.EventMethodSwap).
    Stream(ctx)
```

### Multi-Chain — Scan BSC

```go
client, err := sentinel.New(
    sentinel.WithChainID(uint64(sentinel.ChainBSC)),
    sentinel.WithRPCURL("https://bsc-dataseed.binance.org"),
)
```

Supported chains out of the box:

| Constant | Chain | Chain ID |
|---|---|---|
| `sentinel.ChainETH` | Ethereum Mainnet | 1 |
| `sentinel.ChainBSC` | BNB Smart Chain | 56 |
| `sentinel.ChainPolygon` | Polygon | 137 |
| `sentinel.ChainArbitrum` | Arbitrum One | 42161 |

### Mempool Monitoring

```go
mempoolClient := client.Mempool()

pendingCh, err := mempoolClient.Subscribe(ctx,
    sentinel.FilterByMinValueETH(1.0),       // value > 1 ETH
    sentinel.FilterByMinGas(10),             // gasPrice > 10 Gwei
    sentinel.FilterByMethod("0x38ed1739"),   // swapExactTokensForTokens
)

for ptx := range pendingCh {
    fmt.Printf("pending tx: %s  from: %s  value: %s\n",
        ptx.Tx.Hash().Hex(),
        ptx.From.Hex(),
        ptx.Tx.Value().String(),
    )
}
```

Available filters:

| Filter | Description |
|---|---|
| `FilterByMinValueETH(eth)` | Minimum ETH value (ETH unit) |
| `FilterByMinValueWei(wei)` | Minimum ETH value (wei) |
| `FilterByMinGas(gwei)` | Minimum gas price (Gwei) |
| `FilterByMethod(sig...)` | Match method signature (e.g. `"0x38ed1739"`) |

### Wallet

```go
// Generate a new mnemonic
mnemonic, err := sentinel.GenerateMnemonic()

client, err := sentinel.New(
    sentinel.WithMnemonic(mnemonic),
    sentinel.WithRPCURL("https://mainnet.infura.io/v3/YOUR_KEY"),
)

walletClient, err := client.Wallet()

// Derive account addresses (BIP44: m/44'/60'/0'/0/index)
addr0, _ := walletClient.Address(0)
addr1, _ := walletClient.Address(1)

// Query ETH balance
balance, _ := walletClient.ETHBalance(ctx, addr0)

// Send a transaction
result, err := walletClient.Send(ctx, 0, toAddr, big.NewInt(1e18), nil)
fmt.Println("tx hash:", result.Hash.Hex())

// Sign a message (EIP-191)
sig, _ := walletClient.SignMessage(0, []byte("Hello ETH Sentinel"))
```

### Whale Monitoring

```go
import "github.com/ETHSentinel/internal/service/whale"

whaleSvc := whale.NewService(whale.Config{
    MinUSDValue: decimal.NewFromInt(100_000), // $100,000 threshold
}, logger)

eventCh := make(chan []sentinel.Event, 10)
whaleBehaviors := whaleSvc.Run(ctx, eventCh)

for be := range whaleBehaviors {
    fmt.Printf("[whale] %s  type=%s  value=$%s  confidence=%.2f\n",
        be.Address.Hex(), be.EventType, be.ValueUSD, be.Confidence)
}
```

---

## Adding a New DEX Protocol

Only two files need to be touched. No changes to Scanner, chain of responsibility, or gRPC handlers.

**Step 1 — Create the parser**

```
internal/parser/dex/curve/parser.go
```

```go
package curve

import (
    "github.com/ETHSentinel/internal/parser/comm"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/crypto"
)

var SigTokenExchange = comm.EventSig(
    crypto.Keccak256Hash([]byte("TokenExchange(address,int128,uint256,int128,uint256)")),
)

type Parser struct{ invoker *comm.EventParseInvoker }

func NewParser(client ...*ethclient.Client) *Parser {
    p := &Parser{invoker: comm.NewEventParseInvoker(comm.ProtocolImpl("Curve"))}
    p.invoker.RegisterOne(SigTokenExchange, comm.EventMethodSwap, p.parse)
    return p
}

func (p *Parser) HandleEvent(sig comm.EventSig, log types.Log, meta comm.EventMetadata) (comm.UnifiedEvent, error) {
    return p.invoker.HandleEvent(sig, log, meta)
}
func (p *Parser) ListEventSigs() []comm.EventSig          { return p.invoker.ListEventSigs() }
func (p *Parser) SetFilter(m []comm.EventMethod)          { p.invoker.SetFilter(m) }
func (p *Parser) parse(log types.Log, meta comm.EventMetadata) (comm.UnifiedEvent, error) {
    // decode log.Data using Curve ABI ...
    return &comm.UnifiedEventData{ /* ... */ }, nil
}
```

**Step 2 — Register in `internal/parser/dex/register.go`**

```go
import curve "github.com/ETHSentinel/internal/parser/dex/curve"

func RegisterAll(mgr *comm.ProtocolImplManager, client *ethclient.Client) error {
    mgr.RegisterStrategy(comm.ProtocolImplUniswapV2, uniswapv2.NewParser(client))
    mgr.RegisterStrategy(comm.ProtocolImpl("Curve"), curve.NewParser(client)) // ← add this
    return nil
}
```

**Step 3 — Use it**

```go
sentinel.NewScanBuilder(client).
    WithDEX(comm.ProtocolImpl("Curve"), comm.EventMethodSwap).
    Stream(ctx)
```

---

## Adding a New Chain

Add one function in `internal/chain/chains.go` and one line in `registry.go`:

```go
// chains.go
func Optimism() Chain {
    return NewBaseChain(ChainOptimism, "optimism", "ETH",
        "https://mainnet.optimism.io",
        "wss://ws-mainnet.optimism.io",
    )
}

// registry.go — init()
Register(Optimism())
```

Then use it:

```go
client, _ := sentinel.New(sentinel.WithChainID(10)) // Optimism chain ID
```

---

## What Was Fixed vs Original Code

| Issue | Original | Fixed |
|---|---|---|
| `ListEvents()` type mismatch | Returned `map`, caused for-range compile error | Returns `[]EventSig`, O(1) lookup via `sigIndex` |
| `SetEvents` was a no-op | Filter had zero effect | `SetFilter` actually filters by `EventMethod` name |
| Filterer rebuilt on every call | `NewUniswappairFilterer` called per event | Lazy-loaded + `sync.RWMutex` cache, created once per pair |
| Global panics | `log.Fatalf` / `panic` scattered throughout | Errors returned explicitly, panic only for constant ABI parsing |
| Missing dependency injection | Hardcoded paths, `RegisterAllParser()` had side effects | Constructor injection, `NewEngine(client)` explicit init |
| No gRPC layer | REST + internal direct calls | Full gRPC: unary + server stream + interceptors |
| No SDK entry point | None | `Client` interface + `NewScanBuilder` fluent API |
| Single chain only | Hardcoded Ethereum | Multi-chain pool, supports ETH/BSC/Polygon/Arbitrum |
| No Mempool support | None | `Monitor` with pluggable filters, pending tx subscription |
| No wallet support | Scattered files | HD wallet (BIP44) + keystore + tx builder unified under `wallet.Manager` |

---

## Docker

```bash
docker build -t eth-sentinel:latest .

docker run -e ETH_RPC_URL=https://mainnet.infura.io/v3/YOUR_KEY \
           -e ETH_WS_URL=wss://mainnet.infura.io/ws/v3/YOUR_KEY \
           -p 50051:50051 \
           eth-sentinel:latest
```

---

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `ETH_RPC_URL` | — | Ethereum RPC endpoint (required) |
| `ETH_WS_URL` | — | Ethereum WebSocket endpoint (required for Mempool/Subscribe) |
| `ETH_PROXY` | — | Proxy URL, e.g. `socks5://127.0.0.1:1080` |
| `GRPC_ADDR` | `:50051` | gRPC server listen address |
