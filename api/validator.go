package api

import (
	"github.com/go-playground/validator/v10"

	"github.com/JekaTka/services-api/util"
)

var validServiceSort validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if sort, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedSort(sort)
	}
	return false
}
