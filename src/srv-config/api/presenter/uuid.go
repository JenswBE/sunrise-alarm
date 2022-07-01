package presenter

import (
	"github.com/JenswBE/sunrise-alarm/src/srv-config/api/openapi"
	"github.com/JenswBE/sunrise-alarm/src/srv-config/entities"
	"github.com/google/uuid"
)

func ParseUUID(input string) (uuid.UUID, error) {
	// Parse ID
	id, err := uuid.Parse(input)
	if err != nil {
		return uuid.Nil, entities.NewError(400, openapi.ERRORCODE_INVALID_UUID, input, nil)
	}

	// Parse successful
	return id, nil
}
