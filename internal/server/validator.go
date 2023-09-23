package server

import (
	"github.com/go-playground/validator/v10"
	"net/http"

	"github.com/labstack/echo/v4"
)

type RequestValidator struct {
	v *validator.Validate
}

func NewValidator() *RequestValidator {
	return &RequestValidator{
		v: validator.New(),
	}
}

func (v *RequestValidator) Validate(i any) error {
	if err := v.v.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}
