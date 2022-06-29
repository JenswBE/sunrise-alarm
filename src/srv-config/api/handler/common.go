package handler

import (
	"reflect"

	"github.com/JenswBE/sunrise-alarm/src/srv-config/api/openapi"
	"github.com/JenswBE/sunrise-alarm/src/srv-config/entities"
	"github.com/rs/zerolog/log"
)

// errToResponse checks if the provided error is an entity of type Error.
// If yes, status and embedded error message are returned.
// If no, status is 500 and provided error message are returned.
func errToResponse(e error) (int, *entities.Error) {
	if err, ok := e.(*entities.Error); ok {
		return err.Status, err
	}
	log.Warn().Err(e).Stringer("error_type", reflect.TypeOf(e)).Msg("API received a plain error")
	return 500, entities.NewError(500, openapi.ERRORCODE_UNKNOWN_ERROR, "", e).(*entities.Error)
}
