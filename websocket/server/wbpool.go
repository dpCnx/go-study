package main

import (
	"fmt"
	"log"
)

type WbPool struct {
	Register     chan *WbClient
	UnRegister   chan *WbClient
	ClientMaps   map[*WbClient]bool
	BroadCastMsg chan *Message
}

func newWbpool() *WbPool {
	return &WbPool{
		Register:     make(chan *WbClient),
		UnRegister:   make(chan *WbClient),
		ClientMaps:   make(map[*WbClient]bool),
		BroadCastMsg: make(chan *Message),
	}
}

func (wbpool *WbPool) start() {

	for {
		select {

		case client := <-wbpool.Register:
			wbpool.ClientMaps[client] = true
			log.Println("Size of Connection Pool: ", len(wbpool.ClientMaps))
			for client, _ := range wbpool.ClientMaps {
				fmt.Println(client)
				client.Coon.WriteJSON(Message{Type: 1, Body: "New User Joined..."})
			}

		case client := <-wbpool.UnRegister:
			delete(wbpool.ClientMaps, client)
			log.Println("Size of Connection Pool: ", len(wbpool.ClientMaps))
			for client, _ := range wbpool.ClientMaps {
				fmt.Println(client)
				client.Coon.WriteJSON(Message{Type: 1, Body: "User Disconnected..."})
			}

		case msg := <-wbpool.BroadCastMsg:
			log.Println("Sending message:", msg.Body)
			for client, _ := range wbpool.ClientMaps {
				if err := client.Coon.WriteJSON(msg); err != nil {
					fmt.Println(err)
					return
				}
			}

		}
	}
}
