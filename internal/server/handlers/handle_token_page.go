package handlers

import (
	"bytes"
	"embed"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-github/v48/github"
	"github.com/labstack/echo/v4"
	"html/template"
	"log"
	"net/http"
	"shorter-url/internal/config"
	"shorter-url/internal/model"
)

//go:embed static/*
var static embed.FS

func HandleTokenPage() echo.HandlerFunc {
	tmpl, err := template.ParseFS(static, "static/token.html")
	if err != nil {
		log.Fatalf("error parsing token.html template: %v", err)
	}

	type templateData struct {
		Token           string
		RefreshToken    string
		GitHubUsername  string
		GitHubAvatarURL string
	}

	type request struct {
		Token        string `query:"token"`
		RefreshToken string `query:"refresh-token"`
	}

	return func(c echo.Context) error {
		var req request
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
		}

		tokenObj, err := jwt.ParseWithClaims(
			req.Token,
			&model.UserClaims{},
			func(_ *jwt.Token) (any, error) { return []byte(config.Get().Auth.JWTSecretKey), nil },
		)

		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		if !tokenObj.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		claims, ok := tokenObj.Claims.(*model.UserClaims)
		if !ok {
			log.Printf("error asserting claims to *model.UserClaims")
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		ghClient := github.NewClient(nil)
		ghUser, _, err := ghClient.Users.Get(c.Request().Context(), claims.User.GitHubLogin)
		if err != nil {
			log.Printf("error getting github user: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		var (
			buf  bytes.Buffer
			data = templateData{
				Token:           req.Token,
				RefreshToken:    req.RefreshToken,
				GitHubUsername:  claims.User.GitHubLogin,
				GitHubAvatarURL: ghUser.GetAvatarURL(),
			}
		)

		if err := tmpl.Execute(&buf, data); err != nil {
			log.Printf("error executing token.html template: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		return c.HTML(http.StatusOK, buf.String())
	}
}

func HandleStatic() echo.HandlerFunc {
	return echo.WrapHandler(http.FileServer(http.FS(static)))
}
