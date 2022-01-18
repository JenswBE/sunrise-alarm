package entities

import (
	"fmt"

	"github.com/JenswBE/sunrise-alarm/src/srv-physical/api/openapi"
)

// Error allows to bundle a status with the original error.
// This allows to fine-grained response codes at the API level.
type Error struct {
	// HTTP status code
	Status int `json:"status"`

	// Original error
	Err error `json:"-"`

	// Error code
	Code string `json:"code"`

	// Human-readable description of the error
	Message string `json:"message"`

	// Optional - On which object to error occurred
	Instance string `json:"instance"`
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%d - %s - %s - %s", e.Status, e.Message, e.Instance, e.Err.Error())
	}
	return fmt.Sprintf("%d - %s - %s", e.Status, e.Message, e.Instance)
}

// NewError returns a new GoComError
func NewError(status int, code openapi.ErrorCode, instance string, err error) error {
	return &Error{
		Status:   status,
		Err:      err,
		Code:     string(code),
		Message:  translateCodeToMessage(code),
		Instance: instance,
	}
}

func translateCodeToMessage(code openapi.ErrorCode) string {
	switch code {
	case openapi.ERRORCODE_BRIGHTNESS_OUT_OF_RANGE:
		return `Brightness must be between 0 and 255`
	case openapi.ERRORCODE_UNKNOWN_ERROR:
		return `An unknown error occurred`
	}
	return "" // Covered by exhaustive check
}
