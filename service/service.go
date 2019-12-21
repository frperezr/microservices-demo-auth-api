package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/frperezr/noken-test/src/auth-api"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"

	pb "github.com/frperezr/noken-test/pb"

	user "github.com/frperezr/noken-test/src/users-api"
)

// AuthSvc ...
type AuthSvc struct {
	Client pb.UserServiceClient
	Conn   *grpc.ClientConn
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

	addr := fmt.Sprintf("%v:%v", usersHost, usersPort)

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(fmt.Sprintf("[Auth Service][New Service][Error] %v", err))
	}

	return &AuthSvc{
		Client: pb.NewUserServiceClient(conn),
		Conn:   conn,
	}
}

// Login ...
func (as *AuthSvc) Login(email, password string) (string, error) {
	res, err := as.Client.GetByEmail(context.Background(), &pb.GetUserByEmailRequest{
		Email: email,
	})
	if err != nil {
		return "", err
	}

	if res.GetError() != nil {
		return "", errors.New(res.GetError().GetMessage())
	}

	user := res.GetData()

	expiresAt := time.Now().Add(time.Minute * 100000).Unix()

	err = bcrypt.CompareHashAndPassword([]byte(user.GetPassword()), []byte(password))
	if err != nil {
		return "", errors.New("Invalid login credentials. Please try again")
	}

	tk := &auth.Token{
		UserID: user.Id,
		Name:   user.Name,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Signup handles the creation of a new user
func (as *AuthSvc) Signup(u *user.User) (string, error) {
	if u == nil {
		return "", errors.New("missing user param")
	}

	pwd := u.Password

	pass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	u.Password = string(pass)
	u.Email = strings.ToLower(u.Email)

	res, err := as.Client.Create(context.Background(), &pb.CreateUserRequest{
		Data: u.ToProto(),
	})
	if err != nil {
		return "", err
	}

	if res.GetError() != nil {
		return "", errors.New(res.GetError().GetMessage())
	}

	token, err := as.Login(u.Email, pwd)
	if err != nil {
		return "", err
	}

	return token, nil
}
