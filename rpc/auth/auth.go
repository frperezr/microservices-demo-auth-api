package auth

import (
	"context"

	"github.com/frperezr/microservices-demo/src/users-api"

	pb "github.com/frperezr/microservices-demo/pb"
	"github.com/frperezr/microservices-demo/src/auth-api"
	"github.com/frperezr/microservices-demo/src/auth-api/service"
)

var _ pb.AuthServiceServer = (*Service)(nil)

// Service ...
type Service struct {
	AuthSvc auth.Service
}

// New ...
func New() *Service {
	return &Service{
		AuthSvc: service.New(),
	}
}

// Login ...
func (s *Service) Login(ctx context.Context, gr *pb.LoginRequest) (*pb.LoginResponse, error) {
	email := gr.GetEmail()
	pwd := gr.GetPassword()

	token, err := s.AuthSvc.Login(email, pwd)
	if err != nil {
		if err.Error() == "Invalid login credentials. Please try again" {
			return &pb.LoginResponse{
				Token: "",
				Error: &pb.Error{
					Code:    401,
					Message: err.Error(),
				},
			}, nil
		}

		return &pb.LoginResponse{
			Token: "",
			Error: &pb.Error{
				Code:    500,
				Message: err.Error(),
			},
		}, nil
	}

	return &pb.LoginResponse{
		Token: token,
		Error: nil,
	}, nil
}

// Signup ...
func (s *Service) Signup(ctx context.Context, gr *pb.SignupRequest) (*pb.SignupResponse, error) {
	data := gr.GetData()

	user := &users.User{
		Email:    data.GetEmail(),
		Name:     data.GetName(),
		LastName: data.GetName(),
		Password: data.GetPassword(),
	}

	token, err := s.AuthSvc.Signup(user)
	if err != nil {
		if err.Error() == "user already registered" {
			return &pb.SignupResponse{
				Token: "",
				Error: &pb.Error{
					Code:    400,
					Message: "user already registered",
				},
			}, nil
		}
		return &pb.SignupResponse{
			Token: "",
			Error: &pb.Error{
				Code:    500,
				Message: err.Error(),
			},
		}, nil
	}

	return &pb.SignupResponse{
		Token: token,
		Error: nil,
	}, nil
}
