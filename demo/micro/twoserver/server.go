package main

import (
	"context"
	twoserver_proto "github.com/dpCnx/go-study/demo/micro/twoserver/proto"
	"github.com/micro/go-micro/v2"
)

type ServiceTwo struct {
}

func (s *ServiceTwo) SayHello(ctx context.Context, req *twoserver_proto.HelloRequest, repose *twoserver_proto.HelloReplay) error {

	repose.Message = "你好呀:" + req.Name

	return nil
}

func main() {

	service := micro.NewService(
		micro.Name("service.two"),
		micro.Version("latest"),
	)

	service.Init()

	twoserver_proto.RegisterHelloServerHandler(service.Server(), new(ServiceTwo))

	service.Run()
}

//protoc --proto_path=. --go_out=. --micro_out=. proto/example/example.proto
