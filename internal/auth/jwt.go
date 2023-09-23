package auth

import (
	"shorter-url/internal/config"
	"shorter-url/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func MakeJWT(user model.User, duration time.Duration) (string, error) {
	claims := model.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
		User: user,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Get().Auth.JWTSecretKey))
}
