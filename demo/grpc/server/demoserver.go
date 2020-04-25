package main

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/dpCnx/go-study/demo/grpc/grpcproto"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReplay, error) {
	/*
		由于context超时取消
	*/
	if ctx.Err() == context.Canceled {
		return nil, errors.New("SearchService.Search canceled")
	}

	return &pb.HelloReplay{Message: "Hello " + in.Name}, nil
}

func (s *server) GetHelloMsg(ctx context.Context, in *pb.HelloRequest) (*pb.HelloMessage, error) {

	return &pb.HelloMessage{Msg: "hello" + in.Name}, nil
}

func main() {

	lis, err := net.Listen("tcp", "127.0.0.1:9999")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}

	s := grpc.NewServer()
	pb.RegisterHelloServerServer(s, &server{})
	if err = s.Serve(lis); err != nil {
		fmt.Println("server err:", err)
		return
	}

}

//protoc --go_out=plugins=grpc:./ *.proto #添加grpc插件
