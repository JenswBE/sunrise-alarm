package handler

import (
	"github.com/gin-gonic/gin"
)

const pathPrefixMock = "/mock"

type MockHandler struct{}

func NewMockHandler() *MockHandler {
	return &MockHandler{}
}

func (h *MockHandler) RegisterRoutes(rg *gin.RouterGroup) {
	groupMock := rg.Group(pathPrefixMock)
	groupMock.POST("/button/pressed", h.buttonPressed)
	groupMock.POST("/button/long_pressed", h.buttonLongPressed)
}

func (h *MockHandler) buttonPressed(c *gin.Context) {
	c.String(204, "")
}

func (h *MockHandler) buttonLongPressed(c *gin.Context) {
	c.String(204, "")
}
