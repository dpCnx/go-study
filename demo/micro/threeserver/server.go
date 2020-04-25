package main

import (
	"context"
	threeserver_proto "github.com/dpCnx/go-study/demo/micro/threeserver/proto"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/errors"
	"io/ioutil"
	"log"
)

type ServerThree struct {
}

func (s *ServerThree) GetFile(ctx context.Context, file threeserver_proto.FileServer_GetFileStream) error {

	temp, err := ioutil.TempFile("", "demo")
	if err != nil {
		return errors.InternalServerError("service.three", err.Error())
	}
	for {
		filerequest, err := file.Recv()
		if err != nil {
			return errors.InternalServerError("service.three", err.Error())
		}
		if filerequest.Len == -1 {
			break
		}
		if _, err := temp.Write(filerequest.Byte); err != nil {
			return errors.InternalServerError("service.three", err.Error())
		}
	}

	log.Println(temp.Name())

	return file.SendMsg(&threeserver_proto.FileStreamMsg{
		Filename: temp.Name(),
	})
}

func main() {
	service := micro.NewService(
		micro.Name("service.three"),
		micro.Version("latest"),
	)

	service.Init()
	threeserver_proto.RegisterFileServerHandler(service.Server(), new(ServerThree))
	service.Run()
}
