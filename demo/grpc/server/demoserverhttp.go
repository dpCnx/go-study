package main

import (
	"context"
	pb "github.com/dpCnx/go-study/demo/grpc/grpcproto"
	"google.golang.org/grpc"
	"net/http"
	"strings"
)

type serverhttp struct{}

func (s *serverhttp) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReplay, error) {

	return &pb.HelloReplay{Message: "Hello " + in.Name}, nil
}

func (s *serverhttp) GetHelloMsg(ctx context.Context, in *pb.HelloRequest) (*pb.HelloMessage, error) {

	return &pb.HelloMessage{Msg: "hello" + in.Name}, nil
}

func main() {

	s = grpc.NewServer()
	pb.RegisterHelloServerServer(s, &serverhttp{})

	mux = http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("http:grpc"))
	})

	engin := new(Engine)
	http.ListenAndServe("127.0.0.1:9999", engin)
}

var mux *http.ServeMux
var s *grpc.Server

type Engine struct{}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
		s.ServeHTTP(w, r)
	} else {
		mux.ServeHTTP(w, r)
	}

}
