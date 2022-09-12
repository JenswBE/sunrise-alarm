// Code generated by mockery v2.12.1. DO NOT EDIT.

package mocks

import (
	entities "github.com/JenswBE/sunrise-alarm/src/services/physical/entities"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// Leds is an autogenerated mock type for the Leds type
type Leds struct {
	mock.Mock
}

// GetColorAndBrightness provides a mock function with given fields:
func (_m *Leds) GetColorAndBrightness() (entities.PresetColor, byte) {
	ret := _m.Called()

	var r0 entities.PresetColor
	if rf, ok := ret.Get(0).(func() entities.PresetColor); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(entities.PresetColor)
	}

	var r1 byte
	if rf, ok := ret.Get(1).(func() byte); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(byte)
	}

	return r0, r1
}

// SetColorAndBrightness provides a mock function with given fields: color, brightness
func (_m *Leds) SetColorAndBrightness(color entities.PresetColor, brightness byte) {
	_m.Called(color, brightness)
}

// NewLeds creates a new instance of Leds. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewLeds(t testing.TB) *Leds {
	mock := &Leds{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
