package handler

import (
	"net/http"

	"github.com/JenswBE/sunrise-alarm/src/srv-config/api/presenter"
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
	groupAlarms.GET("", h.listAlarms)
	groupAlarms.GET("/:id", h.getAlarm)
	groupAlarms.POST("", h.createAlarm)
	groupAlarms.PUT("/:id", h.updateAlarm)
	groupAlarms.DELETE("/:id", h.deleteAlarm)
}

func (h *AlarmsHandler) listAlarms(c *gin.Context) {
	// Stub - We're not controlling the backlight yet
	c.String(204, "")
}

func (h *AlarmsHandler) getAlarm(c *gin.Context) {
	// Parse ID
	id, ok := ParseIDParam(c, "id")
	if !ok {
		return
	}

	// Call service
	alarm, err := h.alarms.GetAlarm(id)
	if err != nil {
		c.JSON(ErrToResponse(err))
		return
	}
	c.JSON(http.StatusOK, presenter.AlarmFromEntity(alarm))
}

func (h *AlarmsHandler) createAlarm(c *gin.Context) {
	// Stub - We're not controlling the backlight yet
	c.String(204, "")
}

func (h *AlarmsHandler) updateAlarm(c *gin.Context) {
	// Stub - We're not controlling the backlight yet
	c.String(204, "")
}

func (h *AlarmsHandler) deleteAlarm(c *gin.Context) {
	// Parse ID
	id, ok := ParseIDParam(c, "id")
	if !ok {
		return
	}

	// Call service
	err := h.alarms.DeleteAlarm(id)
	if err != nil {
		c.JSON(ErrToResponse(err))
		return
	}
	c.String(http.StatusNoContent, "")
}
