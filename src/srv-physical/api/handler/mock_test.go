package handler

import (
	"net/http"
	"testing"

	mocks "github.com/JenswBE/sunrise-alarm/src/srv-physical/mocks/usecases/mqtt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupMockHandlerTest() (*MockHandler, *mocks.Usecase) {
	usecase := &mocks.Usecase{}
	handler := NewMockHandler(usecase)
	return handler, usecase
}

func Test_Mock_buttonPressed_Success(t *testing.T) {
	// Setup test
	c, r := setupGinTest(t, "POST", "", nil, nil)
	mockHandler, usecaseMock := setupMockHandlerTest()
	usecaseMock.On("PublishButtonPressed", mock.Anything).Return(nil)

	// Call handler
	mockHandler.buttonPressed(c)

	// Assert result
	require.Equal(t, http.StatusNoContent, r.Code)
	usecaseMock.AssertNumberOfCalls(t, "PublishButtonPressed", 1)
}

func Test_Mock_buttonLongPressed_Success(t *testing.T) {
	// Setup test
	c, r := setupGinTest(t, "POST", "", nil, nil)
	mockHandler, usecaseMock := setupMockHandlerTest()
	usecaseMock.On("PublishButtonLongPressed", mock.Anything).Return(nil)

	// Call handler
	mockHandler.buttonLongPressed(c)

	// Assert result
	require.Equal(t, http.StatusNoContent, r.Code)
	usecaseMock.AssertNumberOfCalls(t, "PublishButtonLongPressed", 1)
}
