package utils

import (
	"errors"
	"strings"
)

func ValidateName(name string) (string, error) {
	return ValidateRequiredString(name, "name")
}

func ValidateRequiredString(
	value string,
	field string,
) (string, error) {

	value = strings.TrimSpace(value)

	if value == "" {
		return "", errors.New(field + " cannot be empty")
	}

	if len(value) > 200 {
		return "", errors.New(field + " too long")
	}

	return value, nil
}
