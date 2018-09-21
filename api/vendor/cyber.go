package vendor

import (
	"feedback/room"
	"feedback/x/rest"
	"github.com/gin-gonic/gin"
)

var hub = room.NewHub()

type CyberServer struct {
	*gin.RouterGroup
	rest.JsonRender
}

func NewCyberServer(parent *gin.RouterGroup, name string) *CyberServer {
	var s = CyberServer{
		RouterGroup: parent.Group(name),
	}
	s.GET("checkin", s.checkinCar)
	return &s
}

func (s *CyberServer) checkinCar(ctx *gin.Context) {
	var carIdentity = ctx.Query("car_identity")
	hub.CyberTransporter <- carIdentity
	s.SendData(ctx, carIdentity)
}
