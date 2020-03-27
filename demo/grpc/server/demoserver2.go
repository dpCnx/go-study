package main

import (
	"fmt"
	dp "github.com/dpCnx/go-study/demo/grpc/grpcproto"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)
/*
	流式stream
*/
type StreamService struct{}

const (
	PORT = "9002"
)

func main() {

	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	server := grpc.NewServer()
	dp.RegisterStreamServiceServer(server, &StreamService{})

	if err = server.Serve(lis); err != nil {
		fmt.Println("server err:", err)
		return
	}

}

func (s *StreamService) List(r *dp.StreamRequest, stream dp.StreamService_ListServer) error {

	for n := 0; n < 6; n++ {
		err := stream.Send(&dp.StreamResponse{
			Pt: &dp.StreamPoint{
				Name:  r.Pt.Name,
				Value: r.Pt.Value + int32(n),
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *StreamService) Record(stream dp.StreamService_RecordServer) error {
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&dp.StreamResponse{Pt: &dp.StreamPoint{Name: "Record", Value: 1}})
		}
		if err != nil {
			return err
		}

		log.Printf("stream.Recv pt.name: %s, pt.value: %d", r.Pt.Name, r.Pt.Value)
	}
}

func (s *StreamService) Route(stream dp.StreamService_RouteServer) error {
	n := 0
	for {
		err := stream.Send(&dp.StreamResponse{
			Pt: &dp.StreamPoint{
				Name:  "server",
				Value: int32(n),
			},
		})
		if err != nil {
			return err
		}

		r, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		n++

		log.Printf("stream.Recv pt.name: %s, pt.value: %d", r.Pt.Name, r.Pt.Value)
	}

	return nil

}
