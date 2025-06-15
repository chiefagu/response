package response

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// StandardResponse defines the uniform response structure
type standardResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// SendSuccessResponse formats and sends a success response
func SendSuccessResponse(c echo.Context, data any) error {
	response := standardResponse{
		Status: "success",
		Data:   data,
	}
	return c.JSON(http.StatusOK, response)
}

func SendNoContentResponse(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

// SendErrorResponse formats and sends an error response
func SendErrorResponse(c echo.Context, code int, message string, data any) error {
	response := standardResponse{
		Status:  "error",
		Message: message,
		Data:    data,
	}
	return c.JSON(code, response)
}

// ValidationErrorResponse formats validation errors into a standard response
func ValidationErrorResponse(c echo.Context, err error) error {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return SendErrorResponse(c, http.StatusBadRequest, "validation failed", map[string]string{"error": "Invalid request format"})
	}

	errorMessages := make(map[string]string)
	for _, vErr := range validationErrors {
		fieldName := vErr.Field()
		customMsg := getCustomValidationMessage(vErr.Tag(), fieldName)
		errorMessages[fieldName] = customMsg
	}

	return SendErrorResponse(c, http.StatusBadRequest, "validation failed", errorMessages)
}

// getCustomValidationMessage returns a user-friendly error message based on the validation tag
func getCustomValidationMessage(tag, field string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("'%s' is required and cannot be empty", field)
	case "email":
		return fmt.Sprintf("'%s' must be a valid email address", field)
	case "min":
		return fmt.Sprintf("'%s' must meet the minimum length requirement", field)
	case "category_enum":
		return fmt.Sprintf("'%s' must be either 'data' or 'airtime'", field)
	default:
		return fmt.Sprintf("'%s' is invalid due to '%s' validation rule", field, tag)
	}
}
