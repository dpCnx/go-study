package demo

import (
	"github.com/go-acme/lego/v3/log"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/file"
	"os"
)

func main() {

	p, _ := os.Getwd()

	if err := config.Load(file.NewSource(
		file.WithPath(p + "/demo/micro/file/demo/config.json"),
	)); err != nil {
		log.Println(err)
		return
	}
	var host Host

	if err := config.Get("hosts", "database").Scan(&host); err != nil {
		log.Println(err)
		return
	}
	log.Println(host.Address, host.Port)
}

type Host struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}
