package rest

import (
	"gopkg.in/go-playground/validator.v9"
)

var validate = validator.New()

func Validate(obj interface{}) error {
	return validate.Struct(obj)
}
