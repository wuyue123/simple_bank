package api

import (
	"github.com/go-playground/validator/v10"
	"pxsemic.com/simplebank/util"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	currency := fl.Field().String()
	return util.IsSupportedCurrency(currency)
}
