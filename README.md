
```
chainTxTracker
├─ .idea
│  ├─ chain_monitor.iml
│  ├─ modules.xml
│  ├─ vcs.xml
│  └─ workspace.xml
├─ contracts
│  ├─ common
│  ├─ ethereum-app
│  │  ├─ contracts
│  │  │  ├─ crocssdexArbitrager.sol
│  │  │  ├─ flashLib.sol
│  │  │  └─ Hello.sol
│  │  ├─ foundry.toml
│  │  ├─ hardhat.config.ts
│  │  ├─ ignition
│  │  │  └─ modules
│  │  │     └─ Lock.ts
│  │  ├─ lib
│  │  │  └─ forge-std
│  │  │     ├─ CONTRIBUTING.md
│  │  │     ├─ foundry.toml
│  │  │     ├─ LICENSE-APACHE
│  │  │     ├─ LICENSE-MIT
│  │  │     ├─ package.json
│  │  │     ├─ README.md
│  │  │     ├─ scripts
│  │  │     │  └─ vm.py
│  │  │     ├─ src
│  │  │     │  ├─ Base.sol
│  │  │     │  ├─ console.sol
│  │  │     │  ├─ console2.sol
│  │  │     │  ├─ interfaces
│  │  │     │  │  ├─ IERC1155.sol
│  │  │     │  │  ├─ IERC165.sol
│  │  │     │  │  ├─ IERC20.sol
│  │  │     │  │  ├─ IERC4626.sol
│  │  │     │  │  ├─ IERC6909.sol
│  │  │     │  │  ├─ IERC721.sol
│  │  │     │  │  ├─ IERC7540.sol
│  │  │     │  │  ├─ IERC7575.sol
│  │  │     │  │  └─ IMulticall3.sol
│  │  │     │  ├─ safeconsole.sol
│  │  │     │  ├─ Script.sol
│  │  │     │  ├─ StdAssertions.sol
│  │  │     │  ├─ StdChains.sol
│  │  │     │  ├─ StdCheats.sol
│  │  │     │  ├─ StdConstants.sol
│  │  │     │  ├─ StdError.sol
│  │  │     │  ├─ StdInvariant.sol
│  │  │     │  ├─ StdJson.sol
│  │  │     │  ├─ StdMath.sol
│  │  │     │  ├─ StdStorage.sol
│  │  │     │  ├─ StdStyle.sol
│  │  │     │  ├─ StdToml.sol
│  │  │     │  ├─ StdUtils.sol
│  │  │     │  ├─ Test.sol
│  │  │     │  └─ Vm.sol
│  │  │     └─ test
│  │  │        ├─ CommonBase.t.sol
│  │  │        ├─ compilation
│  │  │        │  ├─ CompilationScript.sol
│  │  │        │  ├─ CompilationScriptBase.sol
│  │  │        │  ├─ CompilationTest.sol
│  │  │        │  └─ CompilationTestBase.sol
│  │  │        ├─ fixtures
│  │  │        │  ├─ broadcast.log.json
│  │  │        │  ├─ test.json
│  │  │        │  └─ test.toml
│  │  │        ├─ StdAssertions.t.sol
│  │  │        ├─ StdChains.t.sol
│  │  │        ├─ StdCheats.t.sol
│  │  │        ├─ StdConstants.t.sol
│  │  │        ├─ StdError.t.sol
│  │  │        ├─ StdJson.t.sol
│  │  │        ├─ StdMath.t.sol
│  │  │        ├─ StdStorage.t.sol
│  │  │        ├─ StdStyle.t.sol
│  │  │        ├─ StdToml.t.sol
│  │  │        ├─ StdUtils.t.sol
│  │  │        └─ Vm.t.sol
│  │  ├─ package-lock.json
│  │  ├─ package.json
│  │  ├─ README.md
│  │  ├─ scripts
│  │  │  ├─ deploy.ts
│  │  │  └─ testFlashLoan.ts
│  │  ├─ test
│  │  │  ├─ Lock.t.sol
│  │  │  └─ Lock.ts
│  │  └─ tsconfig.json
│  ├─ solana
│  └─ solana-app
│     ├─ .prettierignore
│     ├─ Anchor.toml
│     ├─ Cargo.lock
│     ├─ Cargo.toml
│     ├─ migrations
│     │  └─ deploy.ts
│     ├─ package-lock.json
│     ├─ package.json
│     ├─ programs
│     │  └─ solana-app
│     │     ├─ Cargo.toml
│     │     ├─ src
│     │     │  ├─ dexswaps
│     │     │  ├─ lib.rs
│     │     │  └─ state
│     │     └─ Xargo.toml
│     ├─ tests
│     │  └─ solana-app.ts
│     └─ tsconfig.json
├─ eth_service
│  ├─  utils
│  ├─ .idea
│  │  ├─ eth-backend.iml
│  │  ├─ modules.xml
│  │  ├─ vcs.xml
│  │  └─ workspace.xml
│  ├─ config
│  │  ├─ abi
│  │  │  ├─ UniswapV2_Pair.json
│  │  │  └─ UniswapV2_Router.json
│  │  ├─ CexAddr.json
│  │  ├─ config.yaml
│  │  ├─ DexAddr.json
│  │  └─ whitelist.json
│  ├─ db
│  │  ├─ dbmanager.go
│  │  ├─ mysql.go
│  │  ├─ neo4j.go
│  │  └─ redis.go
│  ├─ executor
│  ├─ go.mod
│  ├─ go.sum
│  ├─ initialize
│  │  ├─ Infuraconn.go
│  │  ├─ mysqlconn.go
│  │  ├─ neo4jconn.go
│  │  └─ redisconn.go
│  ├─ main.go
│  ├─ models
│  │  ├─ arbitrage.go
│  │  ├─ block.go
│  │  ├─ Contract.go
│  │  ├─ Defi.go
│  │  ├─ MonitorPool.go
│  │  ├─ token.go
│  │  ├─ whale.go
│  │  ├─ whaleNode.go
│  │  └─ whaleTrade.go
│  ├─ modules
│  │  ├─ arbitrageur
│  │  │  ├─ arbitrageMain.go
│  │  │  ├─ arbitrageSchedule.go
│  │  │  └─ strategys
│  │  │     ├─ crossdex.go
│  │  │     ├─ strategy.go
│  │  │     └─ triangular.go
│  │  ├─ handler
│  │  │  ├─ addr.go
│  │  │  ├─ contract.go
│  │  │  ├─ defi.go
│  │  │  └─ erc.go
│  │  ├─ monitor
│  │  │  ├─ dexLiquidity.go
│  │  │  ├─ monitorMain.go
│  │  │  └─ tokenHold.go
│  │  └─ parser
│  │     ├─ DefiParser.go
│  │     ├─ Dexparser.go
│  │     ├─ ERCParser.go
│  │     └─ tx.go
│  ├─ monitor
│  │  └─ Arbitrage
│  ├─ server
│  │  ├─ internal
│  │  │  ├─ handler
│  │  │  │  ├─ arbitrage.go
│  │  │  │  ├─ contract.go
│  │  │  │  ├─ mempool.go
│  │  │  │  └─ whale.go
│  │  │  ├─ model
│  │  │  │  ├─ arbitrage.go
│  │  │  │  ├─ contract.go
│  │  │  │  ├─ mempool.go
│  │  │  │  ├─ token.go
│  │  │  │  └─ whale.go
│  │  │  ├─ repository
│  │  │  │  ├─ arbitrage.go
│  │  │  │  ├─ contract.go
│  │  │  │  └─ whale.go
│  │  │  ├─ router
│  │  │  │  ├─ arbitrage.go
│  │  │  │  ├─ contract.go
│  │  │  │  ├─ dexLiquidity.go
│  │  │  │  ├─ mempool.go
│  │  │  │  ├─ router.go
│  │  │  │  └─ whale.go
│  │  │  └─ service
│  │  │     ├─ arbitrage.go
│  │  │     ├─ contract.go
│  │  │     ├─ mempool.go
│  │  │     └─ whale.go
│  │  └─ main.go
│  ├─ storage
│  ├─ test
│  │  ├─ 1.go
│  │  ├─ 2.go
│  │  ├─ dune.go
│  │  ├─ evm_listening.go
│  │  ├─ getpair.go
│  │  ├─ graphql_listeing.go
│  │  ├─ neo4j.go
│  │  ├─ re.go
│  │  ├─ solana_listening.go
│  │  └─ uniswap.go
│  └─ utils
│     ├─ log.go
│     ├─ readconfig.go
│     └─ unitmeasure.go
├─ scripts
└─ sol_service
   ├─ Cargo.lock
   ├─ Cargo.toml
   ├─ config
   │  └─ config.yaml
   ├─ src
   │  ├─ config
   │  │  └─ publicUrl.json
   │  ├─ db
   │  │  ├─ dbmanager.rs
   │  │  └─ schema.rs
   │  ├─ initialize
   │  │  ├─ connect.rs
   │  │  └─ mod.rs
   │  ├─ main.rs
   │  ├─ models
   │  │  ├─ http
   │  │  ├─ mod.rs
   │  │  ├─ notification
   │  │  │  ├─ account.rs
   │  │  │  ├─ chainlog.rs
   │  │  │  ├─ eventTypes.rs
   │  │  │  ├─ mod.rs
   │  │  │  └─ program.rs
   │  │  ├─ spltoken
   │  │  ├─ websocket
   │  │  │  └─ WebsocketEventType.rs
   │  │  └─ whales.rs
   │  ├─ modules
   │  │  ├─ analysis
   │  │  │  ├─ dex
   │  │  │  ├─ DexToken.rs
   │  │  │  ├─ mod.rs
   │  │  │  ├─ monitWhale.rs
   │  │  │  └─ NewProject.rs
   │  │  ├─ arbitrage
   │  │  │  ├─ arbitrage.rs
   │  │  │  ├─ mod.rs
   │  │  │  └─ strategys
   │  │  │     ├─ crossdex_strategy
   │  │  │     ├─ graph_strategy
   │  │  │     └─ mev_strategy
   │  │  ├─ mod.rs
   │  │  ├─ prediction
   │  │  └─ riskmanager
   │  ├─ server
   │  │  ├─ handlers
   │  │  │  └─ mod.rs
   │  │  ├─ main.rs
   │  │  ├─ models
   │  │  │  ├─ contract.rs
   │  │  │  ├─ dex.rs
   │  │  │  ├─ mempool.rs
   │  │  │  ├─ mod.rs
   │  │  │  ├─ whale.rs
   │  │  │  └─ whaletrade.rs
   │  │  ├─ repositories
   │  │  │  └─ mod.rs
   │  │  ├─ routes
   │  │  │  └─ mod.rs
   │  │  └─ services
   │  │     └─ mod.rs
   │  ├─ subscription
   │  │  ├─ accountSubscriper.rs
   │  │  ├─ EventType.rs
   │  │  ├─ logSubscriper.rs
   │  │  ├─ mod.rs
   │  │  ├─ programSubscriper.rs
   │  │  └─ websocket.rs
   │  └─ utils
   │     ├─ heartbeat.rs
   │     ├─ mod.rs
   │     └─ WsStreamTrans.rs
   ├─ target
   │  └─ debug
   │     ├─ build
   │     └─ examples
   └─ tests
      └─ raydium.rs

```