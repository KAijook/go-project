package api

import "github.com/go-playground/validator/v10"

var validCurrencies = map[string]bool{
	"USD": true,
	"EUR": true,
	"VND": true,
}

// Custom Currency Validator Function
var ValidCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return validCurrencies[currency]
	}
	return false
}
