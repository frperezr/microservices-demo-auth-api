package service

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/frperezr/noken-test/lib/clients/users"
	"github.com/frperezr/noken-test/src/auth-api"
	"golang.org/x/crypto/bcrypt"

	user "github.com/frperezr/noken-test/src/users-api"
)

// AuthSvc ...
type AuthSvc struct {
	Client *users.Client
}

// New ...
func New() *AuthSvc {
	usersHost := os.Getenv("USERS_HOST")
	if usersHost == "" {
		log.Fatal("[Auth Service][New Service][Error] missing USERS_HOST env variable")
	}

	usersPort := os.Getenv("USERS_PORT")
	if usersPort == "" {
		log.Fatal("[Auth Service][New Service][Error] missing USERS_PORT env variable")
	}

	return &AuthSvc{
		Client: &users.Client{
			UsersHost: usersHost,
			UsersPort: usersPort,
		},
	}
}

// Login ...
func (as *AuthSvc) Login(email, password string) (string, error) {
	user, err := as.Client.GetByEmail(email)
	if err != nil {
		return "", errors.New(err.Error())
	}

	expiresAt := time.Now().Add(time.Minute * 100000).Unix()

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", errors.New("Invalid login credentials. Please try again")
	}

	tk := &auth.Token{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}

	return tokenString, nil
}

// Signup handles the creation of a new user
func (as *AuthSvc) Signup(user *user.User) (string, error) {
	if user == nil {
		return "", errors.New("missing user param")
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New(err.Error())
	}

	user.Password = string(pass)

	usr, err := as.Client.Create(user)
	if err != nil {
		return "", errors.New(err.Error())
	}

	token, err := as.Login(usr.Email, usr.Password)
	if err != nil {
		return "", errors.New(err.Error())
	}

	return token, nil
}
