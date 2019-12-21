package main

import (
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/frperezr/noken-test/pb"

	authSvc "github.com/frperezr/noken-test/src/auth-api/rpc/auth"

	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Env variable PORT must be defined")
	}

	srv := grpc.NewServer()
	svc := authSvc.New()

	pb.RegisterAuthServiceServer(srv, svc)
	reflection.Register(srv)

	log.Println("Starting auth service...")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Printf("Failed to list: %v", err)
		os.Exit(1)
	}

	log.Println(fmt.Sprintf("Auth service running, Listening on: %v", port))

	if err := srv.Serve(lis); err != nil {
		log.Printf("Fatal to serve: %v", err)
		os.Exit(1)
	}
}
