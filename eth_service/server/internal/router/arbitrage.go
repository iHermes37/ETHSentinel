package router

import (
	"github.com/CryptoQuantX/chain_monitor/server/internal/handler"
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
}
