package main

import (
	"context"
	"fmt"
	"log"
	"test_grpc_go/proto/gen/service/user/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ls, err := grpc.Dial("localhost:6000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	cl := user.NewUserServiceClient(ls)
	resp, err := cl.SayHello(context.Background(), &user.HelloRequest{Name: "test"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
}
