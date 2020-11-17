package main

import (
	"context"
	pt "github.com/dpCnx/go-study/demo/grpc/grpcproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"log"
	"time"
)

func main() {

	demoOne()

	//demoTwo()

	//demoThree()

	//demoFour()

	//demoFive()

	//demoSix()
}

// 实现grpc.PerRPCCredentials接⼝
type AuthToken struct {
	Token string
}

func (a *AuthToken) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{"token": a.Token}, nil
}

//是否开启tls验证
func (a *AuthToken) RequireTransportSecurity() bool {
	return true
}

func demoSix() {
	//生成证书时：Common Name (eg, fully qualified host name) []:d
	tls, err := credentials.NewClientTLSFromFile("demo/grpc/tls/cert.pem", "d")
	if err != nil {
		log.Printf("credentials newclienttlsfromfile err: %v \n", err)
		return
	}

	at := &AuthToken{Token: "D"}

	conn, err := grpc.Dial("127.0.0.1:9999", grpc.WithTransportCredentials(tls), grpc.WithPerRPCCredentials(at))
	if err != nil {
		log.Printf("grpc dial err: %v \n", err)
		return
	}
	defer conn.Close()
	c := pt.NewHelloServerClient(conn)
	r, err := c.SayHello(context.Background(), &pt.HelloRequest{Name: "panda"})
	if err != nil {
		log.Printf("grpc dial err: %v \n", err)
		return
	}
	log.Println("resp: ", r.Message)
}

func demoFive() {
	//生成证书时：Common Name (eg, fully qualified host name) []:d
	tls, err := credentials.NewClientTLSFromFile("demo/grpc/tls/cert.pem", "d")
	if err != nil {
		log.Printf("credentials newclienttlsfromfile err: %v \n", err)
		return
	}

	conn, err := grpc.Dial("127.0.0.1:9999", grpc.WithTransportCredentials(tls))
	if err != nil {
		log.Printf("grpc dial err: %v \n", err)
		return
	}
	defer conn.Close()
	c := pt.NewHelloServerClient(conn)
	r, err := c.SayHello(context.Background(), &pt.HelloRequest{Name: "panda"})
	if err != nil {
		log.Printf("grpc dial err: %v \n", err)
		return
	}
	log.Println("resp: ", r.Message)
}

func demoFour() {
	conn, err := grpc.Dial("127.0.0.1:9002", grpc.WithInsecure())
	if err != nil {
		log.Printf("grpc dial err: %s \n", err.Error())
		return
	}
	defer conn.Close()
	c := pt.NewStreamServiceClient(conn)
	cs, err := c.Route(context.Background())
	if err != nil {
		log.Fatalf("client Route err:%v", err)
	}
	for n := 0; n < 6; n++ {
		err = cs.Send(&pt.StreamRequest{
			Pt: &pt.StreamPoint{
				Name:  "client",
				Value: int32(n),
			},
		})
		if err != nil {
			log.Printf("Send err:%v \n", err)
			return
		}

		resp, err := cs.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Recv err:%v \n", err)
			return
		}

		log.Printf("resp: pj.name: %s, pt.value: %d", resp.Pt.Name, resp.Pt.Value)
	}
	_ = cs.CloseSend()
}

func demoThree() {
	conn, err := grpc.Dial("127.0.0.1:9002", grpc.WithInsecure())
	if err != nil {
		log.Printf("grpc dial err: %s \n", err.Error())
		return
	}
	defer conn.Close()
	c := pt.NewStreamServiceClient(conn)
	cs, err := c.Record(context.Background())
	if err != nil {
		log.Printf("c record err: %s \n", err.Error())
		return
	}
	for n := 0; n < 6; n++ {
		err := cs.Send(&pt.StreamRequest{
			Pt: &pt.StreamPoint{
				Name:  "come",
				Value: 1,
			},
		})
		if err != nil {
			log.Printf("Send err:%v", err)
		}
	}
	resp, err := cs.CloseAndRecv()
	if err != nil {
		log.Printf("Close err:%v \n", err)
	}
	log.Printf("resp: pj.name: %s, pt.value: %d \n", resp.Pt.Name, resp.Pt.Value)
}

func demoTwo() {
	conn, err := grpc.Dial("127.0.0.1:9002", grpc.WithInsecure())
	if err != nil {
		log.Printf("grpc dial err: %s \n", err.Error())
		return
	}
	defer conn.Close()
	c := pt.NewStreamServiceClient(conn)
	sc, err := c.List(context.Background(), &pt.StreamRequest{
		Pt: &pt.StreamPoint{
			Name:  "d demo",
			Value: 1,
		},
	})
	if err != nil {
		log.Printf("sc err: %s \n", err.Error())
		return
	}
	for {
		resp, err := sc.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("recv err:%v \n", err)
			return
		}

		log.Printf("resp: name: %s, value: %d \n", resp.Pt.Name, resp.Pt.Value)
	}
}

func demoOne() {
	conn, err := grpc.Dial("127.0.0.1:9999", grpc.WithInsecure())
	if err != nil {
		log.Printf("grpc dial err: %s \n", err.Error())
		return
	}
	defer conn.Close()
	c := pt.NewHelloServerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := c.SayHello(ctx, &pt.HelloRequest{Name: "d"})
	if err != nil {
		log.Printf("grpc dial err: %s \n", err.Error())
		return
	}
	log.Println("resp: ", res.Message)
}
