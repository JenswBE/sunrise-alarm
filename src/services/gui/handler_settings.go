package gui

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/JenswBE/sunrise-alarm/src/services/gui/entities"
)

func (h *Handler) handleSettings(c *gin.Context) {
	h.htmlWithFlashes(c, http.StatusOK, &entities.SettingsTemplate{
		BaseData: entities.BaseData{
			Title:      "Settings",
			ParentPath: "settings",
		},
	})
}
