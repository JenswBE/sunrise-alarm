package gui

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/JenswBE/sunrise-alarm/src/services/gui/entities"
)

func (h *Handler) handleClock(c *gin.Context) {
	h.htmlWithFlashes(c, http.StatusOK, &entities.ClockTemplate{
		BaseData: entities.BaseData{
			Title:      "Clock",
			ParentPath: "clock",
		},
	})
}
