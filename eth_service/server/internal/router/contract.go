package router

import (
	"github.com/CryptoQuantX/chain_monitor/server/internal/handler"
	"github.com/gin-gonic/gin"
)

func InitContractPoolsRouter(r *gin.Engine) {
	contractHandler := handler.ContractHandler{}

	// 顶层组
	onchainGroup := r.Group("/onchainTianeye")
	onchainGroup.GET("/contractPools", contractHandler.QueryContracts)

	// Defi 子模块
	defiGroup := onchainGroup.Group("/defi")
	// Dex 子模块
	dexGroup := defiGroup.Group("/dex")
	{
		dexGroup.GET("/newTokenPair", contractHandler.QueryNewTokenPair)
		dexGroup.GET("/LiquidityChange", contractHandler.QueryLiquidityChange)
	}
}
