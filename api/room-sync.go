package api

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { // allow all origins
		return true

	},
}

type UserInfo struct {
	UserId         int    `json:"userId"`
	UserName       string `json:"userName"`
	Email          string `json:"email"`
	ProfilePicture string `json:"profilePicture"`
}

type VideoState struct {
	Action      string  `json:"action"`
	IsPlaying   bool    `json:"isPlaying"`
	CurrentTime float64 `json:currentTime`
}
type Message struct {
	Text      string `json:"text"`
	TimeStamp string `json:"timestamp,omitempty"`
}

type WSMessage struct {
	Type        string      `json:"type"`
	RoomId      string      `json:"roomId"`
	Sender      *UserInfo   `json:"sender,omitempty"`
	Video       *VideoState `json:"video,omitempty"`
	Chat        *Message    `json:"chat,omitempty"`
	OnlineUsers []UserInfo  `json:"onlineUsers,omitempty"`
}

// handle Clinents with goroutine
type Client struct {
	Conn    *websocket.Conn
	User    UserInfo
	WriteMu sync.Mutex
}

type WsRoom struct {
	Clients      map[*Client]bool
	CurrentVideo *VideoState
}

var (
	rooms   = make(map[string]*WsRoom)
	roomsMu sync.RWMutex
)

// create room for cliet
func getOrCreateRoom(roomId string) *WsRoom {
	roomsMu.Lock()

	defer roomsMu.Unlock()

	if rooms[roomId] == nil {
		rooms[roomId] = &WsRoom{
			Clients: make(map[*Client]bool),
		}
	}
	return rooms[roomId]
}

func SyncRoom(w http.ResponseWriter, r *http.Request) {
	// getroom id

	roomId := mux.Vars(r)["roomId"]

	// opgrade protocall to websocket
	con, err := upgrader.Upgrade(w, r, nil)
	// listing msg
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return

	}
	//create client space and  room
	client := &Client{Conn: con}
	room := getOrCreateRoom(roomId)

	// add room to client
	roomsMu.Lock()
	room.Clients[client] = true
	roomsMu.Unlock()

	// after user deisconnect cleanup

	defer func() {
		roomsMu.Lock()
		delete(room.Clients, client)
		// handle room clients are emty
		isEmpty := len(room.Clients) == 0

		if isEmpty {
			delete(rooms, roomId)
		}

		roomsMu.Unlock()

		con.Close()

		if !isEmpty {
			// broadcastonline users list
			broadcastOnlineUsers(roomId)

		}
	}()
	// masage listing with create reltime listing
	for {
		var msg WSMessage
		if err := con.ReadJSON(&msg); err != nil {
			log.Printf("[%s] Client disconnected: %v", roomId, err)
			break

		}

		msg.RoomId = roomId

		switch msg.Type {
		// user join
		case "join":
			// get userdata
			if msg.Sender != nil {
				client.User = *msg.Sender

			}
			// add user to online user lit
			broadcastOnlineUsers(roomId)

			// already video loaded in room
			roomsMu.RLock()
			lastVideo := room.CurrentVideo
			roomsMu.RUnlock()

			if lastVideo != nil {
				_ = client.sendJSON(WSMessage{
					Type:   "sync_state",
					RoomId: roomId,
					Video:  lastVideo})
			}
		// video states changes
		case "video":
			// update room latest video state
			if msg.Video != nil {
				roomsMu.Lock()

				room.CurrentVideo = msg.Video
				roomsMu.Unlock()
			}
			broadcastToRoom(roomId, msg, client)
		// chat broadcast
		case "chat":
			broadcastToRoom(roomId, msg, nil)
		}
	}
}

// send json msg to client

func (c *Client) sendJSON(msg WSMessage) error {
	c.WriteMu.Lock()
	defer c.WriteMu.Unlock()
	return c.Conn.WriteJSON(msg)
}

// realtime update online users list
func broadcastOnlineUsers(roomId string) {
	roomsMu.RLock()
	room := rooms[roomId]
	if room == nil {
		roomsMu.RUnlock()
		return
	}

	var onlineList []UserInfo

	clients := make([]*Client, 0, len(room.Clients))

	for c := range room.Clients {
		clients = append(clients, c)
		if c.User.UserName != "" {
			onlineList = append(onlineList, c.User)
		}
	}

	roomsMu.RUnlock()

	msg := WSMessage{
		Type:        "users",
		RoomId:      roomId,
		OnlineUsers: onlineList,
	}

	//broadcast

	for _, c := range clients {
		_ = c.sendJSON(msg)
	}

}

// broadcast masg to all users in room

func broadcastToRoom(roomId string, msg WSMessage, exclude *Client) {
	roomsMu.RLock()

	room := rooms[roomId]

	if room == nil {
		roomsMu.RUnlock()
		return
	}

	// copy client data and lock realse
	clients := make([]*Client, 0, len(room.Clients))

	for c := range room.Clients {
		if c != exclude {
			clients = append(clients, c)
		}
	}

	roomsMu.RUnlock()

	for _, c := range clients {
		// sendmsg
		if err := c.sendJSON(msg); err != nil {
			log.Println("Broadcast write error:", err)

		}

	}

}
