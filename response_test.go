package response_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chiefagu/response"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// Helper function to create an Echo context
func createTestContext(method, path, body string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}

// Test success response
func TestSendSuccessResponse(t *testing.T) {
	c, rec := createTestContext(http.MethodGet, "/", "")

	data := map[string]string{"message": "test success"}
	err := response.SendSuccessResponse(c, data)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"status":"success","data":{"message":"test success"}}`, rec.Body.String())
}

// Test no content response
func TestSendNoContentResponse(t *testing.T) {
	c, rec := createTestContext(http.MethodPost, "/", "")

	err := response.SendNoContentResponse(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

// Test error response
func TestSendErrorResponse(t *testing.T) {
	c, rec := createTestContext(http.MethodGet, "/", "")

	message := "test error"
	err := response.SendErrorResponse(c, http.StatusBadRequest, message, nil)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, `{"status":"error","message":"test error"}`, rec.Body.String())
}

// Test validation error response
func TestValidationErrorResponse(t *testing.T) {
	type TestStruct struct {
		Field1 string `validate:"required"`
		Field2 string `validate:"email"`
	}

	validate := validator.New()
	testData := TestStruct{} // Missing required fields to trigger validation error
	err := validate.Struct(testData)

	assert.Error(t, err)

	c, rec := createTestContext(http.MethodGet, "/", "")

	_ = response.ValidationErrorResponse(c, err)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	expectedJSON := `{
        "status": "error",
        "message": "validation failed",
        "data": {
            "Field1": "'Field1' is required and cannot be empty",
            "Field2": "'Field2' must be a valid email address"
        }
    }`

	assert.JSONEq(t, expectedJSON, rec.Body.String())
}
