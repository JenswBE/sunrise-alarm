package handler

import (
	"github.com/JenswBE/sunrise-alarm/src/srv-config/usecases/alarms"
	"github.com/gin-gonic/gin"
)

const pathPrefixAlarms = "/alarms"

type AlarmsHandler struct {
	alarms alarms.Usecase
}

func NewAlarmsHandler(alarmsService alarms.Usecase) *AlarmsHandler {
	return &AlarmsHandler{
		alarms: alarmsService,
	}
}

func (h *AlarmsHandler) RegisterRoutes(rg *gin.RouterGroup) {
	groupAlarms := rg.Group(pathPrefixAlarms)
	groupAlarms.GET("", h.lockBacklight)
	groupAlarms.POST("", h.lockBacklight)
	groupAlarms.PUT("/lock", h.lockBacklight)
	groupAlarms.DELETE("/lock", h.unlockBacklight)
}

func (h *AlarmsHandler) lockBacklight(c *gin.Context) {
	// Stub - We're not controlling the backlight yet
	c.String(204, "")
}

func (h *AlarmsHandler) unlockBacklight(c *gin.Context) {
	// Stub - We're not controlling the backlight yet
	c.String(204, "")
}
