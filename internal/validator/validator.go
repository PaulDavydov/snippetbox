package validator

import (
	"slices"
	"strings"
	"unicode/utf8"
)

// Validaotor struct that contains a map of validation errors for our form fields
type Validator struct {
	FieldErrors map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

func (v *Validator) AddFieldError(key, message string) {
	// check to see if map has been initialized
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// adds an error message to the FieldErrors map
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// return true if value is not an empty string
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// return true if value contains no more than n characters
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// returns true if value is in a list of permitted integers
func PermittedInt(value int, permittedValues ...int) bool {
	return slices.Contains(permittedValues, value)
}
