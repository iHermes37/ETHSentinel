package handler

import (
	"github.com/CryptoQuantX/chain_monitor/server/internal/service"
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
