package handler

import (
	"github.com/CryptoQuantX/chain_monitor/server/internal/model"
	"github.com/CryptoQuantX/chain_monitor/server/internal/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type WhaleHandler struct {
	service *service.WhaleService
}

func NewWhaleHandler() *WhaleHandler {
	return &WhaleHandler{
		service: service.NewWhaleService(),
	}
}

func (h *WhaleHandler) QueryCapturedWhales(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	var QueryParam model.WhaleQueryParams
	if err := c.ShouldBindQuery(&QueryParam); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	whales, err := h.service.QueryCapturedWhales(QueryParam, page, pageSize)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, whales) // 返回 JSON 响应
}

func (h *WhaleHandler) QueryWhaleTransaction(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	var QueryParam model.WhaleTradeParams
	if err := c.ShouldBindQuery(&QueryParam); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	whales, err := h.service.QueryWhaleTrade(QueryParam, page, pageSize)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, whales) // 返回 JSON 响应

}

func (h *WhaleHandler) QueryWhaleHoldings(c *gin.Context) {

}
