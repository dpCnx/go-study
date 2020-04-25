package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/v2"
)

type ServiceOne struct {
}

func (s *ServiceOne) Hw(ctx context.Context, req *string, rsp *string) error {

	*rsp = fmt.Sprintf("hello:%s", *req)

	return nil
}

func main() {
	service := micro.NewService(
		micro.Name("service.one"),
	)
	service.Init()
	micro.RegisterHandler(service.Server(), new(ServiceOne))
	service.Run()
}
