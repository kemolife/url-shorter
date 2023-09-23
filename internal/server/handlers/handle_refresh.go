package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type RefreshResponse struct {
	AccessToken  string
	RefreshToken string
}

func HandleRefresh() echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.JSON(
			http.StatusOK,
			RefreshResponse{
				"", "",
			},
		)
	}
}
