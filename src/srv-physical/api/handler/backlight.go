package handler

import (
	"github.com/gin-gonic/gin"
)

const pathPrefixBacklight = "/backlight"

type BacklightHandler struct{}

func NewBacklightHandler() *BacklightHandler {
	return &BacklightHandler{}
}

func (h *BacklightHandler) RegisterRoutes(rg *gin.RouterGroup) {
	groupMock := rg.Group(pathPrefixBacklight)
	groupMock.PUT("/lock", h.lockBacklight)
	groupMock.DELETE("/lock", h.unlockBacklight)
}

func (h *BacklightHandler) lockBacklight(c *gin.Context) {
	// Stub - We're not controlling the backlight yet
	c.String(204, "")
}

func (h *BacklightHandler) unlockBacklight(c *gin.Context) {
	// Stub - We're not controlling the backlight yet
	c.String(204, "")
}
