package main

import (
	"context"
	"fmt"
	esvc "github.com/dpCnx/go-study/demo/grpc-etcd/etcd"
	pb "github.com/dpCnx/go-study/demo/grpc-etcd/proto"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const servername = "serverone"

type server struct{}

func (s *server) SayHi(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReplay, error) {
	fmt.Printf("SayHi: %s, %s\n", in.Name, time.Now().Format("2006-01-02 15:04:05"))
	return &pb.HelloReplay{
		Message: "127.0.0.1:8888",
	}, nil
}

func (s *server) GetMsg(ctx context.Context, in *pb.HelloRequest) (*pb.HelloMessage, error) {
	fmt.Printf("GetMsg: %s, %s\n", in.Name, time.Now().Format("2006-01-02 15:04:05"))
	return &pb.HelloMessage{
		Msg: "127.0.0.1:8888",
	}, nil
}

func main() {

	ln, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("网络异常：", err)
		return
	}
	defer ln.Close()

	srv := grpc.NewServer()
	defer srv.GracefulStop()

	pb.RegisterHelloServerServer(srv, &server{})

	addr := fmt.Sprintf("%s:%s", "127.0.0.1", "8888")
	fmt.Printf("server addr:%s\n", addr)
	_ = esvc.Register(servername, addr, 15)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		esvc.UnRegister(servername, addr)

		if i, ok := s.(syscall.Signal); ok {
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}

	}()

	err = srv.Serve(ln)
	if err != nil {
		fmt.Println("监听异常：", err)
		return
	}

}
