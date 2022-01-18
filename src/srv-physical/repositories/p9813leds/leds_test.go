package p9813leds

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_scaleColor(t *testing.T) {
	subtests := []struct {
		color      byte
		brightness byte
		expected   byte
	}{
		{
			color:      0,
			brightness: 0,
			expected:   0,
		},
		{
			color:      255,
			brightness: 255,
			expected:   255,
		},
		{
			color:      255,
			brightness: 0,
			expected:   0,
		},
		{
			color:      0,
			brightness: 255,
			expected:   0,
		},
		{
			color:      255,
			brightness: 128,
			expected:   128,
		},
		{
			color:      200,
			brightness: 64,
			expected:   50,
		},
	}

	for _, subtest := range subtests {
		actual := scaleColor(subtest.color, subtest.brightness)
		require.EqualValues(t, subtest.expected, actual)
	}
}
