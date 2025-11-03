package router

import (
	"github.com/Crypto-ChainSentinel/server/handler"
	"github.com/gin-gonic/gin"
)

func InitWhaleRouter(r *gin.Engine) {
	WhaleHandler := handler.WhaleHandler{}

	onchainGroup := r.Group("/onchainTianeye")

	whales := onchainGroup.Group("/whales")
	{
		whales.POST("/AddWhale", WhaleHandler.AddWhale)
		whales.GET("/MonitoredWhale", WhaleHandler.GetMonitoredWhale)
	}

	task := whales.Group("/task")
	{
		task.POST("/capturedWhale", WhaleHandler.CapturedWhales)
		task.POST("/trackWhale", WhaleHandler.TrackWhales)
		// whalesGroup.GET("/whaleTransaction", WhaleHandler.QueryWhaleTransaction)
		// whalesGroup.GET("/whaleHoldings", WhaleHandler.QueryWhaleHoldings)
	}

}
