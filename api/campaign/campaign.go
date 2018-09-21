package campaign

import (
	"feedback/o/campaign"
	"feedback/x/rest"
	"fmt"
	"strconv"

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
	s.GET("/list", s.handleList)
	s.POST("/create", s.handleCreate)
	s.POST("/update", s.handleUpdate)
	s.GET("/delete", s.handleDelete)
	s.POST("/device/create", s.handleAddDevice)
	s.GET("/device/get", s.handleGetCampaign)
	s.GET("/channel/list", s.handleListChannel)
	return &s
}

func (s *CampaignServer) handleListChannel(ctx *gin.Context) {
	var channels = campaign.GetAllChannelLink()
	s.SendData(ctx, channels)
}

func (s *CampaignServer) handleAddDevice(ctx *gin.Context) {
	var c *campaign.Campaign
	ctx.BindJSON(&c)
	rest.AssertNil(c.Create())
	s.SendData(ctx, c)
}
func (s *CampaignServer) handleCreate(ctx *gin.Context) {
	var c *campaign.Campaign
	rest.AssertNil(ctx.BindJSON(&c))
	rest.AssertNil(c.Create())
	s.SendData(ctx, c)
}
func (s *CampaignServer) handleUpdate(ctx *gin.Context) {
	var c *campaign.Campaign
	ctx.BindJSON(&c)
	rest.AssertNil(campaign.UpdateByID(c))
	s.SendData(ctx, c)
}

func (s *CampaignServer) handleList(ctx *gin.Context) {
	campaigns, err := campaign.GetCampaigns()
	rest.AssertNil(err)
	s.SendData(ctx, campaigns)
}

func (s *CampaignServer) handleDelete(ctx *gin.Context) {
	var id = ctx.Query("id")
	fmt.Println(id)
	rest.AssertNil(campaign.DeleteByID(id))
	s.Success(ctx)
}

func (s *CampaignServer) handleGetCampaign(ctx *gin.Context) {
	var deviceID = ctx.Query("device_id")
	var channel = ctx.Query("channel")
	var at, _ = strconv.Atoi(ctx.Query("at"))
	var c, err = campaign.GetCampaignByDeviceOrChannel(deviceID, channel, int64(at))
	rest.AssertNil(err)
	s.SendData(ctx, c)
}
