package router

import (
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	// 注册不同模块的路由
	InitContractPoolRouter(r)
	InitWhaleRouter(r)
	InitArbitrageRouter(r)
	InitMemPoolsRouter(r)
	return r
}
