package esvc

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"log"
	"time"
)

var cli *clientv3.Client

const schema = "ns"

func Register(name string, addr string, ttl int64) error {

	var err error
	if cli == nil {
		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   []string{"127.0.0.1:2379"},
			DialTimeout: 5 * time.Second,
		})
		if err != nil {
			fmt.Printf("connect to etcd err:%s", err)
			return err
		}
	}

	ticker := time.NewTicker(time.Second * time.Duration(ttl))

	go func() {
		for {
			//  /ns/sername/addr
			getResp, err := cli.Get(context.Background(), "/"+schema+"/"+name+"/"+addr)
			if err != nil {
				log.Println(err)
				fmt.Printf("Register:%s", err)
			} else if getResp.Count == 0 {
				err = withAlive(name, addr, ttl)
				if err != nil {
					log.Println(err)
					fmt.Printf("keep alive:%s", err)
				}
			} else {
				//fmt.Println(getResp.Kvs)
			}

			<-ticker.C
		}
	}()

	return nil
}

func withAlive(name string, addr string, ttl int64) error {

	leaseResp, err := cli.Grant(context.Background(), ttl)
	if err != nil {
		return err
	}

	_, err = cli.Put(context.Background(), "/"+schema+"/"+name+"/"+addr, addr, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		fmt.Printf("put etcd error:%s", err)
		return err
	}

	ch, err := cli.KeepAlive(context.Background(), leaseResp.ID)
	if err != nil {
		fmt.Printf("keep alive error:%s", err)
		return err
	}

	go func() {
		for {
			_ = <-ch
			if ch == nil {
				fmt.Println("租约已经失效了")
				return
			} else { // 每秒会续租一次, 所以就会受到一次应答
				//fmt.Println("收到自动续租应答:", keepResp.ID)
			}
		}
	}()

	return nil
}

func UnRegister(name string, addr string) {
	if cli != nil {
		_, _ = cli.Delete(context.Background(), "/"+schema+"/"+name+"/"+addr)
	}
}
