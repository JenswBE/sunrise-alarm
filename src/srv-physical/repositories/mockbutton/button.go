package mockbutton

import (
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/repositories"
)

var _ repositories.Button = &MockButton{}

type MockButton struct {
	Pressed bool
}

func NewMockButton() *MockButton {
	return &MockButton{}
}

func (b *MockButton) IsPressed() bool {
	return b.Pressed
}
