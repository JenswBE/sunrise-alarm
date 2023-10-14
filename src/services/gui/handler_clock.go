package gui

import (
	"github.com/gin-gonic/gin"

	"github.com/JenswBE/sunrise-alarm/src/services/gui/entities"
)

func (h *Handler) handleClock(c *gin.Context) {
	h.html200WithFlashes(c, &entities.ClockTemplate{
		BaseData: entities.BaseData{
			Title:      "Clock",
			ParentPath: "clock",
		},
	})
}
