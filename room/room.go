package room

import (
	"encoding/json"
	"feedback/o/result"
	"feedback/x/rest"
	"feedback/x/socket"
	"fmt"
	"g/x/math"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type RoomServer struct {
	*gin.RouterGroup
	rest.JsonRender
	box *socket.Box
}

func NewRoomServer(parent *gin.RouterGroup) *RoomServer {
	var s = RoomServer{
		RouterGroup: parent,
		box:         socket.NewBox(),
	}
	s.box.Handle("/chat", s.handleChat)
	s.GET("/join", s.handleJoin)
	return &s
}

func (s *RoomServer) handleJoin(ctx *gin.Context) {
	s.serveWS(ctx.Writer, ctx.Request)
}

func (s *RoomServer) handleChat(r *socket.Request) {
	var msg = struct {
		To   string
		Text string
	}{}
	r.MustDecodeBody(&msg)
	glog.Info(msg)
	s.box.SendToOther(msg.To, []byte(msg.Text))
}

var clients = map[string]*websocket.Conn{}

var (
	storeReport []*result.QuantityReport
	staffReport []*result.QuantityReport
	poorReport  []*result.SurveyResult
)

//update realtime
func updateStoreReport(r *result.SurveyResult) {
	for _, item := range storeReport {
		if item.Actor == r.Store {

		}
	}
}

//end update global retime
func (s *RoomServer) serveWS(w http.ResponseWriter, r *http.Request) {
	var id = r.URL.Query().Get("id")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	if id != "admin" {
		id = math.RandNumString(6)
	} else {
		go func() {
			time.Sleep(1000)
			storeReport, _ = result.GetQuantityReportRealtime("store")
			staffReport, _ = result.GetQuantityReportRealtime("uname")
			poorReport, _ = result.GetPoorFeedback(true)
			var inPoorReport, _ = result.GetPoorFeedback(false)
			var payload = map[string]interface{}{
				"store_report":           storeReport,
				"staff_report":           staffReport,
				"poor_feedback_report":   poorReport,
				"inpoor_feedback_report": inPoorReport,
			}
			var payloadByte, _ = json.Marshal(payload)
			conn.WriteMessage(websocket.TextMessage, payloadByte)
		}()
	}
	clients[id] = conn
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			delete(clients, id)
			glog.Info("socket closed", err)
			return
		}
		if id == "admin" {
			sendToAll()
		}
	}
}

func sendToAdmin(data interface{}) {
	if client, ok := clients["admin"]; ok {
		var payload, _ = json.Marshal(data)
		client.WriteMessage(websocket.TextMessage, payload)
	}
}

func sendToAll() {
	for id, item := range clients {
		if id != "admin" {
			if err := item.WriteMessage(websocket.TextMessage, []byte("ping")); err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
