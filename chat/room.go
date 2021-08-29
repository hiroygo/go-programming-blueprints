package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/hiroygo/go-programming-blueprints/trace"
)

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func newRoom(t trace.Tracer) *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]struct{}),
		tracer:  t,
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
	// ロガー
	tracer trace.Tracer
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = struct{}{}
			r.tracer.Trace("新しいクライアントが参加しました")
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("クライアントが退室しました")
		case b := <-r.forward:
			r.tracer.Trace("メッセージを受信しました: ", string(b))
			for client := range r.clients {
				select {
				case client.send <- b:
					r.tracer.Trace(" -- クライアントに送信しました")

				// client.send のバッファに空きが無いときに実行される
				default:
					close(client.send)
					delete(r.clients, client)
					r.tracer.Trace(" -- クライアントに送信出来ません。クライアントを削除しました")
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
