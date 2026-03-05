package validator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ParseAndValidate reads JSON from request and validates the struct tags
func ParseAndValidate(r *http.Request, dst interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return fmt.Errorf("invalid JSON: %v", err)
	}

	return ValidateStruct(dst)
}

// ValidateStruct manually triggers validation on a struct
func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	// Format validation errors to be user-friendly
	var errorMsgs []string
	for _, err := range err.(validator.ValidationErrors) {
		field := strings.ToLower(err.Field())
		switch err.Tag() {
		case "required":
			errorMsgs = append(errorMsgs, fmt.Sprintf("%s is required", field))
		case "email":
			errorMsgs = append(errorMsgs, fmt.Sprintf("%s must be a valid email", field))
		case "min":
			errorMsgs = append(errorMsgs, fmt.Sprintf("%s must be at least %s characters", field, err.Param()))
		case "max":
			errorMsgs = append(errorMsgs, fmt.Sprintf("%s must be at most %s characters", field, err.Param()))
		default:
			errorMsgs = append(errorMsgs, fmt.Sprintf("%s failed on %s", field, err.Tag()))
		}
	}

	return fmt.Errorf(strings.Join(errorMsgs, ", "))
}
