package auth

import (
	user "github.com/frperezr/microservices-demo/src/users-api"

	"github.com/dgrijalva/jwt-go"
)

// Token is a jwt returned by auth service
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
