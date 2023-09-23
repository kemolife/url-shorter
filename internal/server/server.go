package server

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"shorter-url/internal/config"
	"shorter-url/internal/model"
	"shorter-url/internal/server/handlers"
	"shorter-url/internal/shorter"

	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"shorter-url/internal/auth"
)

type CloseFunc func(context.Context) error

type Server struct {
	e         *echo.Echo
	shortener *shorter.Service
	auth      *auth.Service
	closers   []CloseFunc
}

func New(shortener *shorter.Service, auth *auth.Service) *Server {
	s := &Server{
		shortener: shortener,
		auth:      auth,
	}
	s.setupRouter()

	return s
}

func (s *Server) AddCloser(closer CloseFunc) {
	s.closers = append(s.closers, closer)
}

func (s *Server) setupRouter() {
	s.e = echo.New()
	s.e.HideBanner = true
	s.e.Validator = NewValidator()

	s.e.Pre(middleware.RemoveTrailingSlash())
	s.e.Use(middleware.RequestID())

	s.e.GET("/auth/oauth/github/link", handlers.HandleGetGitHubAuthLink(s.auth))
	s.e.GET("/auth/oauth/github/callback", handlers.HandleGitHubAuthCallback(s.auth))
	s.e.GET("/auth/token.html", handlers.HandleTokenPage())
	s.e.GET("/auth/refresh", handlers.HandleRefresh())
	s.e.GET("/static/*", handlers.HandleStatic())

	restricted := s.e.Group("/api")
	{
		restricted.Use(echojwt.WithConfig(makeJWTConfig()))
		restricted.POST("/short", handlers.HandleShorten(s.shortener))
		restricted.GET("/stats/:identifier", handlers.HandleStats(s.shortener))
	}

	s.e.GET("/:identifier", handlers.HandleRedirect(s.shortener))

	s.AddCloser(s.e.Shutdown)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.e.ServeHTTP(w, r)
}

func (s *Server) Shutdown(ctx context.Context) error {
	for _, fn := range s.closers {
		if err := fn(ctx); err != nil {
			return err
		}
	}

	return nil
}

func makeJWTConfig() echojwt.Config {
	return echojwt.Config{
		SigningKey: []byte(config.Get().Auth.JWTSecretKey),
		ErrorHandler: func(c echo.Context, err error) error {
			return echo.NewHTTPError(http.StatusUnauthorized)
		},
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(model.UserClaims)
		},
	}
}
