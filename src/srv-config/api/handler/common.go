package handler

import (
	"reflect"

	"github.com/JenswBE/sunrise-alarm/src/srv-config/api/openapi"
	"github.com/JenswBE/sunrise-alarm/src/srv-config/api/presenter"
	"github.com/JenswBE/sunrise-alarm/src/srv-config/entities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// ParseIDParam tries to parse parameter with the given name as an UUID.
// On failure, an error is set on the Gin context.
//
//   id, ok := ParseIDParam(c, "id")
//   if !ok {
// 	   return // Response already set on Gin context
//   }
func ParseIDParam(c *gin.Context, name string) (uuid.UUID, bool) {
	// Parse param
	pID, ok := c.Params.Get(name)
	if !ok {
		err := entities.NewError(400, openapi.ERRORCODE_PARAMETER_MISSING, name, nil)
		c.JSON(ErrToResponse(err))
		return uuid.Nil, false
	}

	// Parse ID
	id, err := presenter.ParseUUID(pID)
	if err != nil {
		c.JSON(ErrToResponse(err))
		return uuid.Nil, false
	}

	// Parse successful
	return id, true
}

// ErrToResponse checks if the provided error is an entity of type Error.
// If yes, status and embedded error message are returned.
// If no, status is 500 and provided error message are returned.
func ErrToResponse(e error) (int, *entities.Error) {
	if err, ok := e.(*entities.Error); ok {
		return err.Status, err
	}
	log.Warn().Err(e).Stringer("error_type", reflect.TypeOf(e)).Msg("API received a plain error")
	return 500, entities.NewError(500, openapi.ERRORCODE_UNKNOWN_ERROR, "", e).(*entities.Error)
}
