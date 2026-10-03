package wshelper

import (
	"github.com/gofiber/contrib/websocket"
)

type InitialWebsocketMessage struct {
	Authorization string `json:"authorization"`
}

func InitialCheck(conn *websocket.Conn) (InitialWebsocketMessage, error) {
	var msg InitialWebsocketMessage
	if err := conn.ReadJSON(&msg); err != nil {
		return msg, err
	}
	return msg, nil
}
