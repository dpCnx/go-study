package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/source/file"
	"github.com/micro/go-micro/util/log"
	proto "github.com/micro/go-plugins/config/source/grpc/proto"
	"google.golang.org/grpc"
	"net"
	"os"
	"strings"
	"time"
)

var (
	apps = []string{"micro", "extra"}
)

func main() {
	// load config files
	err := loadConfigFile()
	if err != nil {
		log.Fatal(err)
	}

	// new service
	service := grpc.NewServer()
	proto.RegisterSourceServer(service, new(Service))
	ts, err := net.Listen("tcp", ":8600")
	if err != nil {
		log.Fatal(err)
	}

	log.Logf("configServer started")
	err = service.Serve(ts)
	if err != nil {
		log.Fatal(err)
	}
}

type Service struct{}

func (s Service) Read(ctx context.Context, req *proto.ReadRequest) (rsp *proto.ReadResponse, err error) {
	appName := parsePath(req.Path)
	switch appName {
	case "micro", "extra":
		rsp = &proto.ReadResponse{
			ChangeSet: getConfig(appName),
		}
		return
	default:
		err = fmt.Errorf("[Read] the first path is invalid")
		return
	}

	return
}

func (s Service) Watch(req *proto.WatchRequest, server proto.Source_WatchServer) (err error) {
	appName := parsePath(req.Path)
	rsp := &proto.WatchResponse{
		ChangeSet: getConfig(appName),
	}
	if err = server.Send(rsp); err != nil {
		log.Logf("[Watch] watch files error，%s", err)
		return err
	}

	return
}

func loadConfigFile() (err error) {

	p, _ := os.Getwd()

	for _, app := range apps {
		if err := config.Load(file.NewSource(
			file.WithPath(p + "/conf/" + app + ".yaml"),
		)); err != nil {
			log.Fatalf("[loadConfigFile] load files error，%s", err)
			return err
		}
	}

	// watch changes
	watcher, err := config.Watch()
	if err != nil {
		log.Fatalf("[loadConfigFile] start watching files error，%s", err)
		return err
	}

	go func() {
		for {
			v, err := watcher.Next()
			if err != nil {
				log.Fatalf("[loadConfigFile] watch files error，%s", err)
				return
			}

			log.Logf("[loadConfigFile] file change， %s", string(v.Bytes()))
		}
	}()

	return
}

func getConfig(appName string) *proto.ChangeSet {
	bytes := config.Get(appName).Bytes()

	log.Logf("[getConfig] appName，%s", string(bytes))
	return &proto.ChangeSet{
		Data:      bytes,
		Checksum:  fmt.Sprintf("%x", md5.Sum(bytes)),
		Format:    "yaml",
		Source:    "file",
		Timestamp: time.Now().Unix()}
}

func parsePath(path string) (appName string) {
	paths := strings.Split(path, "/")

	if paths[0] == "" && len(paths) > 1 {
		return paths[1]
	}

	return paths[0]
}
