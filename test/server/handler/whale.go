package handler

import (
	model "github.com/Crypto-ChainSentinel/server/schemas"
	"github.com/Crypto-ChainSentinel/server/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
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

// ====================================================

func (h *WhaleHandler) WhaleTransactions(c *gin.Context) {
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

// ==============================================================

func (h *WhaleHandler) CapturedWhales(c *gin.Context) {
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

func (h *WhaleHandler) TrackWhale(c *gin.Context) {
	var QueryParam model.TrackWhale
	if err := c.ShouldBindQuery(&QueryParam); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
}

func (h *WhaleHandler) TrackAllWhale(c *gin.Context) {

}

// ============================================================

func (h *WhaleHandler) CreateWhale(c *gin.Context) {

}

func (h *WhaleHandler) GetWhale(c *gin.Context) {
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

func (h *WhaleHandler) ListWhales(c *gin.Context) {}

// ================================================================

func (h *WhaleHandler) GetWhaleAssets(c *gin.Context) {
	// 从 URL 查询参数获取地址
	whaleAddr := c.Param("address")
	// 调用服务层
	response, err := h.service.QueryWhaleAssets(whaleAddr)
	if err != nil {
		// 根据错误类型返回不同的状态码
		if err == gorm.ErrRecordNotFound {
			c.JSON(404, gin.H{
				"code": 404,
				"msg":  "未找到",
				"data": "该地址不存在",
			})
		} else {
			// 记录错误日志
			log.Printf("查询失败: %v", err)
			c.JSON(500, gin.H{
				"code": 500,
				"msg":  "查询失败",
				"data": err.Error(), // 生产环境建议隐藏具体错误
			})
		}
		return
	}
	// 成功响应
	c.JSON(200, gin.H{
		"code": 0,
		"msg":  "success",
		"data": response,
	})
}
