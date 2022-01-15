package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func setupGinTest(t *testing.T, method, path string, params gin.Params, body []byte) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.ReleaseMode)
	c, _ := gin.CreateTestContext(w)
	bodyReader := bytes.NewReader(body)
	var err error
	c.Request, err = http.NewRequest(method, path, bodyReader)
	require.NoError(t, err)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = params
	return c, w
}

// RequireEqualJSON unmarshals the body from the provided recorder to the same type as expected.
// Next, it asserts this result against the expected value.
func RequireEqualJSON(t *testing.T, expected interface{}, recorder *httptest.ResponseRecorder) {
	actual := reflect.New(reflect.TypeOf(expected))
	body := recorder.Body.Bytes()
	err := json.Unmarshal(body, actual.Interface())
	require.NoErrorf(t, err, `Response body: %s`, string(body))
	// Has better support for time.Time than require.Equal
	// See https://github.com/stretchr/testify/issues/1078 for more info
	require.Empty(t, cmp.Diff(expected, reflect.Indirect(actual).Interface()))
}
