package gui

import (
	"os/exec"

	"github.com/gin-gonic/gin"

	"github.com/JenswBE/sunrise-alarm/src/services/gui/entities"
	"github.com/JenswBE/sunrise-alarm/src/utils/pubsub"
)

func (h *Handler) handleDebug(c *gin.Context) {
	h.html200WithFlashes(c, &entities.DebugTemplate{
		BaseData: entities.BaseData{
			Title:      "Debug",
			ParentPath: "debug",
		},
	})
}

func (h *Handler) handleSimulateButtonPressedShort(c *gin.Context) {
	h.pubSub.Publish(pubsub.EventButtonPressedShort{})
	redirect(c, "/debug")
}

func (h *Handler) handleSimulateButtonPressedLong(c *gin.Context) {
	h.pubSub.Publish(pubsub.EventButtonPressedLong{})
	redirect(c, "/debug")
}

func (h *Handler) handleReboot(c *gin.Context) {
	if err := exec.Command("reboot").Run(); err != nil {
		h.redirectWithErrorMessage(c, "/debug", "Failed to reboot device: %v", err)
		return
	}
	redirect(c, "/debug")
}
