package ctxhelper

import (
	"fmt"
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	idtranslations "github.com/go-playground/validator/v10/translations/id"
)

var validate *validator.Validate

type language struct {
	name       string
	translator func() locales.Translator
	register   func(v *validator.Validate, t ut.Translator)
}

var (
	validatorTranslators = map[string]ut.Translator{}

	validatorLanguages = []language{
		{
			name:       "en",
			translator: en.New,
			register: func(v *validator.Validate, t ut.Translator) {
				if err := entranslations.RegisterDefaultTranslations(v, t); err != nil {
					panic(fmt.Errorf("failed to register validator translations for %s: %w", "en", err))
				}
			},
		},
		{
			name:       "id",
			translator: id.New,
			register: func(v *validator.Validate, t ut.Translator) {
				if err := idtranslations.RegisterDefaultTranslations(v, t); err != nil {
					panic(fmt.Errorf("failed to register validator translations for %s: %w", "id", err))
				}
			},
		},
	}
)

func init() {
	validate = validator.New()

	for _, lang := range validatorLanguages {
		name := lang.name
		translator := lang.translator()
		universalTranslator := ut.New(translator, translator)
		tr, found := universalTranslator.GetTranslator(name)
		if !found {
			panic(fmt.Errorf("failed to get validation translator for %s", name))
		}
		lang.register(validate, tr)
		validatorTranslators[name] = tr
	}
}
