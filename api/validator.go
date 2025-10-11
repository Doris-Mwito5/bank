package api

import (
	"github.com/Doris-Mwito5/simple-bank/internal/utils"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return utils.IsSupported(currency)
	}
	return false
}
