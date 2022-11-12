package gui

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"

	globalEntities "github.com/JenswBE/sunrise-alarm/src/entities"
	"github.com/JenswBE/sunrise-alarm/src/services/gui/entities"
)

var weekdays = []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday}
var weekdaysPresets = map[string][]time.Weekday{
	"NONE":    {},
	"WEEK":    weekdays[0:5],
	"WEEKEND": weekdays[5:7],
	"ALL":     weekdays,
}

func (h *Handler) handleAlarmsList(c *gin.Context) {
	// List alarms
	alarms, err := h.alarmService.ListAlarms()
	if err != nil {
		c.String(http.StatusInternalServerError, `Failed to list alarms: %v`, err.Error())
		return
	}

	h.htmlWithFlashes(c, http.StatusOK, &entities.AlarmsListTemplate{
		BaseData: entities.BaseData{
			Title:      "Alarms",
			ParentPath: "alarms",
		},
		Alarms: alarms,
	})
}

func (h *Handler) handleAlarmsForm(c *gin.Context) {
	// Setup
	isNew := true
	var alarm globalEntities.Alarm

	rawAlarmID := c.Param("alarm_id")
	if rawAlarmID != "new" {
		// Parse alarm ID
		alarmID, err := uuid.Parse(rawAlarmID)
		if err != nil {
			h.redirectWithErrorMessage(c, "/alarms", `Failed to parse invalid alarm ID "%s": %v`, rawAlarmID, err.Error())
			return
		}

		// Get alarm
		alarm, err = h.alarmService.GetAlarm(alarmID)
		if err != nil {
			h.redirectWithErrorMessage(c, "/alarms", `Failed to fetch alarm "%s": %v`, rawAlarmID, err.Error())
			return
		}
		isNew = false
	}

	h.htmlWithFlashes(c, http.StatusOK, &entities.AlarmsFormTemplate{
		BaseData: entities.BaseData{
			Title:      "New/Edit alarm",
			ParentPath: "alarms",
		},
		AlarmBody:       entities.AlarmBodyFromEntity(alarm),
		IsNew:           isNew,
		Weekdays:        weekdays,
		WeekdaysPresets: weekdaysPresets,
	})
}

func (h *Handler) handleAlarmsFormPOST(c *gin.Context) {
	// Check if new alarm
	rawAlarmID := c.Param("alarm_id")
	isNew := rawAlarmID == "new"

	// Parse body
	alarmBody := entities.AlarmBody{}
	err := c.MustBindWith(&alarmBody, binding.FormPost)
	if err != nil {
		renderAlarmsFormWithError(c, isNew, alarmBody, fmt.Sprintf("Received invalid data: %v", err))
		return
	}

	// Create new entity
	alarmEntity, err := alarmBody.ToEntity()
	if err != nil {
		renderAlarmsFormWithError(c, isNew, alarmBody, fmt.Sprintf("Failed to parse alarm body into an entity: %v", err))
		return
	}
	if isNew {
		// Enable new alarms by default
		alarmEntity.Enabled = true

		// Create alarm
		_, err := h.alarmService.CreateAlarm(alarmEntity)
		if err != nil {
			renderAlarmsFormWithError(c, isNew, alarmBody, fmt.Sprintf("Failed to add alarm: %v", err))
			return
		}
	} else {
		// Parse ID parameter
		alarmID, err := uuid.Parse(rawAlarmID)
		if err != nil {
			renderAlarmsFormWithError(c, isNew, alarmBody, fmt.Sprintf("Invalid alarm ID %s: %v", rawAlarmID, err))
			return
		}

		// Fetch alarm
		current, err := h.alarmService.GetAlarm(alarmID)
		if err != nil {
			renderAlarmsFormWithError(c, isNew, alarmBody, fmt.Sprintf("Alarm %s not found: %v", rawAlarmID, err))
			return
		}

		// Update alarm
		current.Name = alarmEntity.Name
		current.Hour = alarmEntity.Hour
		current.Minute = alarmEntity.Minute
		current.Days = alarmEntity.Days
		err = h.alarmService.UpdateAlarm(current)
		if err != nil {
			renderAlarmsFormWithError(c, isNew, alarmBody, fmt.Sprintf("Failed to update alarm: %v", err))
			return
		}
	}

	// Upsert successful
	h.redirectWithSuccessMessage(c, "/alarms", "Alarm successfully added/updated")
}

func renderAlarmsFormWithError(c *gin.Context, isNew bool, alarmBody entities.AlarmBody, message string) {
	html(c, http.StatusOK, &entities.AlarmsFormTemplate{
		BaseData: entities.BaseData{
			Title:      "New/Edit alarm",
			ParentPath: "alarms",
			Messages: []entities.Message{{
				Type:    entities.MessageTypeError,
				Content: message,
			}},
		},
		IsNew:           isNew,
		AlarmBody:       alarmBody,
		Weekdays:        weekdays,
		WeekdaysPresets: weekdaysPresets,
	})
}

func (h *Handler) handleAlarmsSetEnabled(c *gin.Context) {
	h.setBooleanAlarmAttribute(c, func(alarm *globalEntities.Alarm, value bool) { alarm.Enabled = value })
}

func (h *Handler) handleAlarmsSetSkipNext(c *gin.Context) {
	h.setBooleanAlarmAttribute(c, func(alarm *globalEntities.Alarm, value bool) { alarm.SkipNext = value })
}

func (h *Handler) setBooleanAlarmAttribute(c *gin.Context, updateAlarmFunc func(alarm *globalEntities.Alarm, value bool)) {
	// Parse alarm ID
	rawAlarmID := c.Param("alarm_id")
	alarmID, err := uuid.Parse(rawAlarmID)
	if err != nil {
		h.redirectWithErrorMessage(c, "/alarms", `Failed to parse invalid alarm ID "%s": %v`, rawAlarmID, err.Error())
		return
	}

	// Fetch boolean value
	rawValue, set := c.GetPostForm("value")
	if !set {
		h.redirectWithErrorMessage(c, "/alarms", `Boolean missing in POST body field "value"`)
		return
	}

	// Parse boolean value
	value, err := strconv.ParseBool(rawValue)
	if err != nil {
		h.redirectWithErrorMessage(c, "/alarms", `Failed to parse boolean value "%s": %v`, rawValue, err.Error())
		return
	}

	// Fetch alarm
	alarm, err := h.alarmService.GetAlarm(alarmID)
	if err != nil {
		h.redirectWithErrorMessage(c, "/alarms", `Failed to get alarm "%s" for update: %v`, alarmID, err.Error())
		return
	}

	// Update alarm
	updateAlarmFunc(&alarm, value)
	err = h.alarmService.UpdateAlarm(alarm)
	if err != nil {
		h.redirectWithErrorMessage(c, "/alarms", `Failed to update alarm "%s": %v`, alarmID, err.Error())
		return
	}

	// Redirect back to overview
	redirect(c, "/alarms")
}
