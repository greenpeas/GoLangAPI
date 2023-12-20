package validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

type Validate struct {
	validate *validator.Validate
}

func NewValidator() *Validate {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		if name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]; name != "-" {
			return name
		}

		if name := strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]; name != "-" {
			return name
		}

		return ""
	})

	return &Validate{validate}
}

func (v *Validate) Struct(s interface{}) map[string]string {
	err := v.validate.Struct(s)

	if err == nil {
		return nil
	}

	errors := make(map[string]string)

	for _, err := range err.(validator.ValidationErrors) {
		errors[err.Field()] = getErrorText(err)
		//fmt.Println(err.Namespace()) // can differ when a custom TagNameFunc is registered or
		//fmt.Println(err.Field())     // by passing alt name to ReportError like below
		//fmt.Println(err.StructNamespace())
		//fmt.Println(err.StructField())
		//fmt.Println(err.Tag())
		//fmt.Println(err.ActualTag())
		//fmt.Println(err.Kind())
		//fmt.Println(err.Type())
		//fmt.Println(err.Value())
		//fmt.Println(err.Param())
		//fmt.Println()
	}

	// from here you can create your own error messages in whatever language you wish
	return errors
}

func getErrorText(err validator.FieldError) string {
	switch err.Tag() {
	case "max":
		return fmt.Sprintf("Максимум %v", err.Param())
	case "min":
		return fmt.Sprintf("Минимум %v", err.Param())
	case "required":
		return fmt.Sprintf("Обязательно к заполнению")
	default:
		return fmt.Sprintf("%s is %v", err.Tag(), err.Param())
	}
}
