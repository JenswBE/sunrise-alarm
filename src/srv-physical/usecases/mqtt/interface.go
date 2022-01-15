package mqtt

import (
	"context"

	"github.com/JenswBE/sunrise-alarm/src/srv-physical/entities"
)

type Usecase interface {
	PublishButtonPressed(ctx context.Context) error
	PublishButtonLongPressed(ctx context.Context) error
	PublishTempHumidUpdated(ctx context.Context, e entities.TempHumidReading) error
}
