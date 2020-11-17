package main

import (
	"context"
	pb "github.com/dpCnx/go-study/demo/grpc/grpcproto"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
	"runtime/debug"
	"time"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReplay, error) {

	data := make(chan *pb.HelloReplay)
	go handler(in, data)
	select {
	case res := <-data:
		return res, nil
	case <-ctx.Done():
		log.Printf("ctx done err:%s \n", ctx.Err())
		return nil, status.Errorf(codes.Canceled, "time out")
	}
}

func handler(request *pb.HelloRequest, replays chan *pb.HelloReplay) {

	time.Sleep(7 * time.Second)

	res := pb.HelloReplay{
		Message: "HELLO--->" + request.Name,
	}

	replays <- &res

}

func (s *server) GetHelloMsg(ctx context.Context, in *pb.HelloRequest) (*pb.HelloMessage, error) {

	return &pb.HelloMessage{Msg: "hello" + in.Name}, nil
}

func demoOne() {
	ln, err := net.Listen("tcp", "127.0.0.1:9999")
	if err != nil {
		log.Printf("net listen err: %v \n", err)
		return
	}
	s := grpc.NewServer()
	pb.RegisterHelloServerServer(s, &server{})

	if err = s.Serve(ln); err != nil {
		log.Printf("s serve err: %v \n", err)
		return
	}
}

type StreamService struct{}

func (s *StreamService) List(r *pb.StreamRequest, stream pb.StreamService_ListServer) error {

	for n := 0; n < 6; n++ {
		if err := stream.Send(&pb.StreamResponse{
			Pt: &pb.StreamPoint{
				Name:  r.Pt.Name,
				Value: r.Pt.Value + int32(n),
			},
		}); err != nil {
			return err
		}
	}

	return nil
}

func (s *StreamService) Record(stream pb.StreamService_RecordServer) error {

	for {
		r, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.StreamResponse{Pt: &pb.StreamPoint{Name: "Record", Value: 1}})
		}
		if err != nil {
			return err
		}

		log.Printf("stream.Recv pt.name: %s, pt.value: %d \n", r.Pt.Name, r.Pt.Value)
	}

}

func (s *StreamService) Route(stream pb.StreamService_RouteServer) error {

	n := 0

	for {
		err := stream.Send(&pb.StreamResponse{
			Pt: &pb.StreamPoint{
				Name:  "server",
				Value: int32(n),
			},
		})
		if err != nil {
			return err
		}

		r, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		n++

		log.Printf("stream.Recv pt.name: %s, pt.value: %d \n", r.Pt.Name, r.Pt.Value)
	}

	return nil

}

func demoTwo() {
	ln, err := net.Listen("tcp", ":9002")
	if err != nil {
		log.Printf("net listen err: %v \n", err)
		return
	}
	server := grpc.NewServer()
	pb.RegisterStreamServiceServer(server, &StreamService{})
	if err = server.Serve(ln); err != nil {
		log.Printf("s serve err: %v \n", err)
		return
	}
}

type serverTls struct{}

func (s *serverTls) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReplay, error) {

	return &pb.HelloReplay{Message: "Hello " + in.Name}, nil
}

func (s *serverTls) GetHelloMsg(ctx context.Context, in *pb.HelloRequest) (*pb.HelloMessage, error) {

	return &pb.HelloMessage{Msg: "hello" + in.Name}, nil
}

func demothree() {
	ln, err := net.Listen("tcp", "127.0.0.1:9999")
	if err != nil {
		log.Printf("net listen err: %v \n", err)
		return
	}
	c, err := credentials.NewServerTLSFromFile("demo/grpc/tls/cert.pem", "demo/grpc/tls/cert.key")
	if err != nil {
		log.Printf("credentials newservertlsfromfile err: %v \n", err)
		return
	}
	s := grpc.NewServer(grpc.Creds(c))
	pb.RegisterHelloServerServer(s, &serverTls{})
	if err = s.Serve(ln); err != nil {
		log.Printf("s serve err: %v \n", err)
		return
	}
}

type serverIntercepte struct{}

func (s *serverIntercepte) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReplay, error) {
	return &pb.HelloReplay{Message: "Hello " + in.Name}, nil
}

func (s *serverIntercepte) GetHelloMsg(ctx context.Context, in *pb.HelloRequest) (*pb.HelloMessage, error) {

	return &pb.HelloMessage{Msg: "hello" + in.Name}, nil
}

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	m, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "token is nil")
	}

	if val, ok := m["token"]; ok {
		log.Println("t:", val[0])
	}

	resp, err := handler(ctx, req)

	return resp, err
}

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	start := time.Now()
	log.Printf("server:%v  method:%s request:%v \n", info.Server, info.FullMethod, req)

	resp, err := handler(ctx, req)

	cost := time.Since(start)
	log.Printf("time:%d server:%v  method:%s request:%v respose:%v \n", cost, info.Server, info.FullMethod, req, resp)

	return resp, err
}

func RecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("error:%v \n", e)
			log.Printf("stack:%s \n", string(debug.Stack()))
			err = status.Errorf(codes.Internal, "Panic err: %v", e)
		}
	}()

	return handler(ctx, req)

}

func demofour() {
	lis, err := net.Listen("tcp", "127.0.0.1:9999")
	if err != nil {
		log.Printf("net listen err: %v \n", err)
		return
	}
	c, err := credentials.NewServerTLSFromFile("demo/grpc/tls/cert.pem", "demo/grpc/tls/cert.key")
	if err != nil {
		log.Printf("credentials newservertlsfromfile err: %v \n", err)
		return
	}
	opts := []grpc.ServerOption{
		grpc.Creds(c),
		grpcMiddleware.WithUnaryServerChain(
			AuthInterceptor,
			LoggingInterceptor,
			RecoveryInterceptor,
		),
	}
	s := grpc.NewServer(opts...)
	pb.RegisterHelloServerServer(s, &serverIntercepte{})
	if err = s.Serve(lis); err != nil {
		log.Printf("s server err: %v \n", err)
		return
	}
}

func main() {

	//demoOne()

	//demoTwo()

	//demothree()

	//demofour()
}

/*
	protoc --go_out=plugins=grpc:./ *.proto #添加grpc插件
	openssl genrsa > cert.key //生成私钥
	openssl req -new -x509 -sha256 -key cert.key -out cert.pem -days 3650 //生成证书
*/
