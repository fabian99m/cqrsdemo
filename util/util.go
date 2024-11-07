package util

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

var (
	structValidator = validator.New(validator.WithRequiredStructEnabled(), func(v *validator.Validate) {
		v.RegisterValidation("notblank", validators.NotBlank)
	})
)

func As[T error](err error) (*T, bool) {
	var as T
	ok := errors.As(err, &as)
	return &as, ok
}

func IsType[T error](err error) bool {
	_, ok := As[T](err)
	return ok
}

func UnmarshalTo[T any](data []byte) (dest *T, err error) {
	dest = new(T)
	err = json.Unmarshal(data, dest)
	if err != nil {
		return
	}
	return
}

func ValidateStruct(s any) error {
	return structValidator.Struct(s)
}

func GetValidations(err error) []string {
	validations := []string{}
	for _, err := range err.(validator.ValidationErrors) {
		validations = append(validations, fmt.Sprintf("field: %s - error: %v", err.Field(), err.Error()))
	}
	return validations
}
