package router

import (
    "github.com/gin-gonic/gin"
    "your_project/internal/handler"
)




func InitRouter() *gin.Engine {
    r := gin.Default()

    // 注册不同模块的路由
    InitContractPoolsRouter(r)
    // InitContractRouter(r)

    return r
}

