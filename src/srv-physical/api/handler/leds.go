package handler

import (
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/api/openapi"
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/api/presenter"
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/usecases/leds"
	"github.com/gin-gonic/gin"
)

const pathPrefixLeds = "/leds"

type LedsHandler struct {
	service leds.Usecase
}

func NewLedsHandler(service leds.Usecase) *LedsHandler {
	return &LedsHandler{service: service}
}

func (h *LedsHandler) RegisterRoutes(rg *gin.RouterGroup) {
	groupMock := rg.Group(pathPrefixLeds)
	groupMock.GET("/", h.getLeds)
	groupMock.PUT("/", h.setLeds)
	groupMock.DELETE("/", h.clearLeds)
}

func (h *LedsHandler) getLeds(c *gin.Context) {
	// TODO
	c.String(204, "")
}

func (h *LedsHandler) setLeds(c *gin.Context) {
	// Parse body
	var body openapi.Leds
	if err := c.BindJSON(&body); err != nil {
		c.JSON(errToResponse(err))
		return
	}

	// Convert body to entity
	color, brightness, err := presenter.LedsToEntity(body)
	if err != nil {
		c.JSON(errToResponse(err))
		return
	}

	// Call service
	h.service.SetColorAndBrightness(color, brightness)
	c.String(204, "")
}

func (h *LedsHandler) clearLeds(c *gin.Context) {
	h.service.Clear()
	c.String(204, "")
}
