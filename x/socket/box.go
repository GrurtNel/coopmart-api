package socket

import (
	"encoding/json"
	"strings"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

type Box struct {
	Clients  WSClients
	Groups   WSGroups
	handlers Handlers
}
type Handlers map[string]func(r *Request)

func NewBox() *Box {
	return &Box{
		Clients:  WSClients{},
		Groups:   WSGroups{},
		handlers: Handlers{},
	}
}
func (b *Box) Accept(id string, conn *websocket.Conn) {
	defer func() {
		b.removeClient(id)
		rc := recover()
		if err, ok := rc.(error); ok {
			b.Recover(conn, err)
		}
	}()
	client := b.newClient(id, conn)
	glog.Info(b.Clients)
	//send data
	go func() {
		for {
			data := <-client.msg
			err := conn.WriteMessage(websocket.TextMessage, data)
			client.handleError(err)
		}
	}()
	//recive data
	for {
		if _, msg, err := conn.ReadMessage(); err != nil {
			glog.Error(err)
			panic(err)
		} else {
			b.serve(NewRequest(msg))
		}
	}

}

func (b *Box) Recover(conn *websocket.Conn, err error) {
	conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
}

func (b *Box) serve(r *Request) {
	if handler, ok := b.handlers[r.URI]; ok {
		handler(r)
	}
}
func (b *Box) Handle(path string, handler func(r *Request)) {
	b.handlers[path] = handler
}

//request
type Request struct {
	URI  string
	Data string
}

func NewRequest(payload []byte) *Request {
	var temp = strings.Split(string(payload), " ")
	return &Request{
		URI:  temp[0],
		Data: temp[1],
	}
}
func (r *Request) MustDecodeBody(object interface{}) {
	json.Unmarshal([]byte(r.Data), object)
}
