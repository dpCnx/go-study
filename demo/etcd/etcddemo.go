package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"time"
)

/*
	https://www.jianshu.com/p/2c1f56814ea5 //简单使用
	https://blog.csdn.net/u010649766/article/details/89643699 //docker部署
	https://www.jianshu.com/p/44022c67f117 //docker部署集群
	https://blog.csdn.net/bbwangj/article/details/82584988 //介绍
*/

func main() {

}

func txnclient() {
	var (
		config         clientv3.Config
		client         *clientv3.Client
		err            error
		lease          clientv3.Lease
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId        clientv3.LeaseID
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
		keepResp       *clientv3.LeaseKeepAliveResponse
		ctx            context.Context
		cancelFunc     context.CancelFunc
		kv             clientv3.KV
		txn            clientv3.Txn
		txnResp        *clientv3.TxnResponse
	)
	// 客户端配置
	config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}
	// 建立连接
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
	}
	// lease实现锁自动过期:
	// op操作
	// txn事务: if else then
	// 1, 上锁 (创建租约, 自动续租, 拿着租约去抢占一个key)
	lease = clientv3.NewLease(client)
	// 申请一个5秒的租约
	if leaseGrantResp, err = lease.Grant(context.TODO(), 5); err != nil {
		fmt.Println(err)
	}
	// 拿到租约的ID
	leaseId = leaseGrantResp.ID
	// 准备一个用于取消自动续租的context
	ctx, cancelFunc = context.WithCancel(context.TODO())
	// 确保函数退出后, 自动续租会停止
	defer cancelFunc()
	defer lease.Revoke(context.TODO(), leaseId)
	// 5秒后会取消自动续租
	if keepRespChan, err = lease.KeepAlive(ctx, leaseId); err != nil {
		fmt.Println(err)
	}
	// 处理续约应答的协程
	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepRespChan == nil {
					fmt.Println("租约已经失效了")
					return
				} else { // 每秒会续租一次, 所以就会受到一次应答
					fmt.Println("收到自动续租应答:", keepResp.ID)
				}
			}
		}
	}()
	//  if 不存在key， then 设置它, else 抢锁失败
	kv = clientv3.NewKV(client)
	// 创建事务
	txn = kv.Txn(context.TODO())
	// 定义事务
	// 如果key不存在
	// /demo/lock这个key对应的value”必须等于”lockvaule ” 返回的是成功就执行then 失败就执行else
	txn.If(clientv3.Compare(clientv3.CreateRevision("/demo/lock"), "=", 0)).
		Then(clientv3.OpPut("/demo/lock", "lockvaule", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet("/demo/lock"))
	// 提交事务
	if txnResp, err = txn.Commit(); err != nil {
		fmt.Println("err")
	}
	fmt.Println("res:", txnResp.Succeeded)
	// 判断是否抢到了锁
	if !txnResp.Succeeded {
		fmt.Println("锁被占用:", string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
	}
	// 2, 处理业务
	fmt.Println("处理任务")
	time.Sleep(50 * time.Second)
	// 3, 释放锁(取消自动续租, 释放租约)
	// defer 会把租约释放掉, 关联的KV就被删除了
}

func opclient() {
	var (
		config clientv3.Config
		client *clientv3.Client
		err    error
		kv     clientv3.KV
		putOp  clientv3.Op
		getOp  clientv3.Op
		opResp clientv3.OpResponse
	)
	// 客户端配置
	config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}
	// 建立连接
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}
	kv = clientv3.NewKV(client)
	// 创建Op: operation
	putOp = clientv3.OpPut("/demo/op", "opvaule")
	// 执行OP
	if opResp, err = kv.Do(context.TODO(), putOp); err != nil {
		fmt.Println(err)
		return
	}
	// kv.Do(op)
	// kv.Put
	// kv.Get
	// kv.Delete
	fmt.Println("写入Revision:", opResp.Put().Header.Revision)
	// 创建Op
	getOp = clientv3.OpGet("/demo/op")
	// 执行OP
	if opResp, err = kv.Do(context.TODO(), getOp); err != nil {
		fmt.Println(err)
		return
	}
	// 打印
	fmt.Println("数据Revision:", opResp.Get().Kvs[0].ModRevision)
	// create rev == mod rev
	fmt.Println("数据value:", string(opResp.Get().Kvs[0].Value))
}

