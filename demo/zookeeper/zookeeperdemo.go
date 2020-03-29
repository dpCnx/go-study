package main

import (
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"time"
)

func main() {
	//conn := getConn()
	//conn := watchConn()
	//getAllNodes(conn)
	//createNode(conn)
	//updateNode(conn)
	//deleteNode(conn)
	//watchNodeChild(conn)
}

/*
	获取连接
*/
func getConn() *zk.Conn {
	var hosts = []string{"192.168.60.100:2181", "192.168.60.101:2181", "192.168.60.102:2181"} //server端host
	conn, _, err := zk.Connect(hosts, time.Second*5)
	//defer conn.Close()
	if err != nil {
		log.Printf("connect err:%v", err)
		return nil
	}
	return conn
}

/*
	获取连接 watch -->没成功
*/
func watchConn() *zk.Conn {
	option := zk.WithEventCallback(func(event zk.Event) {
		log.Printf("path:%s", event.Path)
		log.Printf("type:%s", event.Type.String())
		log.Printf("state:%s", event.State.String())
	})

	var hosts = []string{"192.168.60.100:2181", "192.168.60.101:2181", "192.168.60.102:2181"} //server端host
	conn, _, err := zk.Connect(hosts, time.Second*5, option)
	if err != nil {
		log.Printf("connect err:%v", err)
		return nil
	}
	//defer conn.Close()
	return conn
}

/*
	获取节点
*/
func getAllNodes(conn *zk.Conn) {
	children, stat, err := conn.Children("/")
	if err != nil {
		log.Printf("connect childrenw err:%v", err)
		return
	}

	log.Printf("child:%v\n", children)
	log.Printf("stat:%+v\n", stat)
}

/*
	创建节点
*/
func createNode(conn *zk.Conn) {
	//参数3：模式 	FlagEphemeral 临时节点  FlagSequence 永久节点，添加序列号
	//参数4：控制访问权限模式
	node, err := conn.Create("/test", []byte("hello"), zk.FlagSequence, zk.WorldACL(zk.PermAll))
	if err != nil {
		log.Printf("create err:%v", err)
		return
	}

	log.Printf("node:%v", node)
}

/*
	修改节点
*/
func updateNode(conn *zk.Conn) {
	b, stat, err := conn.Exists("/test0000000001")
	if err != nil {
		log.Printf("conn exists err:%v", err)
		return
	}

	log.Printf("stat:%+v", stat)

	if b {
		stat, err = conn.Set("/test0000000001", []byte("hello2"), stat.Version)
		if err != nil {
			log.Printf("conn set err:%v", err)
			return
		}

		log.Printf("stat:%+v", stat)
	}
}

/*
	修改节点
*/
func deleteNode(conn *zk.Conn) {
	b, stat, err := conn.Exists("/test0000000001")
	if err != nil {
		log.Printf("conn exists err:%v", err)
		return
	}

	log.Printf("stat:%+v", stat)

	if b {
		err = conn.Delete("/test0000000001", stat.Version)
		if err != nil {
			log.Printf("conn set err:%v", err)
			return
		}
	}
}

/*
	watch 子节点 -->不好用
	shell:get /test watch
*/

func watchNodeChild(conn *zk.Conn) {
	childrens, stat, event, err := conn.ChildrenW("/")
	if err != nil {
		log.Printf("connect childrenw err:%v", err)
		return
	}

	log.Printf("childrens:%v\n", childrens)
	log.Printf("stat:%+v\n", stat)

	go func(e <-chan zk.Event) {

		for {
			select {
			case e := <-event:
				switch e.Type {
				case zk.EventNodeCreated:
					log.Println("EventNodeCreated")
					log.Printf("path:%s\n", e.Path)
				case zk.EventNodeDeleted:
					log.Println("EventNodeDeleted")
					log.Printf("path:%s\n", e.Path)
				case zk.EventNodeDataChanged:
					log.Println("EventNodeDataChanged")
					log.Printf("path:%s\n", e.Path)
				case zk.EventNodeChildrenChanged:
					log.Println("EventNodeChildrenChanged")
					log.Printf("path:%s\n", e.Path)
				}

			}

			time.Sleep(10 * time.Second)
		}
	}(event)

	for {
		log.Println("loading")
		time.Sleep(15 * time.Second)
	}
}
