package request

import (
	"errors"
	"fmt"
	"learn/internal/pkg/response"
	"log/slog"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// BindJSONOrError attempts to bind the request body to the given input struct.
// If binding fails, it logs the error, sends a consistent validation response,
// and returns false. If binding succeeds, it returns true.
func BindJSONOrError(c *gin.Context, input interface{}, logger *slog.Logger, action string) bool {
	if err := c.ShouldBindJSON(input); err != nil {
		logger.Warn("failed to bind JSON for "+action, slog.String("error", err.Error()))
		response.SendValidationError(c, validationDetails(input, err))
		return false
	}
	return true
}

func validationDetails(input interface{}, err error) []response.ValidationErrorDetail {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		details := make([]response.ValidationErrorDetail, 0, len(validationErrors))
		for _, fieldErr := range validationErrors {
			details = append(details, response.ValidationErrorDetail{
				Field:   jsonFieldName(input, fieldErr.StructNamespace()),
				Message: validationMessage(fieldErr),
			})
		}
		return details
	}

	return []response.ValidationErrorDetail{{
		Field:   "body",
		Message: "request body must be valid JSON and match the expected format",
	}}
}

func validationMessage(fieldErr validator.FieldError) string {
	switch fieldErr.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address"
	case "len":
		return fmt.Sprintf("must be exactly %s characters", fieldErr.Param())
	case "min":
		return fmt.Sprintf("must be at least %s", fieldErr.Param())
	case "max":
		return fmt.Sprintf("must be at most %s", fieldErr.Param())
	default:
		return fmt.Sprintf("failed validation rule %q", fieldErr.Tag())
	}
}

func jsonFieldName(input interface{}, namespace string) string {
	parts := strings.Split(namespace, ".")
	if len(parts) <= 1 {
		return strings.ToLower(namespace)
	}

	t := reflect.TypeOf(input)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	jsonParts := make([]string, 0, len(parts)-1)
	for _, part := range parts[1:] {
		fieldName := strings.Split(part, "[")[0]

		if t.Kind() == reflect.Slice {
			t = t.Elem()
			for t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
		}

		if t.Kind() != reflect.Struct {
			jsonParts = append(jsonParts, strings.ToLower(fieldName))
			continue
		}

		field, ok := t.FieldByName(fieldName)
		if !ok {
			jsonParts = append(jsonParts, strings.ToLower(fieldName))
			continue
		}

		jsonName := strings.Split(field.Tag.Get("json"), ",")[0]
		if jsonName == "" || jsonName == "-" {
			jsonName = strings.ToLower(field.Name)
		}
		jsonParts = append(jsonParts, jsonName)

		t = field.Type
		for t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
	}

	return strings.Join(jsonParts, ".")
}
