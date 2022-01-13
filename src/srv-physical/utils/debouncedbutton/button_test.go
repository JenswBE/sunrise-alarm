package debouncedbutton_test

import (
	"testing"
	"time"

	"github.com/JenswBE/sunrise-alarm/src/srv-physical/repositories/mockbutton"
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/utils/debouncedbutton"
	"github.com/stretchr/testify/require"
)

func TestDebounceShort(t *testing.T) {
	// Setup test
	t.Parallel()
	mockButton := mockbutton.MockButton{}
	notifyChannel := make(chan debouncedbutton.ButtonPress)
	button := debouncedbutton.NewDebouncedButton(&mockButton, notifyChannel)
	_ = button

	// Test short press
	mockButton.Pressed = true
	time.Sleep(2 * debouncedbutton.DebounceDuration)
	mockButton.Pressed = false
	time.Sleep(2 * debouncedbutton.DebounceDuration)
	select {
	case notif := <-notifyChannel:
		require.Equal(t, debouncedbutton.ButtonPressShort, notif)
	default:
		require.FailNow(t, "No short button press detected")
	}
}

func TestDebounceLong(t *testing.T) {
	// Setup test
	t.Parallel()
	mockButton := mockbutton.MockButton{}
	notifyChannel := make(chan debouncedbutton.ButtonPress)
	button := debouncedbutton.NewDebouncedButton(&mockButton, notifyChannel)
	_ = button

	// Test long press
	mockButton.Pressed = true
	time.Sleep(debouncedbutton.DebounceDuration*2 + debouncedbutton.LongPressDuration)
	mockButton.Pressed = false
	time.Sleep(2 * debouncedbutton.DebounceDuration)
	select {
	case notif := <-notifyChannel:
		require.Equal(t, debouncedbutton.ButtonPressLong, notif)
	default:
		require.FailNow(t, "No long button press detected")
	}
}

func TestDebounceNoPress(t *testing.T) {
	// Setup test
	t.Parallel()
	mockButton := mockbutton.MockButton{}
	notifyChannel := make(chan debouncedbutton.ButtonPress)
	button := debouncedbutton.NewDebouncedButton(&mockButton, notifyChannel)
	_ = button

	// Test debouncing
	mockButton.Pressed = true
	time.Sleep(debouncedbutton.DebounceDuration)
	mockButton.Pressed = false
	time.Sleep(debouncedbutton.DebounceDuration)
	mockButton.Pressed = true
	time.Sleep(debouncedbutton.DebounceDuration)
	mockButton.Pressed = false
	time.Sleep(debouncedbutton.DebounceDuration)
	require.Empty(t, notifyChannel)
}
