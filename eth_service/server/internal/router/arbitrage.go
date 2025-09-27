package router

import (
	"github.com/Crypto-ChainSentinel/server/internal/handler"
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
