package main

import (
	"github.com/Crypto-ChainSentinel/server/router"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// 创建默认路由引擎
	r := gin.Default()

	// 定义 GET 路由
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// 定义 POST 路由
	r.POST("/hello", func(c *gin.Context) {
		name := c.PostForm("name")
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello " + name,
		})
	})

	router.InitRouter()

	// 启动服务，默认监听 0.0.0.0:8080
	r.Run()
}
