package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"shorter-url/internal/auth"
	"shorter-url/internal/config"
	"shorter-url/internal/model"
	"time"

	"github.com/labstack/echo/v4"
)

type callbackProvider interface {
	GitHubAuthCallback(ctx context.Context, sessionCode string) (*model.User, error)
}

func HandleGitHubAuthCallback(cbProvider callbackProvider) echo.HandlerFunc {
	// TODO: add tests
	return func(c echo.Context) error {
		sessionCode := c.QueryParam("code")
		if sessionCode == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "missing code")
		}

		user, err := cbProvider.GitHubAuthCallback(c.Request().Context(), sessionCode)
		jwt, err := auth.MakeJWT(*user, time.Minute*15)
		jwtRefresh, err := auth.MakeJWT(*user, time.Hour*24)
		if err != nil {
			log.Printf("error handling github auth callback: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		redirectURL := fmt.Sprintf("%s/auth/token.html?token=%s&refresh-token=%s", config.Get().Address, jwt, jwtRefresh)
		return c.Redirect(http.StatusMovedPermanently, redirectURL)
	}
}
