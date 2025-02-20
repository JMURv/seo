package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb "github.com/JMURv/protos/par-pro"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func sendAuthRequest(ctx context.Context, email, password string) string {
	address := "localhost:50050"

	cli, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to create client: %v", err)
		return ""
	}
	defer func(cli *grpc.ClientConn) {
		if err := cli.Close(); err != nil {
			log.Println("failed to close client")
		}
	}(cli)

	res, err := pb.NewSSOClient(cli).Authenticate(
		ctx, &pb.SSO_EmailAndPasswordRequest{
			Email:    email,
			Password: password,
		},
	)
	if err != nil {
		log.Printf("authentication error: %v", err)
		return ""
	}

	return res.Token
}

func main() {
	email := flag.String("email", "", "User email")
	password := flag.String("password", "", "User password")
	flag.Parse()

	if *email == "" || *password == "" {
		log.Fatal("email and password must be provided")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	token := sendAuthRequest(ctx, *email, *password)
	if token == "" {
		log.Fatal("failed to obtain token")
	}

	fmt.Println(token)
}
