package router

import (
	"github.com/Crypto-ChainSentinel/server/handler"
	"github.com/gin-gonic/gin"
)

func InitContractRouter(r *gin.Engine) {
	contractHandler := handler.NewContractHandler()

	// 顶层组
	onchainGroup := r.Group("/onchain_sentinel")
	contractGroup := onchainGroup.Group("/contracts") // 使用复数形式，符合RESTful惯例

	// ========监控合约======================
	contractGroup.GET("", contractHandler.ListContracts)              // 获取所有监控合约列表
	contractGroup.POST("", contractHandler.CreateContract)            // 添加新监控合约
	contractGroup.GET("/:address", contractHandler.GetContract)       // 获取特定合约详情
	contractGroup.PUT("/:address", contractHandler.UpdateContract)    // 更新合约信息
	contractGroup.DELETE("/:address", contractHandler.DeleteContract) // 删除监控合约

	// 5. 合约事件/交易监控
	contractGroup.GET("/:address/events", contractHandler.GetContractEvents)             // 获取合约事件
	contractGroup.GET("/:address/transactions", contractHandler.GetContractTransactions) // 获取合约交易
	contractGroup.GET("/:address/calls", contractHandler.GetContractCalls)               // 获取合约调用

	tokenGroup := contractGroup.Group("/tokens")
	{
		tokenGroup.GET("", contractHandler.ListNewTokens)              // 获取所有新代币
		tokenGroup.POST("", contractHandler.AddNewToken)               // 添加新代币
		tokenGroup.GET("/:address", contractHandler.GetToken)          // 获取代币详情
		tokenGroup.GET("/recent", contractHandler.GetRecentTokens)     // 获取最近添加的代币
		tokenGroup.GET("/trending", contractHandler.GetTrendingTokens) // 获取热门代币
	}

	//monitor := contract.Group("/monitor")
	//{
	//	monitor.GET("/deployed", contractHandler.GetMonitorContracts)
	//	monitor.PUT("/new_defi", contractHandler.AddDefiContracts)
	//	//monitor.GET("/arbitrageBot", contractHandler.GetArbitrageBot)
	//	//monitor.GET("/alphaProjects", contractHandler.GetAlphaProjects)
	//	//monitor.GET("/newToken", contractHandler.GetNewToken)
	//}

}

//func InitContractRouter(r *gin.Engine) {
//	contractHandler := handler.NewContractHandler()
//
//	// 顶层组
//	onchainGroup := r.Group("/onchain_sentinel")
//	contractGroup := onchainGroup.Group("/contracts") // 使用复数形式，符合RESTful惯例
/
//
//	// 4. 合约特定操作
//	contractGroup.POST("/:address/monitor", contractHandler.StartMonitoring)   // 开始监控合约
//	contractGroup.POST("/:address/pause", contractHandler.PauseMonitoring)     // 暂停监控
//	contractGroup.POST("/:address/resume", contractHandler.ResumeMonitoring)   // 恢复监控
//	contractGroup.GET("/:address/status", contractHandler.GetMonitoringStatus) // 获取监控状态
//
//	// 5. 合约事件/交易监控
//	contractGroup.GET("/:address/events", contractHandler.GetContractEvents)             // 获取合约事件
//	contractGroup.GET("/:address/transactions", contractHandler.GetContractTransactions) // 获取合约交易
//	contractGroup.GET("/:address/calls", contractHandler.GetContractCalls)               // 获取合约调用
/
//
//	// 7. 套利机器人监控
//	arbitrageGroup := contractGroup.Group("/arbitrage")
//	{
//		arbitrageGroup.GET("", contractHandler.ListArbitrageBots)             // 获取所有套利机器人
//		arbitrageGroup.POST("", contractHandler.AddArbitrageBot)              // 添加套利机器人
//		arbitrageGroup.GET("/:address", contractHandler.GetArbitrageBot)      // 获取套利机器人详情
//		arbitrageGroup.GET("/active", contractHandler.GetActiveArbitrageBots) // 获取活跃的套利机器人
//	}
//
//	// 8. Alpha 项目监控
//	alphaGroup := contractGroup.Group("/alpha")
//	{
//		alphaGroup.GET("", contractHandler.ListAlphaProjects)         // 获取所有Alpha项目
//		alphaGroup.POST("", contractHandler.AddAlphaProject)          // 添加Alpha项目
//		alphaGroup.GET("/:address", contractHandler.GetAlphaProject)  // 获取Alpha项目详情
//		alphaGroup.GET("/trending", contractHandler.GetTrendingAlpha) // 获取热门Alpha项目
//	}

//
//	// ======== 合约部署监控 ======================
//
//	// 10. 部署监控（对应原 deployed）
//	deployments := contractGroup.Group("/deployments")
//	{
//		deployments.GET("", contractHandler.GetRecentDeployments)              // 获取最近部署的合约
//		deployments.GET("/monitored", contractHandler.GetMonitoredDeployments) // 获取已监控的部署
//		deployments.POST("/watch", contractHandler.WatchDeployment)            // 开始监控新部署
//		deployments.GET("/stats", contractHandler.GetDeploymentStats)          // 获取部署统计
//		deployments.GET("/:tx_hash", contractHandler.GetDeploymentByTx)        // 根据交易哈希获取部署详情
//	}

//}
