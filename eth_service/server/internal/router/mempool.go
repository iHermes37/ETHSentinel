package router

import (
	"github.com/CryptoQuantX/chain_monitor/server/internal/handler"
	"github.com/gin-gonic/gin"
)

func InitMemPoolsRouter(r *gin.Engine) {
	memPoolHandler := handler.MempoolHandler{}

	onchainGroup := r.Group("/onchainTianeye")

	arbitrageGroup := onchainGroup.Group("/mempools")
	{
		arbitrageGroup.GET("/mempoolTx", memPoolHandler.QueryMempoolTx)
		arbitrageGroup.GET("/mempoolStats", memPoolHandler.QueryMempoolStats)
	}
}
