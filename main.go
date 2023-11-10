package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"test_grpc_go/proto/gen/service/task/v1"
	"test_grpc_go/proto/gen/service/user/v1"
)

type TaskService struct {
	task.UnimplementedTaskServiceServer
}

func (svc *TaskService) SayHello(ctx context.Context, req *task.HelloRequest) (*task.HelloResponse, error) {

	return &task.HelloResponse{Msg: "service task say hello"}, nil
}

type UserService struct {
	user.UnimplementedUserServiceServer
}

func (svc *UserService) SayHello(ctx context.Context, req *user.HelloRequest) (*user.HelloResponse, error) {

	return &user.HelloResponse{Msg: "service user say hello"}, nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	lis, err := net.Listen("tcp", ":6000")
	if err != nil {
		log.Fatal(err)
	}

	grpcSvr := grpc.NewServer()
	task.RegisterTaskServiceServer(grpcSvr, &TaskService{})
	user.RegisterUserServiceServer(grpcSvr, &UserService{})

	go func() {
		if err := grpcSvr.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err = task.RegisterTaskServiceHandlerFromEndpoint(ctx, mux, "localhost:6000", opts)
	if err != nil {
		log.Fatal(err)
	}

	err = user.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "localhost:6000", opts)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	<-sigChan
	// grpcSvr.GracefulStop()
	fmt.Println("graceful stop")
}
