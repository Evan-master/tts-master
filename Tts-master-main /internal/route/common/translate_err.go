package common

import (
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
)

var (
	uni   *ut.UniversalTranslator
	trans ut.Translator
)

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		en2 := en.New()
		uni = ut.New(en2, en2)
		trans, _ = uni.GetTranslator("en")
		err := entranslations.RegisterDefaultTranslations(v, trans)
		if err != nil {
			panic(err)
		}
	}
}

func TranslateErr(err error) string {
	if err == nil {
		return ""
	}
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err.Error()
	}
	var builder strings.Builder
	for _, e := range errs {
		builder.WriteString(e.Translate(trans))
		builder.WriteByte('\n')
	}
	return builder.String()
}
