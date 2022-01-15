package mockbutton

import (
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/repositories"
)

var _ repositories.Button = &MockButton{}

type MockButton struct {
	Pressed bool
}

func (b *MockButton) IsPressed() bool {
	return b.Pressed
}
