package api

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateStruct comprueba las etiquetas de validación de los DTOs
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}