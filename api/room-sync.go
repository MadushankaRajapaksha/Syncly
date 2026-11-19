package api

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { // allow all origins
		return true

	},
}

func SyncRoom(w http.ResponseWriter, r *http.Request) {
	// opgrade protocall to websocket
	con, err := upgrader.Upgrade(w, r, nil)
	// listing msg
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return

	}
	defer con.Close()

	for {
		msgType, msg, err := con.ReadMessage()
		if err != nil {
			WriteJSON(w, 500, map[string]string{
				"error": "Erro update to Websocket!",
			})
			// break
			return

		}

		// Broadcast msg

		WriteJSON(w, 201, map[string]string{
			"message": string(msg),
		})

		// write msg to browsers
		if err = con.WriteMessage(msgType, msg); err != nil {
			return
		}
	}
}
