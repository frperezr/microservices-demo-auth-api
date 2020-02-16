package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	pb "github.com/frperezr/microservices-demo/pb"
	"github.com/frperezr/microservices-demo/src/users-api"
	"google.golang.org/grpc"
)

func main() {
	flag.Parse()

	usersHost := os.Getenv("USERS_HOST")
	if usersHost == "" {
		fmt.Print(`{"error": "missing env USERS_HOST"}`)
		os.Exit(1)
	}

	usersPort := os.Getenv("USERS_PORT")
	if usersPort == "" {
		fmt.Print(`{"error": "missing env USERS_PORT"}`)
		os.Exit(1)
	}

	addr := fmt.Sprintf("%v:%v", usersHost, usersPort)

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	c := pb.NewAuthServiceClient(conn)

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	subcmd := flag.Arg(0)
	var result string

	switch subcmd {
	case "login":
		result, err = Login(c, flag.Args()[1:])
		if err != nil {
			fmt.Print(fmt.Sprintf(`{"error": "%v"}`, err.Error()))
		}
	case "signup":
		result, err = Signup(c, flag.Args()[1:])
		if err != nil {
			fmt.Print(fmt.Sprintf(`{"error": "%v"}`, err.Error()))
		}
	default:
		fmt.Print(`{"error": "invalid command"}`)
		os.Exit(1)
	}

	fmt.Print(result)
	os.Exit(0)
}

// Login ...
func Login(as pb.AuthServiceClient, args []string) (string, error) {
	if len(args) != 1 {
		flag.Usage()
		return "", errors.New("missing params")
	}

	jsonStr := args[0]
	data := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return "", errors.New("invalid JSON")
	}

	res, err := as.Login(context.Background(), &pb.LoginRequest{
		Email:    data.Email,
		Password: data.Password,
	})

	if err != nil {
		return "", err
	}

	if res.Error != nil {
		return "", err
	}

	json, err := json.Marshal(res.GetToken())
	if err != nil {
		return "", errors.New("cant marshal data")
	}

	return string(json), nil
}

// Signup ....
func Signup(as pb.AuthServiceClient, args []string) (string, error) {
	if len(args) != 1 {
		flag.Usage()
		return "", errors.New("missing params")
	}

	jsonStr := args[0]
	data := struct {
		User *users.User `json:"user"`
	}{}

	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return "", errors.New("invalid JSON")
	}

	res, err := as.Signup(context.Background(), &pb.SignupRequest{
		Data: data.User.ToProto(),
	})

	if err != nil {
		return "", err
	}

	if res.Error != nil {
		return "", fmt.Errorf(res.Error.GetMessage())
	}

	json, err := json.Marshal(res.GetToken())
	if err != nil {
		return "", errors.New("cant marshal data")
	}

	return string(json), nil
}
