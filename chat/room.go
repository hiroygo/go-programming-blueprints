package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]struct{}),
	}
}

type room struct {
	// ある client のメッセージを別の client に送る
	forward chan []byte
	// room に参加しようとしている client のためのチャネル
	join chan *client
	// room から退室しようとしている client のためのチャネル
	leave chan *client
	// 在室しているすべての client
	clients map[*client]struct{}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = struct{}{}
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
		case b := <-r.forward:
			for client := range r.clients {
				select {
				case client.send <- b:
				// client.send のバッファに空きが無いときに実行される
				default:
					close(client.send)
					delete(r.clients, client)
				}
			}
		}
	}
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// WebSocket ではクライアントからハンドシェイクを開始する
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal(err)
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
