package main

import (
	"context"
	"fmt"
	pt "github.com/dpCnx/go-study/demo/grpc/grpcproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

func main() {
	//生成证书时：Common Name (eg, fully qualified host name) []:d
	ctls, err := credentials.NewClientTLSFromFile("demo/grpc/tls/cert.pem", "d")
	if err != nil {
		log.Fatalf("credentials.NewClientTLSFromFile err: %v", err)
	}

	conn, err := grpc.Dial("127.0.0.1:9999", grpc.WithTransportCredentials(ctls))

	if err != nil {
		fmt.Println("连接服务器失败", err)
	}
	defer conn.Close()

	c := pt.NewHelloServerClient(conn)

	r1, err := c.SayHello(context.Background(), &pt.HelloRequest{Name: "panda"})

	if err != nil {
		fmt.Println("cloud not get Hello server ..", err)
		return
	}

	fmt.Println("HelloServer resp: ", r1.Message)

	r2, err := c.GetHelloMsg(context.Background(), &pt.HelloRequest{Name: "panda2"})

	if err != nil {
		fmt.Println("cloud not get hello msg ..", err)
		return

	}

	fmt.Println("HelloServer resp: ", r2.Msg)

}
