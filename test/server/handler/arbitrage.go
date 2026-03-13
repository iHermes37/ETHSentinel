package handler

import (
	"github.com/Crypto-ChainSentinel/server/service"
	"github.com/gin-gonic/gin"
)

type ArbitrageHandler struct {
	service *service.ArbitrageService
}

func NewArbitrageHandler() *ArbitrageHandler {
	return &ArbitrageHandler{
		service: service.NewArbitrageService(),
	}
}

func (h *ArbitrageHandler) QueryArbitrageOpportunity(c *gin.Context) {

}

func (h *ArbitrageHandler) QueryTransactionResult(c *gin.Context) {

}

func (h *ArbitrageHandler) QueryNewTokenPair(c *gin.Context) {

}

func (h *ArbitrageHandler) QueryLiquidityChange(c *gin.Context) {

}