func watchclient() {
	var (
		config             clientv3.Config
		client             *clientv3.Client
		kv                 clientv3.KV
		getResp            *clientv3.GetResponse
		watcher            clientv3.Watcher
		watchRespChan      <-chan clientv3.WatchResponse
		watchResp          clientv3.WatchResponse
		watchStartRevision int64
		event              *clientv3.Event
		err                error
	)
	config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}
	if client, err = clientv3.New(config); err != nil {
		fmt.Println("new client err:", err)
	}
	kv = clientv3.NewKV(client)
	// 模拟etcd中KV的变化
	go func() {
		for {
			kv.Put(context.TODO(), "/demo/watch", "watchvaule")

			time.Sleep(1 * time.Second)

			kv.Delete(context.TODO(), "/demo/watch")
		}
	}()
	if getResp, err = kv.Get(context.TODO(), "/demo/watch"); err != nil {
		fmt.Println(err)
		return
	}
	if len(getResp.Kvs) != 0 {
		fmt.Println("当前值:", string(getResp.Kvs[0].Value))
	}
	// 当前etcd集群事务ID, 单调递增的
	watchStartRevision = getResp.Header.Revision + 1
	watcher = clientv3.NewWatcher(client)
	watchRespChan = watcher.Watch(context.Background(), "/demo/watch", clientv3.WithRev(watchStartRevision))
	// 处理kv变化事件
	for watchResp = range watchRespChan {
		for _, event = range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("修改为:", string(event.Kv.Value), "Revision:", event.Kv.CreateRevision, event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println("删除了", "Revision:", event.Kv.ModRevision)
			}
		}
	}
}

func leaseclient() {
	var (
		config         clientv3.Config
		client         *clientv3.Client
		lease          clientv3.Lease
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId        clientv3.LeaseID
		kv             clientv3.KV
		getResp        *clientv3.GetResponse
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
		keepResp       *clientv3.LeaseKeepAliveResponse
		err            error
	)
	config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}
	if client, err = clientv3.New(config); err != nil {
		fmt.Println("new client err:", err)
	}
	// 租约lease
	lease = clientv3.NewLease(client)
	// 申请一个10秒的租约
	if leaseGrantResp, err = lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println(err)
	}
	leaseId = leaseGrantResp.ID
	fmt.Println("leaseId:", leaseId)
	// 自动续租
	if keepRespChan, err = lease.KeepAlive(context.TODO(), leaseId); err != nil {
		fmt.Println(err)
	}
	// 处理续约应答的协程
	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepRespChan == nil {
					fmt.Println("租约已经失效了")
					return
				} else { // 自动定时的续约某个租约
					fmt.Println("收到自动续租应答:", time.Now())
				}
			}
		}
	}()
	kv = clientv3.NewKV(client)
	if _, err = kv.Put(context.TODO(), "/demo/leaseone", "leaseonevaule", clientv3.WithLease(leaseId)); err != nil {
		fmt.Println(err)
	}
	// 定时的看一下key过期了没有
	for {
		if getResp, err = kv.Get(context.TODO(), "/demo/leaseone"); err != nil {
			fmt.Println(err)
			return
		}
		if getResp.Count == 0 {
			fmt.Println("kv过期了")
			break
		}
		fmt.Println("还没过期:", getResp.Kvs)
		time.Sleep(2 * time.Second)
	}
}

func kvClient() {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		kv      clientv3.KV
		putResp *clientv3.PutResponse
		getResp *clientv3.GetResponse
		delResp *clientv3.DeleteResponse
		kvpair  *mvccpb.KeyValue
		err     error
	)
	config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}
	if client, err = clientv3.New(config); err != nil {
		fmt.Println("new client err:", err)
	}
	// 读写etcd的键值对
	kv = clientv3.NewKV(client)
	if putResp, err = kv.Put(context.Background(), "/demo/twokey", "twokeyvaule"); err != nil {
		fmt.Println("kv put err:", err)
	} else {
		fmt.Println("Revision:", putResp.Header.Revision)
		if putResp.PrevKv != nil {
			fmt.Println("PrevValue:", string(putResp.PrevKv.Value))
			fmt.Println("PrevValue:", string(putResp.PrevKv.Version))
		}
	}
	if getResp, err = kv.Get(context.Background(), "/demo/twokey"); err != nil {
		fmt.Println("kv get err:", err)
	} else {
		fmt.Println(getResp.Kvs)
	}

	// clientv3.WithLimit(2)
	if delResp, err = kv.Delete(context.TODO(), "/demo/onekey", clientv3.WithFromKey()); err != nil {
		fmt.Println(err)
	} else {
		if len(delResp.PrevKvs) != 0 {
			for _, kvpair = range delResp.PrevKvs {
				fmt.Println("删除了:", string(kvpair.Key), string(kvpair.Value))
			}
		}
	}

	/*
		if getResp, err = kv.Get(context.Background(), "/demo/twokey", clientv3.WithPrefix()); err != nil {
			fmt.Println("kv get err:", err)
		} else {
			fmt.Println(getResp.Kvs)
		}

		clientv3.WithPrefix()
		clientv3.WithCountOnly()
	*/

	/*
		putResp.Header.Revision --> Revision: 7  --> 每次存入一个值  Revision都会加1
		[key:"/demo/onekey" create_revision:5 mod_revision:7 version:3 value:"onekeyvaule" ]
		create_revision --> key在创建时  Revision的值为key的create_revision
		mod_revision --> mod_revision = Revision
		version --> key每改变一次 会增加一次
	*/
}
