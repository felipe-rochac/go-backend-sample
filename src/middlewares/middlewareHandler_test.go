package middlewares

import (
	"backend-sample/common"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMiddlewareHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Test formatRespose JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		response := gin.H{"message": "success"}
		c.Set("response", response)

		formatResponse(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message": "success"}`, w.Body.String())
	})

	t.Run("Test formatRespose YAML", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest("GET", "/", nil)
		c.Request.Header.Set("Accept", "application/x-yaml")

		response := gin.H{"message": "success"}
		c.Set("response", response)

		formatResponse(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/x-yaml", w.Header().Get("Content-Type"))
		assert.Contains(t, w.Body.String(), "message: success")
	})

	t.Run("Test handleError with BackendError", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		ErrorCodeKey = KeyValue{Key: "testkey", Value: "testvalue"}
		backendError := &common.BackendError{
			Code:       http.StatusBadRequest,
			Identifier: "test_identifier",
		}
		c.Error(backendError)

		handleError(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "error_code")
		assert.Contains(t, w.Body.String(), "An error occurred")
	})

	t.Run("Test handleError with generic error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Error(assert.AnError)

		handleError(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), assert.AnError.Error())
	})
}
