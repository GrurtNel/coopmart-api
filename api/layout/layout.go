package layout

import (
	"feedback/x/rest"
	"github.com/gin-gonic/gin"
)

type CampaignServer struct {
	*gin.RouterGroup
	rest.JsonRender
}

func NewCampaignServer(parent *gin.RouterGroup, name string) *CampaignServer {
	var s = CampaignServer{
		RouterGroup: parent.Group(name),
	}
	s.POST("/create", s.handleSetSetting)
	s.POST("/upload/logo", s.handleUploadLogo)
	s.POST("/upload/background", s.handleUploadBackground)
	s.GET("/get", s.handleGetSetting)
	return &s
}