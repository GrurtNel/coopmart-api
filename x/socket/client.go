package socket

import (
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

type WSGroups map[string]*WSClients
type WSClients map[string]*WSClient
type WSClient struct {
	*websocket.Conn
	msg chan []byte
}

func (b *Box) newClient(id string, conn *websocket.Conn) *WSClient {
	var client = &WSClient{
		Conn: conn,
		msg:  make(chan []byte, 0),
	}
	b.Clients[id] = client
	return client
}
func (b *Box) removeClient(id string) {
	b.Clients[id].Close()
	delete(b.Clients, id)
	glog.Info("client da roi")
	glog.Info(b.Clients)
}
func (b *Box) SendToOther(id string, msg []byte) {
	b.Clients[id].msg <- msg
}
func (client *WSClient) handleError(err error) {
	if err != nil {
		client.Conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
	}
}
