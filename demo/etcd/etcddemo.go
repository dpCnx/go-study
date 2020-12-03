package main

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"log"
	"time"
)

var (
	client *clientv3.Client
	err    error
)

func main() {
	initEtcd()

	kvClient()

}

func initEtcd() {

	config := clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}
	if client, err = clientv3.New(config); err != nil {
		log.Printf("init etcd err: %v\n", err.Error())
		return
	}

	log.Println("init etcd successful")
}

func kvClient() {

	kv := clientv3.NewKV(client)

	if putResp, err := kv.Put(context.Background(), "/test/oneKey", "oneValue"); err != nil {
		log.Printf("kv put err:%v\n", err.Error())
		return
	} else {
		log.Println("revision:", putResp.Header.Revision)
		if putResp.PrevKv != nil {
			log.Println("prevvalue:", string(putResp.PrevKv.Value))
			log.Println("prevvalue:", string(putResp.PrevKv.Version))
		}
	}

	if getResp, err := kv.Get(context.Background(), "/test/oneKey"); err != nil {
		log.Printf("kv get err:%v\n", err.Error())
		return
	} else {
		log.Println("getResp >>", getResp.Kvs)
	}

	if getResp, err := kv.Get(context.Background(), "/test/oneKey", clientv3.WithPrefix()); err != nil {
		log.Printf("kv get err:%v\n", err.Error())
		return
	} else {
		log.Println("getResp >>", getResp.Kvs)
	}

	if delResp, err := kv.Delete(context.TODO(), "/test/oneKey", clientv3.WithFromKey()); err != nil {
		log.Printf("kv delete err:%v\n", err.Error())
		return
	} else {

		log.Println(delResp)

		if len(delResp.PrevKvs) != 0 {
			for _, kvpair := range delResp.PrevKvs {
				log.Printf("delete: %s  %s \n", string(kvpair.Key), string(kvpair.Value))
			}
		}
	}

	/*
		clientv3.WithPrefix()
		clientv3.WithCountOnly()
		clientv3.WithLimit(2)
	*/

}

func opclient() {

	kv := clientv3.NewKV(client)

	putOp := clientv3.OpPut("/test/op", "opvaule")
	if opResp, err := kv.Do(context.TODO(), putOp); err != nil {
		log.Printf("kv putop err :%v\n", err.Error())
		return
	} else {
		log.Println("revision:", opResp.Put().Header.Revision)
		if opResp.Put().PrevKv != nil {
			log.Println("prevvalue:", string(opResp.Put().PrevKv.Value))
			log.Println("prevvalue:", string(opResp.Put().PrevKv.Version))
		}
	}

	getOp := clientv3.OpGet("/test/op")
	if opResp, err := kv.Do(context.TODO(), getOp); err != nil {
		log.Printf("kv getop err :%v\n", err.Error())
		return
	} else {
		log.Println("revision:", opResp.Get().Header.Revision)
		log.Println("Kvs:", opResp.Get().Kvs)
	}

	// kv.Do(op)
	// kv.Put
	// kv.Get
	// kv.Delete
}

func watchclient() {

	kv := clientv3.NewKV(client)

	go func() {
		for {
			_, err = kv.Put(context.TODO(), "/test/watch", "watchValue")
			if err != nil {
				log.Printf("kv put err:%v\n", err.Error())
				return
			}

			time.Sleep(3 * time.Second)

			_, err = kv.Delete(context.TODO(), "/test/watch")
			if err != nil {
				log.Printf("kv put err:%v\n", err.Error())
				return
			}
		}
	}()

	if getResp, err := kv.Get(context.TODO(), "/test/watch"); err != nil {
		log.Printf("kv get err:%v\n", err.Error())
		return
	} else {
		if len(getResp.Kvs) != 0 {
			log.Println("getResp:", getResp.Kvs)
		}

		watcher := clientv3.NewWatcher(client)
		watchRespChan := watcher.Watch(context.Background(), "/test/watch", clientv3.WithPrefix())

		for watchResp := range watchRespChan {
			for _, event := range watchResp.Events {
				switch event.Type {
				case mvccpb.PUT:
					log.Println("update:", string(event.Kv.Key), string(event.Kv.Key), "Revision:", event.Kv.CreateRevision, event.Kv.ModRevision)
				case mvccpb.DELETE:
					log.Println("delete:", string(event.Kv.Key), "Revision:", event.Kv.ModRevision)
				}
			}
		}
	}
}

func leaseclient() {

	// 租约lease
	lease := clientv3.NewLease(client)
	// 申请一个10秒的租约
	leaseGrantResp, err := lease.Grant(context.Background(), 10)
	if err != nil {
		log.Printf("lease grant err :%v \n", err.Error())
		return
	}
	leaseId := leaseGrantResp.ID
	log.Println("leaseId:", leaseId)
	// 自动续租
	keepRespChan, err := lease.KeepAlive(context.Background(), leaseId)
	if err != nil {
		log.Printf("lease keepalive err :%v \n", err.Error())
		return
	}
	// 处理续约应答的协程
	go func() {
		for {
			select {
			case keepResp := <-keepRespChan:
				if keepRespChan == nil {
					log.Println("租约已经失效了")
					return
				} else {
					// 自动定时的续约某个租约
					log.Println("收到自动续租应答:", keepResp.ID, "time:", time.Now())
				}
			}
		}
	}()
	kv := clientv3.NewKV(client)
	if _, err = kv.Put(context.TODO(), "/test/lease", "leaseValue", clientv3.WithLease(leaseId)); err != nil {
		log.Printf("kv put err :%v \n", err.Error())
		return
	}
	// 定时的看一下key过期了没有
	for {
		getResp, err := kv.Get(context.Background(), "/test/lease")

		if err != nil {
			log.Printf("kv get err :%v \n", err.Error())
			return
		}
		if getResp.Count == 0 {
			log.Println("kv过期")
			break
		}
		log.Println("还没过期:", getResp.Kvs)
		time.Sleep(2 * time.Second)
	}
}

func txnclient() {

	kv := clientv3.NewKV(client)
	// 创建事务
	txn := kv.Txn(context.Background())

	/*
		func CreateRevision(key string) Cmp：key=xxx的创建版本必须满足…
		func LeaseValue(key string) Cmp：key=xxx的Lease ID必须满足…
		func ModRevision(key string) Cmp：key=xxx的最后修改版本必须满足…
		func Value(key string) Cmp：key=xxx的创建值必须满足…
		func Version(key string) Cmp：key=xxx的累计更新次数必须满足…
	*/

	txnResp, err := txn.If(clientv3.Compare(clientv3.Value("/test/txn"), "=", "txvValue")).
		Then(clientv3.OpGet("/test/txn")).
		Else(clientv3.OpPut("/test/txn", "txvValue")).
		Commit()

	if err != nil {
		log.Printf("txn err:%v \n", err.Error())
		return
	}

	log.Println(txnResp.Succeeded)

	if txnResp.Succeeded {
		log.Println("~~~", txnResp.Responses[0].GetResponseRange().Kvs)
	} else {
		log.Println("!!!", txnResp.Responses[0].GetResponseRange().Kvs)
	}
}

/*
	putResp.Header.Revision --> Revision: 7  --> 每次存入一个值  Revision都会加1
	[key:"/demo/onekey" create_revision:5 mod_revision:7 version:3 value:"onekeyvaule" ]
	create_revision --> key在创建时  Revision的值为key的create_revision
	mod_revision --> mod_revision = Revision
	version --> key每改变一次 会增加一次
*/
