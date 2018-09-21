package room

import (
	"encoding/json"
	"feedback/o/result"
	"github.com/gorilla/websocket"
)

type Hub struct {
	Transporter      chan *result.SurveyResult
	CyberTransporter chan interface{}
}

var hub *Hub

func NewHub() *Hub {
	if hub == nil {
		hub = &Hub{
			Transporter:      make(chan *result.SurveyResult),
			CyberTransporter: make(chan interface{}),
		}
	}
	return hub
}

func (h *Hub) Loop() {
	for {
		select {
		case <-h.Transporter:
			if clients["admin"] != nil {
				var storeReport, _ = result.GetQuantityReportRealtime("store")
				var staffReport, _ = result.GetQuantityReportRealtime("uname")
				var poorReport, _ = result.GetPoorFeedback(true)
				var inPoorReport, _ = result.GetPoorFeedback(false)
				var payload = map[string]interface{}{
					"store_report":           storeReport,
					"staff_report":           staffReport,
					"poor_feedback_report":   poorReport,
					"inpoor_feedback_report": inPoorReport,
				}
				var payloadByte, _ = json.Marshal(payload)
				clients["admin"].WriteMessage(websocket.TextMessage, payloadByte)
			}
			// clients["admin"].WriteMessage(websocket.TextMessage, payload)
		case <-h.CyberTransporter:
			if clients["admin"] != nil {
				var storeReport, _ = result.GetQuantityReportRealtime("store")
				var staffReport, _ = result.GetQuantityReportRealtime("uname")
				var poorReport, _ = result.GetPoorFeedback(true)
				var inPoorReport, _ = result.GetPoorFeedback(false)
				var payload = map[string]interface{}{
					"store_report":           storeReport,
					"staff_report":           staffReport,
					"poor_feedback_report":   poorReport,
					"inpoor_feedback_report": inPoorReport,
				}
				var payloadByte, _ = json.Marshal(payload)
				clients["admin"].WriteMessage(websocket.TextMessage, payloadByte)
			}
			// clients["admin"].WriteMessage(websocket.TextMessage, payload)
		}
	}
}
