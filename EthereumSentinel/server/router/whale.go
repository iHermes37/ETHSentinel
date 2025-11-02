package router

import (
	"github.com/Crypto-ChainSentinel/server/handler"
	"github.com/gin-gonic/gin"
)

func InitWhaleRouter(r *gin.Engine) {
	WhaleHandler := handler.WhaleHandler{}

	whalesGroup := r.Group("/onchainTianeye")
	{
		// userGroup.GET("/:id", userHandler.GetUser)
		whalesGroup.POST("/capturedWhale", WhaleHandler.CapturedWhales)
		whalesGroup.POST("/trackWhale", WhaleHandler.TrackWhales)
		// whalesGroup.GET("/whaleTransaction", WhaleHandler.QueryWhaleTransaction)
		// whalesGroup.GET("/whaleHoldings", WhaleHandler.QueryWhaleHoldings)
	}
}
