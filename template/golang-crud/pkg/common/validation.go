package common

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

func Validate(data interface{}, validations map[string]struct {
	Message   string
	Validator validator.Func
},
) validator.ValidationErrorsTranslations {
	eng := en.New()
	uni := ut.New(eng, eng)
	trans, _ := uni.GetTranslator("en")

	validate := validator.New()
	if err := en_translations.RegisterDefaultTranslations(validate, trans); err != nil {
		return nil, err
	}

	// Get the JSON tags in the keys instead of struct tags
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Register the custom validation functions
	for key, item := range validations {
		if err := validate.RegisterValidation(key, item.Validator); err != nil {
			return nil, err
		}

		if err := validate.RegisterTranslation(key, trans, func(ut ut.Translator) error {
			return ut.Add(key, item.Message, true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			field := fe.Field()
			// This might not get the value for every data type
			value := fmt.Sprintf("%s", fe.Value())

			t, _ := ut.T(key, field, value)

			return t
		}); err != nil {
			return nil, err
		}
	}

	if err := validate.Struct(data); err != nil {
		errs := make(validator.ValidationErrorsTranslations)

		// Build our own version of the Translate function to get the JSON key
		for _, err := range err.(validator.ValidationErrors) {
			errs[err.Field()] = err.Translate(trans)
		}

		return errs
	}

	return nil
}
