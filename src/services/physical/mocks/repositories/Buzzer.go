// Code generated by mockery v2.12.1. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// Buzzer is an autogenerated mock type for the Buzzer type
type Buzzer struct {
	mock.Mock
}

// Off provides a mock function with given fields:
func (_m *Buzzer) Off() {
	_m.Called()
}

// On provides a mock function with given fields:
func (_m *Buzzer) On() {
	_m.Called()
}

// NewBuzzer creates a new instance of Buzzer. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewBuzzer(t testing.TB) *Buzzer {
	mock := &Buzzer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}