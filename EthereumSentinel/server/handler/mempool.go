package handler

import (
	"github.com/Crypto-ChainSentinel/server/internal/service"
	"github.com/gin-gonic/gin"
)

type MempoolHandler struct {
	service *service.ContractService
}

func NewMempoolHandler() *MempoolHandler {
	return &MempoolHandler{
		service: service.NewMempoolService(),
	}
}

func (h *MempoolHandler) QueryMempoolTx(c *gin.Context) {

}

func (h *MempoolHandler) QueryMempoolStats(c *gin.Context) {

}
