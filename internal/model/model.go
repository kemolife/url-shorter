package model

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	. "github.com/samber/mo"
	"time"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrIdentifierExists = errors.New("identifier already exists")
	ErrInvalidURL       = errors.New("invalid url")
	ErrUserIsNotMember  = errors.New("user is not member of the organization")
	ErrInvalidToken     = errors.New("invalid token")
)

type ShortenInput struct {
	RawURL     string
	Identifier Option[string]
	CreatedBy  string
}

type Shortening struct {
	Identifier  string    `json:"identifier"`
	CreatedBy   string    `json:"created_by"`
	OriginalURL string    `json:"original_url"`
	Visits      int64     `json:"visits"`
	Alias       string    `json:"alias"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type User struct {
	IsActive    bool   `json:"is_verified,omitempty"`
	GitHubLogin string `json:"gh_login"`

	// TODO: should we store it in something like Vault?
	GitHubAccessKey string    `json:"gh_access_key,omitempty"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
}

type UserClaims struct {
	jwt.RegisteredClaims
	User `json:"user_data"`
}
