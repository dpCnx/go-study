package main

import (
	"github.com/micro/go-micro/v2"
)

func main() {
	service := micro.NewService(
		micro.Name("one.server"),
	)
	service.Init()
	//micro.RegisterHandler(service)
}
