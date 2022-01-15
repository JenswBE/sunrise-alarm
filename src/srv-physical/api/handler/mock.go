package handler

import (
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/usecases/mqtt"
	"github.com/gin-gonic/gin"
)

const pathPrefixMock = "/mock"

type MockHandler struct {
	service mqtt.Usecase
}

func NewMockHandler(service mqtt.Usecase) *MockHandler {
	return &MockHandler{service: service}
}

func (h *MockHandler) RegisterRoutes(rg *gin.RouterGroup) {
	groupMock := rg.Group(pathPrefixMock)
	groupMock.POST("/button/pressed", h.buttonPressed)
	groupMock.POST("/button/long_pressed", h.buttonLongPressed)
}

func (h *MockHandler) buttonPressed(c *gin.Context) {
	err := h.service.PublishButtonPressed(c.Request.Context())
	if err != nil {
		c.JSON(errToResponse(err))
		return
	}
	c.String(204, "")
}

func (h *MockHandler) buttonLongPressed(c *gin.Context) {
	err := h.service.PublishButtonLongPressed(c.Request.Context())
	if err != nil {
		c.JSON(errToResponse(err))
		return
	}
	c.String(204, "")
}
