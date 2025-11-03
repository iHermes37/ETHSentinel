package router

import (
	"github.com/Crypto-ChainSentinel/server/handler"
	"github.com/gin-gonic/gin"
)

func InitContractPoolRouter(r *gin.Engine) {
	contractHandler := handler.ContractHandler{}

	// 顶层组
	onchainGroup := r.Group("/onchainTianeye")
	contract := onchainGroup.Group("/contract")

	// ========监控合约======================
	monitor := contract.Group("/monitor")
	{
		monitor.GET("/deployed", contractHandler.GetMonitorContracts)
		monitor.GET("/arbitrageBot", contractHandler.GetArbitrageBot)
		monitor.GET("/alphaProjects", contractHandler.GetAlphaProjects)
		monitor.GET("/newToken", contractHandler.GetNewToken)
	}

	// =======已知合约表======================
	defi := contract.Group("/defi")
	{
		defi.PUT("/newcontract", contractHandler.AddKnownContracts)
	}
}
