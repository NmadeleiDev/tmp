package wsClient

import (
	"backend/model"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"encoding/json"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 30 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = pongWait * 9 / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024
)

var Clients map[string]*Client

type Client struct {
	Id				string
	IsOnline		bool
	Connection      *websocket.Conn
	ReadMessageChan chan model.SocketMessage
}

func	RegisterNewClient(connection *websocket.Conn, user *model.FullUserData) (client *Client) {
	client = &Client{Id: user.Id, Connection: connection, ReadMessageChan:make(chan model.SocketMessage), IsOnline: true}
	Clients[user.Id] = client
	return client
}

func GetWsMessageType(message []byte) int {
	dest := struct {
		MessageType int    `json:"messageType"`
	}{}

	if err := json.Unmarshal(message, &dest); err != nil {
		logrus.Errorf("Error unmarshal message (GetWsMessageType): %v", err)
		return 0
	}
	return dest.MessageType
}

func	SendMessageToClient(message model.SocketMessage) {
	Clients[message.ToChat].ReadMessageChan <- message
	//for _, id := range message.ToChat {
	//	if Clients[id] != nil && Clients[id].IsOnline {
	//		Clients[id].ReadMessageChan <- message
	//	}
	//}
}

