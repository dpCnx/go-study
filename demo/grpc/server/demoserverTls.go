package main

import (
	"context"
	"fmt"
	pb "github.com/dpCnx/go-study/demo/grpc/grpcproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type serverTls struct{}

func (s *serverTls) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReplay, error) {

	return &pb.HelloReplay{Message: "Hello " + in.Name}, nil
}

func (s *serverTls) GetHelloMsg(ctx context.Context, in *pb.HelloRequest) (*pb.HelloMessage, error) {

	return &pb.HelloMessage{Msg: "hello" + in.Name}, nil
}

func main() {

	lis, err := net.Listen("tcp", "127.0.0.1:9999")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}

	c, err := credentials.NewServerTLSFromFile("demo/grpc/tls/cert.pem", "demo/grpc/tls/cert.key")
	if err != nil {
		log.Fatalf("credentials.NewServerTLSFromFile err: %v", err)
	}

	s := grpc.NewServer(grpc.Creds(c))
	pb.RegisterHelloServerServer(s, &serverTls{})
	reflection.Register(s)
	if err = s.Serve(lis); err != nil {
		fmt.Println("server err:", err)
		return
	}

}
