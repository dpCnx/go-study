package main

import (
	"context"
	"fmt"
	proto "github.com/dpCnx/go-study/demo/micro/wrapper/proto"
	"github.com/micro/go-micro/util/log"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/server"
)

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloResponse) error {
	log.Logf("[hello] 收到请求")
	rsp.Greeting = "你好呀！" + req.Name
	return nil
}

// logWrapper1 包装HandlerFunc类型的接口
func logW1(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		log.Logf("[logWrapper1] %s 收到请求", req.Endpoint())
		err := fn(ctx, req, rsp)
		return err
	}
}

// logWrapper2 包装HandlerFunc类型的接口
func logW2(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		log.Logf("[logWrapper2] %s 收到请求", req.Endpoint())
		err := fn(ctx, req, rsp)
		return err
	}
}

func main() {
	service := micro.NewService(
		micro.Name("greeter"),
		// 声明包装器
		micro.WrapHandler(logW1, logW2),
	)

	service.Init()

	proto.RegisterGreeterHandler(service.Server(), new(Greeter))

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
