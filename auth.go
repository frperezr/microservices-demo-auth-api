package auth

import (
	user "github.com/frperezr/noken-test/src/users-api"

	"github.com/dgrijalva/jwt-go"
)

// Token is a jwt returned by noken auth service
type Token struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	*jwt.StandardClaims
}

// Service ...
type Service interface {
	Login(email, password string) (string, error)
	Signup(user *user.User) (string, error)
}
