package main

import (
	"context"
	"fmt"
	pt "github.com/dpCnx/go-study/demo/grpc/grpcproto"
	"google.golang.org/grpc"
	"time"
)

func main() {

	conn, err := grpc.Dial("127.0.0.1:9999", grpc.WithInsecure())

	if err != nil {
		fmt.Println("连接服务器失败", err)
	}
	defer conn.Close()

	c := pt.NewHelloServerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r1, err := c.SayHello(ctx, &pt.HelloRequest{Name: "panda"})

	if err != nil {
		fmt.Println("cloud not get Hello server ..", err)
		return
	}

	fmt.Println("HelloServer resp: ", r1.Message)

	r2, err := c.GetHelloMsg(context.Background(), &pt.HelloRequest{Name: "panda"})

	if err != nil {
		fmt.Println("cloud not get hello msg ..", err)
		return

	}

	fmt.Println("HelloServer resp: ", r2.Msg)

}
