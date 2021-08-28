package main

import (
	"github.com/gorilla/websocket"
)

type client struct {
	socket *websocket.Conn
	// room から届いたメッセージを受信する
	// close は room が行う
	send chan []byte
	// この client が参加している room
	room *room
}

func (c *client) read() {
	for {
		// 相手のソケットが閉じると抜ける
		_, b, err := c.socket.ReadMessage()
		if err != nil {
			break
		}
		c.room.forward <- b
	}
	c.socket.Close()
}

func (c *client) write() {
	for b := range c.send {
		// 相手のソケットが閉じると抜ける
		err := c.socket.WriteMessage(websocket.TextMessage, b)
		if err != nil {
			break
		}
	}
	c.socket.Close()
}
