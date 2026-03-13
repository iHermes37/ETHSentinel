package router

import (
	"github.com/Crypto-ChainSentinel/server/handler"
	"github.com/gin-gonic/gin"
)

func InitWhaleRouter(r *gin.Engine) {
	WhaleHandler := handler.NewWhaleHandler()

	onchainGroup := r.Group("/onchain_sentinel")
	whales := onchainGroup.Group("/whale")

	whales.POST("/:address", WhaleHandler.CreateWhale)
	whales.GET("", WhaleHandler.ListWhales) // 获取列表
	whales.GET("/:address", WhaleHandler.GetWhale)
	whales.GET("/:address/assets", WhaleHandler.GetWhaleAssets)
	whales.GET("/:address/transactions", WhaleHandler.WhaleTransactions)

	whales.POST("/:address/track", WhaleHandler.TrackWhale)
	whales.POST("/batch/track", WhaleHandler.TrackAllWhale)
	whales.POST("/batch/captured", WhaleHandler.CapturedWhales)
}
