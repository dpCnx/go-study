package main

import (
	"context"
	threeserver_proto "github.com/dpCnx/go-study/demo/micro/threeserver/proto"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"io"
	"log"
	"net/http"
)

var fileservice threeserver_proto.FileServerService
var c client.Client

func main() {
	service := micro.NewService()
	service.Init()
	c = service.Client()
	fileservice = threeserver_proto.NewFileServerService("service.three", c)

	http.HandleFunc("/", uploadfile)
	if err := http.ListenAndServe(":9966", nil); err != nil {
		log.Println(err)
	}
}

func uploadfile(respose http.ResponseWriter, resquest *http.Request) {
	resquest.ParseMultipartForm(10 << 20)
	files, _ := resquest.MultipartForm.File["file"]
	file, _ := files[0].Open()
	next, _ := c.Options().Selector.Select("service.three")
	node, _ := next()
	stream, _ := fileservice.GetFile(context.Background(), func(options *client.CallOptions) {
		options.Address = []string{node.Address}
	})

	for {
		buff := make([]byte, 1024*1024)
		n, err := file.Read(buff)
		if err != nil && err == io.EOF {
			err = stream.Send(&threeserver_proto.FileRequest{
				Byte: nil,
				Len:  -1,
			})
			break
		}
		_ = stream.Send(&threeserver_proto.FileRequest{
			Byte: buff[:n],
			Len:  int64(n),
		})

	}

	var res threeserver_proto.FileStreamMsg
	_ = stream.RecvMsg(&res)
	log.Println(res.Filename)
}
