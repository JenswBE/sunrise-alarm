package gui

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/JenswBE/sunrise-alarm/src/services/alarm"
	"github.com/JenswBE/sunrise-alarm/src/services/gui/entities"
	"github.com/JenswBE/sunrise-alarm/src/utils/flashstore"
	"github.com/JenswBE/sunrise-alarm/src/utils/pubsub"
)

type Handler struct {
	pubSub         pubsub.PubSub
	alarmService   alarm.Service
	isDebugEnabled bool
	// Since this is a single user system, we can implement
	// in-memory flash messages instead of relying on session cookies.
	messages *flashstore.FlashStore[entities.Message]
}

func NewHandler(pubSub pubsub.PubSub, alarmService alarm.Service, isDebugEnabled bool) *Handler {
	return &Handler{
		pubSub:         pubSub,
		alarmService:   alarmService,
		isDebugEnabled: isDebugEnabled,
		messages:       &flashstore.FlashStore[entities.Message]{},
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	// Register static routes
	staticFS, err := fs.Sub(htmlContent, "html/static")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to define sub FS for static content")
	}
	r.StaticFS("static", http.FS(staticFS))

	// Register GUI routes
	r.GET("/", h.handleClock)
	r.GET("/alarms", h.handleAlarmsList)
	r.GET("/alarms/:alarm_id", h.handleAlarmsForm)
	r.POST("/alarms/:alarm_id", h.handleAlarmsFormPOST)
	r.POST("/alarms/:alarm_id/enabled", h.handleAlarmsSetEnabled)    // HTML forms only allow GET or POST
	r.POST("/alarms/:alarm_id/skip-next", h.handleAlarmsSetSkipNext) // HTML forms only allow GET or POST
	r.POST("/alarms/:alarm_id/delete", h.handleDeleteAlarm)          // HTML forms only allow GET or POST

	// Register API routes
	api := r.Group("/api")
	api.GET("/next-ring-time", h.handleAPINextRingTime)

	// Register debug routes
	if h.isDebugEnabled {
		debug := r.Group("/debug")
		debug.GET("/", h.handleDebug)
		debug.POST("/simulate-button-pressed-short", h.handleSimulateButtonPressedShort)
		debug.POST("/simulate-button-pressed-long", h.handleSimulateButtonPressedLong)
		debug.POST("/reboot", h.handleReboot)
	}
}

func html(c *gin.Context, code int, template entities.Template) {
	c.HTML(code, template.GetTemplateName(), template)
}

func (h *Handler) html200WithFlashes(c *gin.Context, template entities.Template) {
	// Get and convert flashes
	template.SetMessages(h.messages.Get())

	// Display page
	html(c, http.StatusOK, template)
}

func redirect(c *gin.Context, redirectLocation string) {
	c.Redirect(http.StatusSeeOther, redirectLocation)
}

func (h *Handler) redirectWithMessage(c *gin.Context, redirectLocation string, messageType entities.MessageType, messageFormat string, messageArgs ...any) {
	h.messages.Add(entities.Message{
		Type:    messageType,
		Content: fmt.Sprintf(messageFormat, messageArgs...),
	})
	c.Redirect(http.StatusSeeOther, redirectLocation)
}

func (h *Handler) redirectWithErrorMessage(c *gin.Context, redirectLocation, messageFormat string, messageArgs ...any) {
	h.redirectWithMessage(c, redirectLocation, entities.MessageTypeError, messageFormat, messageArgs...)
}
