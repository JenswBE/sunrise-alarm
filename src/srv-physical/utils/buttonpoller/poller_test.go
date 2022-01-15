package buttonpoller_test

import (
	"testing"
	"time"

	"github.com/JenswBE/sunrise-alarm/src/srv-physical/repositories/mockbutton"
	"github.com/JenswBE/sunrise-alarm/src/srv-physical/utils/buttonpoller"
	"github.com/stretchr/testify/require"
)

func TestDebounceShort(t *testing.T) {
	// Setup test
	t.Parallel()
	mockButton := mockbutton.MockButton{}
	notifyChannel := make(chan buttonpoller.ButtonPress)
	button := buttonpoller.NewButtonPoller(&mockButton, notifyChannel)
	_ = button

	// Test short press
	mockButton.Pressed = true
	time.Sleep(2 * buttonpoller.DebounceDuration)
	mockButton.Pressed = false
	time.Sleep(2 * buttonpoller.DebounceDuration)
	select {
	case notif := <-notifyChannel:
		require.Equal(t, buttonpoller.ButtonPressShort, notif)
	default:
		require.FailNow(t, "No short button press detected")
	}
}

func TestDebounceLong(t *testing.T) {
	// Setup test
	t.Parallel()
	mockButton := mockbutton.MockButton{}
	notifyChannel := make(chan buttonpoller.ButtonPress)
	button := buttonpoller.NewButtonPoller(&mockButton, notifyChannel)
	_ = button

	// Test long press
	mockButton.Pressed = true
	time.Sleep(buttonpoller.DebounceDuration*2 + buttonpoller.LongPressDuration)
	mockButton.Pressed = false
	time.Sleep(2 * buttonpoller.DebounceDuration)
	select {
	case notif := <-notifyChannel:
		require.Equal(t, buttonpoller.ButtonPressLong, notif)
	default:
		require.FailNow(t, "No long button press detected")
	}
}

func TestDebounceNoPress(t *testing.T) {
	// Setup test
	t.Parallel()
	mockButton := mockbutton.MockButton{}
	notifyChannel := make(chan buttonpoller.ButtonPress)
	button := buttonpoller.NewButtonPoller(&mockButton, notifyChannel)
	_ = button

	// Test debouncing
	mockButton.Pressed = true
	time.Sleep(buttonpoller.DebounceDuration)
	mockButton.Pressed = false
	time.Sleep(buttonpoller.DebounceDuration)
	mockButton.Pressed = true
	time.Sleep(buttonpoller.DebounceDuration)
	mockButton.Pressed = false
	time.Sleep(buttonpoller.DebounceDuration)
	require.Empty(t, notifyChannel)
}
