package main

import (
	"github.com/gorilla/websocket"
	"log"
)

type WbClient struct {
	WbId string
	Coon *websocket.Conn
	pool *WbPool
}

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

func (w *WbClient) read() {

	defer func() {
		w.Coon.Close()
		w.pool.UnRegister <- w
		log.Println("w.Coon.Close()")
	}()

	for {
		mgType, p, err := w.Coon.ReadMessage()
		if err != nil {
			log.Println("w.Coon.ReadMessage() err:", err)
			return
		}

		msg := &Message{
			Type: mgType,
			Body: string(p),
		}

		w.pool.BroadCastMsg <- msg

		log.Println("w.Coon.ReadMessage():", msg)
	}
}
