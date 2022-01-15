package handler

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func setupMockHandlerTest() *MockHandler {
	return NewMockHandler()
}

func Test_Mock_buttonPressed_Success(t *testing.T) {
	// Setup test
	c, r := setupGinTest(t, "POST", "", nil, nil)
	mockHandler := setupMockHandlerTest()

	// Call handler
	mockHandler.buttonPressed(c)

	// Assert result
	require.Equal(t, http.StatusNoContent, r.Code)
}

func Test_Mock_buttonLongPressed_Success(t *testing.T) {
	// Setup test
	c, r := setupGinTest(t, "POST", "", nil, nil)
	mockHandler := setupMockHandlerTest()

	// Call handler
	mockHandler.buttonLongPressed(c)

	// Assert result
	require.Equal(t, http.StatusNoContent, r.Code)
}
