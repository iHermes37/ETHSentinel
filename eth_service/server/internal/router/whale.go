package router

import (
	"github.com/Crypto-ChainSentinel/server/internal/handler"
	"github.com/gin-gonic/gin"
)

func InitWhaleRouter(r *gin.Engine) {
	WhaleHandler := handler.WhaleHandler{}

	whalesGroup := r.Group("/onchainTianeye")
	{
		// userGroup.GET("/:id", userHandler.GetUser)
		whalesGroup.GET("/capturedWhale", WhaleHandler.QueryCapturedWhales)
		whalesGroup.GET("/whaleTransaction", WhaleHandler.QueryWhaleTransaction)
		whalesGroup.GET("/whaleHoldings", WhaleHandler.QueryWhaleHoldings)
	}
}
