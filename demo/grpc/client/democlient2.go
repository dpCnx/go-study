package main

import (
	"context"
	dp "github.com/dpCnx/go-study/demo/grpc/grpcproto"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:9002", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}

	defer conn.Close()

	client := dp.NewStreamServiceClient(conn)

	//demo1(err, client)

	demo2(err, client)

	//demo3(err, client)
}

func demo3(err error, client dp.StreamServiceClient) {
	steamcs, err := client.Route(context.Background())
	if err != nil {
		log.Fatalf("client Route err:%v", err)
	}
	for n := 0; n < 6; n++ {
		err = steamcs.Send(&dp.StreamRequest{
			Pt: &dp.StreamPoint{
				Name:  "client",
				Value: 0,
			},
		})
		if err != nil {
			log.Printf("Send err:%v", err)
		}

		resp, err := steamcs.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Recv err:%v", err)
		}

		log.Printf("resp: pj.name: %s, pt.value: %d", resp.Pt.Name, resp.Pt.Value)
	}
	steamcs.CloseSend()
}

func demo2(err error, client dp.StreamServiceClient) {
	stramclient, err := client.Record(context.Background())
	if err != nil {
		log.Fatalf("client Record err:%v", err)
	}
	for n := 0; n < 6; n++ {
		err := stramclient.Send(&dp.StreamRequest{
			Pt: &dp.StreamPoint{
				Name:  "come",
				Value: 1,
			},
		})
		if err != nil {
			log.Printf("Send err:%v", err)
		}
	}
	resp, err := stramclient.CloseAndRecv()
	if err != nil {
		log.Printf("Close err:%v", err)
	}
	log.Printf("resp: pj.name: %s, pt.value: %d", resp.Pt.Name, resp.Pt.Value)
}

func demo1(err error, client dp.StreamServiceClient) {
	streamclient, err := client.List(context.Background(), &dp.StreamRequest{
		Pt: &dp.StreamPoint{
			Name:  "list demo",
			Value: 2019,
		},
	})
	if err != nil {
		log.Fatalf("client list err:%v", err)
	}
	for {
		resp, err := streamclient.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("recv err:%v", err)
			return
		}

		log.Printf("resp: pj.name: %s, pt.value: %d", resp.Pt.Name, resp.Pt.Value)
	}
}
