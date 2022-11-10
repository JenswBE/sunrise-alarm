package gui

import (
	"net/http"

	"github.com/JenswBE/sunrise-alarm/src/services/gui/entities"
	"github.com/gin-gonic/gin"
)

func (h *Handler) handleAPINextAlarmToRing(c *gin.Context) {
	nextAlarm := h.alarmService.GetNextAlarmToRing()
	if nextAlarm == nil {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, entities.NextAlarmToRingResponseFromEntity(*nextAlarm))
}
