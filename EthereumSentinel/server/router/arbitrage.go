package router

import (
	"github.com/Crypto-ChainSentinel/server/handler"
	"github.com/gin-gonic/gin"
)

func InitArbitrageRouter(r *gin.Engine) {
	arbitrageHandler := handler.ArbitrageHandler{}

	onchainGroup := r.Group("/onchainTianeye")

	arbitrageGroup := onchainGroup.Group("/arbitrage")
	{
		arbitrageGroup.GET("/opportunity", arbitrageHandler.QueryArbitrageOpportunity)
		arbitrageGroup.GET("/transactionResult", arbitrageHandler.QueryTransactionResult)
	}

	// // ============Defi 子模块===================
	defiGroup := onchainGroup.Group("/defi")
	// Dex 子模块
	dexGroup := defiGroup.Group("/dex")
	{
		dexGroup.GET("/newTokenPair", arbitrageHandler.QueryNewTokenPair)
		dexGroup.GET("/LiquidityChange", arbitrageHandler.QueryLiquidityChange)
	}
}
