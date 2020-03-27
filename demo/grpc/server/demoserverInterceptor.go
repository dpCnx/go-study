package main

import (
	"context"
	"fmt"
	pb "github.com/dpCnx/go-study/demo/grpc/grpcproto"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"runtime/debug"
)
/*
	添加拦截器
*/
type serverIntercepte struct{}

func (s *serverIntercepte) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReplay, error) {

	return &pb.HelloReplay{Message: "Hello " + in.Name}, nil
}

func (s *serverIntercepte) GetHelloMsg(ctx context.Context, in *pb.HelloRequest) (*pb.HelloMessage, error) {

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

	opts := []grpc.ServerOption{
		grpc.Creds(c),
		grpc_middleware.WithUnaryServerChain(
			RecoveryInterceptor,
			LoggingInterceptor,
		),
	}

	s := grpc.NewServer(opts...)
	pb.RegisterHelloServerServer(s, &serverIntercepte{})
	reflection.Register(s)
	if err = s.Serve(lis); err != nil {
		fmt.Println("server err:", err)
		return
	}

}

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("gRPC method: %s, %v", info.FullMethod, req)
	resp, err := handler(ctx, req)
	log.Printf("gRPC method: %s, %v", info.FullMethod, resp)
	return resp, err
}

func RecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "Panic err: %v", e)
		}
	}()

	return handler(ctx, req)
}
