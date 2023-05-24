// Package svalid содержитметоды-обертки для валидации структур.
package svalid

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

// SValid структура валидации.
type SValid struct {
	validator *validator.Validate
}

// New конструктор.
// Добавляется новый тип валидируемых данных "enum".
func New() *SValid {
	v := validator.New()
	_ = v.RegisterValidation("enum", func(fl validator.FieldLevel) bool {
		vType, ok := fl.Field().Interface().(interface{ IsValid() bool })
		if !ok {
			return false
		}
		return vType.IsValid()
	})
	return &SValid{
		validator: v,
	}
}

// Validate метод валидации структуры.
func (s SValid) Validate(v any) error {
	err := s.validator.Struct(v)
	if err != nil {
		var errorList []string
		for _, err := range err.(validator.ValidationErrors) {
			errorList = append(errorList, err.Error())
		}
		return errors.New(strings.Join(errorList, "\n"))
	}
	return nil
}
