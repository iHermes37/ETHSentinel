package handler

import (
	"github.com/Crypto-ChainSentinel/server/internal/model"
	"github.com/Crypto-ChainSentinel/server/internal/service"
	"github.com/gin-gonic/gin"
	"strconv"
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

func (h *ContractHandler) QueryNewTokenPair(c *gin.Context) {

}

func (h *ContractHandler) QueryLiquidityChange(c *gin.Context) {

}
