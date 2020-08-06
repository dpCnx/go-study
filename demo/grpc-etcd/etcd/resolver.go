package esvc

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"google.golang.org/grpc/resolver"
	"log"
	"strings"
	"time"
)

type etcdResolver struct {
	rawAddr string
	cc      resolver.ClientConn
}

func NewResolver(etcdAddr string) resolver.Builder {
	return &etcdResolver{rawAddr: etcdAddr}
}

func (r *etcdResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {

	var err error
	if cli == nil {
		cli, err = clientv3.New(clientv3.Config{
			Endpoints:   strings.Split(r.rawAddr, ";"),
			DialTimeout: 15 * time.Second,
		})
		if err != nil {
			return nil, err
		}
	}

	r.cc = cc

	fmt.Println("target.Scheme==>", target.Scheme)
	fmt.Println("target.Endpoint==>", target.Endpoint)

	go r.watch("/" + target.Scheme + "/" + target.Endpoint + "/")

	return r, nil
}

func (r etcdResolver) Scheme() string {
	return schema
}

func (r etcdResolver) ResolveNow(rn resolver.ResolveNowOptions) {
	log.Println("ResolveNow") // TODO check
}

func (r etcdResolver) Close() {
	log.Println("Close")
}

func (r *etcdResolver) watch(keyPrefix string) {

	var addrList []resolver.Address

	getResp, err := cli.Get(context.Background(), keyPrefix, clientv3.WithPrefix())
	if err != nil {
		log.Println(err)
	} else {
		for i := range getResp.Kvs {
			fmt.Println("kvs:", getResp.Kvs)
			fmt.Println(strings.TrimPrefix(string(getResp.Kvs[i].Key), keyPrefix))
			addrList = append(addrList, resolver.Address{Addr: strings.TrimPrefix(string(getResp.Kvs[i].Key), keyPrefix)})
		}
	}

	r.cc.NewAddress(addrList)

	rch := cli.Watch(context.Background(), keyPrefix, clientv3.WithPrefix())
	for n := range rch {
		for _, ev := range n.Events {
			addr := strings.TrimPrefix(string(ev.Kv.Key), keyPrefix)
			switch ev.Type {
			case mvccpb.PUT:
				if !exist(addrList, addr) {
					addrList = append(addrList, resolver.Address{Addr: addr})
					r.cc.NewAddress(addrList)
				}
			case mvccpb.DELETE:
				if s, ok := remove(addrList, addr); ok {
					addrList = s
					r.cc.NewAddress(addrList)
				}
			}
			//log.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}

func exist(l []resolver.Address, addr string) bool {
	for i := range l {
		if l[i].Addr == addr {
			return true
		}
	}
	return false
}

func remove(s []resolver.Address, addr string) ([]resolver.Address, bool) {
	for i := range s {
		if s[i].Addr == addr {
			s[i] = s[len(s)-1]
			return s[:len(s)-1], true
		}
	}
	return nil, false
}
