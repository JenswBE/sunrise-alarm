package handler

import (
	"reflect"

	"github.com/JenswBE/sunrise-alarm/src/srv-physical/api/openapi"
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/entities"
	"github.com/rs/zerolog/log"
)

// ErrToResponse checks if the provided error is a GoComError.
// If yes, status and embedded error message are returned.
// If no, status is 500 and provided error message are returned.
func ErrToResponse(e error) (int, *entities.Error) {
	if err, ok := e.(*entities.Error); ok {
		return err.Status, err
	}
	log.Warn().Err(e).Stringer("error_type", reflect.TypeOf(e)).Msg("API received an non-GoComError error")
	return 500, entities.NewError(500, openapi.ERRORCODE_UNKNOWN_ERROR, "", e).(*entities.Error)
}
