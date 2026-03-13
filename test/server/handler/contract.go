package handler

import (
	"strconv"

	model "github.com/Crypto-ChainSentinel/server/schemas"
	"github.com/Crypto-ChainSentinel/server/service"
	"github.com/gin-gonic/gin"
)

type ContractHandler struct {
	service *service.ContractService
}

func NewContractHandler() *ContractHandler {
	return &ContractHandler{
		service: service.NewContractService(),
	}
}

func (h *ContractHandler) QueryContracts(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	var QueryParam model.ContractQueryParams

	if err := c.ShouldBindQuery(&QueryParam); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	contracts, err := h.service.QueryContracts(QueryParam, page, pageSize)
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, contracts) // 返回 JSON 响应
}

// =======================================================================

func (h *ContractHandler) GetMonitorContracts(c *gin.Context) {

}

func (h *ContractHandler) GetArbitrageBot(c *gin.Context) {

}

func (h *ContractHandler) GetAlphaProjects(c *gin.Context) {

}

func (h *ContractHandler) GetNewToken(c *gin.Context) {

}

// =============================================================================

func (h *ContractHandler) AddDefiContracts(c *gin.Context) {

}
