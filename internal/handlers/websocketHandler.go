package handlers

import (
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Error(err)
		return
	}

	log.Println("Client connected")

	err = ws.WriteMessage(websocket.TextMessage, []byte("Hi Client!"))
	if err != nil {
		log.Error(err)
	}

	reader(ws)
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Error(err)
			return
		}

		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, []byte(fmt.Sprintf("your message: %v", string(p)))); err != nil {
			log.Error(err)
			return
		}
	}
}
