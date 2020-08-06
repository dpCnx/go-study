package main

import (
	"context"
	"fmt"
	esvc "github.com/dpCnx/go-study/demo/grpc-etcd/etcd"
	pb "github.com/dpCnx/go-study/demo/grpc-etcd/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

func main() {

	r := esvc.NewResolver("127.0.0.1:2379")
	resolver.Register(r)

	// 客户端连接服务器
	conn, err := grpc.Dial(r.Scheme()+"://author/"+"serverone", grpc.WithBalancerName("round_robin"), grpc.WithInsecure())

	if err != nil {
		fmt.Println("连接服务器失败", err)
	}
	defer conn.Close()

	// 获得grpc句柄
	c := pb.NewHelloServerClient(conn)

	// 远程单调用 SayHi 接口
	r1, err := c.SayHi(
		context.Background(),
		&pb.HelloRequest{
			Name: "Kitty",
		},
	)
	if err != nil {
		fmt.Println("Can not get SayHi:", err)
		return
	}
	fmt.Printf("SayHi 响应：%s\n", r1.GetMessage())

	// 远程单调用 GetMsg 接口
	r2, err := c.GetMsg(
		context.Background(),
		&pb.HelloRequest{
			Name: "Kitty",
		},
	)
	if err != nil {
		fmt.Println("Can not get GetMsg:", err)
		return
	}
	fmt.Printf("GetMsg 响应：%s\n",r2.GetMsg())

}
