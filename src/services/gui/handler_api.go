package gui

import (
	"net/http"

	"github.com/JenswBE/sunrise-alarm/src/services/gui/entities"
	"github.com/gin-gonic/gin"
)

func (h *Handler) handleAPINextRingTime(c *gin.Context) {
	nextRingTime := h.alarmService.GetNextRingTime()
	if nextRingTime.IsZero() {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, entities.NextRingTimeResponseFromEntity(nextRingTime))
}
