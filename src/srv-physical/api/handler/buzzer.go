package handler

import (
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/usecases/buzzer"
	"github.com/gin-gonic/gin"
)

const pathPrefixBuzzer = "/buzzer"

type BuzzerHandler struct {
	service buzzer.Usecase
}

func NewBuzzerHandler(service buzzer.Usecase) *BuzzerHandler {
	return &BuzzerHandler{service: service}
}

func (h *BuzzerHandler) RegisterRoutes(rg *gin.RouterGroup) {
	groupMock := rg.Group(pathPrefixBuzzer)
	groupMock.PUT("", h.startBuzzer)
	groupMock.DELETE("", h.stopBuzzer)
}

func (h *BuzzerHandler) startBuzzer(c *gin.Context) {
	h.service.Start()
}

func (h *BuzzerHandler) stopBuzzer(c *gin.Context) {
	h.service.Stop()
}
