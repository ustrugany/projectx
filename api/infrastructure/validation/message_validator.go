package validation

import (
	"errors"
	"fmt"

	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"

	"github.com/ustrugany/projectx/api/service"
)

type messageValidator struct {
	validator  validator.Validate
	translator ut.Translator
}

func CreateMessageValidator(validator validator.Validate, translator ut.Translator) service.MassageValidator {
	return messageValidator{validator: validator, translator: translator}
}

func (v messageValidator) ValidateForCreate(title, content, email string, magicNumber int) (map[string]string, error) {
	var violations map[string]string
	message := struct {
		Title       string `json:"title" validate:"required,gte=0,lte=140"`
		Email       string `json:"email" validate:"required,email"`
		Content     string `json:"content" validate:"required,gte=0,lte=1500"`
		MagicNumber int    `json:"magic_number" validate:"required,gte=0,lte=999"`
	}{
		Title:       title,
		Email:       email,
		Content:     content,
		MagicNumber: magicNumber,
	}

	err := v.validator.Struct(&message)
	if err != nil {
		switch e := err.(type) {
		case validator.ValidationErrors:
			violations = make(map[string]string, len(e))
			for _, validationError := range e {
				violations[validationError.Field()] = validationError.Translate(v.translator)
			}
			return violations, nil
		case *validator.InvalidValidationError:
			return violations, fmt.Errorf("failed to validate: %w", e)
		default:
			return violations, errors.New("failed to validate")
		}
	}

	return violations, nil
}
