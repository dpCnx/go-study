package main

import (
	"log"
	"net/http"
)

var wbpool *WbPool

func main() {

	wbpool = newWbpool()
	go wbpool.start()

	http.HandleFunc("/ws", wsServer)

	http.ListenAndServe(":9999", nil)
}

func wsServer(writer http.ResponseWriter, request *http.Request) {

	coon, err := upgrade(writer, request)
	if err != nil {
		log.Println("ws server:", err)
		return
	}

	client := &WbClient{
		Coon: coon,
		pool: wbpool,
	}

	wbpool.Register <- client

	client.read()
}
