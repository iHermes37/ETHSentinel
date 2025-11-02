package router

import (
	"github.com/Crypto-ChainSentinel/server/handler"
	"github.com/gin-gonic/gin"
)

func InitContractPoolRouter(r *gin.Engine) {
	contractHandler := handler.ContractHandler{}

	// 顶层组
	onchainGroup := r.Group("/onchainTianeye")
	// onchainGroup.GET("/contracts", contractHandler.QueryContracts)
	ContractGroup := onchainGroup.Group("/contract")
	{
		ContractGroup.GET("/deployed", contractHandler.GetDeployContracts)
		ContractGroup.GET("/arbitrageBot", contractHandler.GetArbitrageBot)
		ContractGroup.GET("/alphaProjects", contractHandler.GetAlphaProjects)
		ContractGroup.GET("/newToken", contractHandler.GetNewToken)
	}
}
