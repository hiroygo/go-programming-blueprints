package main

import (
	"time"

	"github.com/gorilla/websocket"
)

type client struct {
	socket *websocket.Conn
	// room から届いたメッセージを受信する
	// close は room が行う
	roomMsg chan *message
	// この client が参加している room
	room     *room
	userData map[string]interface{}
}

func (c *client) read() {
	for {
		msg := &message{}
		err := c.socket.ReadJSON(msg)
		// 相手のソケットが閉じると抜ける
		if err != nil {
			break
		}

		msg.When = time.Now()
		msg.Name = c.userData["name"].(string)
		c.room.forward <- msg
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.roomMsg {
		err := c.socket.WriteJSON(msg)
		// 相手のソケットが閉じると抜ける
		if err != nil {
			break
		}
	}
	c.socket.Close()
}
