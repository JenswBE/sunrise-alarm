// Based on https://github.com/metakeule/loop

package loop

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoop(t *testing.T) {
	testCases := map[string]struct {
		givenData    string
		expectedData string
	}{
		"shorter length": {
			givenData:    "12345",
			expectedData: "123",
		},
		"exact length": {
			givenData:    "12345",
			expectedData: "12345",
		},
		"longer partial length": {
			givenData:    "12345",
			expectedData: "1234512345123",
		},
	}
	for desc, tC := range testCases {
		t.Run(desc, func(t *testing.T) {
			loop := New([]byte(tC.givenData))
			buffer := make([]byte, len(tC.expectedData))
			length, _ := loop.Read(buffer)
			require.Equal(t, tC.expectedData, string(buffer))
			require.Equal(t, len(tC.expectedData), length)
		})
	}
}

func TestLoopNewNil(t *testing.T) {
	require.Panics(t, func() { New(nil) })
}

func TestLoopNewEmpty(t *testing.T) {
	require.Panics(t, func() { New([]byte{}) })
}
