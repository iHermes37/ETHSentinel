package handler

import (
	"strconv"
	"strings"

	"github.com/Crypto-ChainSentinel/server/model"
	"github.com/Crypto-ChainSentinel/server/service"
	"github.com/gin-gonic/gin"
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

// ============================

func (h *WhaleHandler) CapturedWhales(c *gin.Context) {
	var QueryParam model.CapturedWhale
	if err := c.ShouldBindQuery(&QueryParam); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.service.CapturedWhales(QueryParam)
	// ✅ 根据错误类型返回合适的HTTP状态码
	if strings.Contains(err.Error(), "不支持的检测方法") {
		c.JSON(400, gin.H{"error": err.Error()}) // 客户端错误
	} else {
		c.JSON(500, gin.H{"error": "内部服务器错误"}) // 服务器错误
	}
	c.JSON(200, whales) // 返回 JSON 响应
}

func (h *WhaleHandler) TrackWhales(c *gin.Context) {
	var QueryParam model.TrackWhale
	if err := c.ShouldBindQuery(&QueryParam); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.service.TrackWhales(QueryParam)
	//根据错误类型返回合适的HTTP状态码
	if strings.Contains(err.Error(), "不支持的检测方法") {
		c.JSON(400, gin.H{"error": err.Error()}) // 客户端错误
	} else {
		c.JSON(500, gin.H{"error": "内部服务器错误"}) // 服务器错误
	}
	c.JSON(200, whales) // 返回 JSON 响应
}

// =================================

func (h *WhaleHandler) AddWhale(c *gin.Context) {

}

func (h *WhaleHandler) SearchWhale(c *gin.Context) {
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

func (h *WhaleHandler) GetMonitoredWhale(c *gin.Context) {

}

// =================================
