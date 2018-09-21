package setting

import (
	"feedback/o/result"
	"feedback/o/setting"
	"feedback/x/rest"
	"g/x/web"

	"github.com/gin-gonic/gin"
)

type SettingServer struct {
	*gin.RouterGroup
	rest.JsonRender
}

func NewSettingServer(parent *gin.RouterGroup, name string) *SettingServer {
	var s = SettingServer{
		RouterGroup: parent.Group(name),
	}
	s.POST("/create", s.handleCreateSetting)
	s.GET("/get", s.handleGetSetting)
	s.POST("/lucky/create", s.handleCreateLuckySetting)
	s.GET("/lucky/get", s.handleGetLuckySetting)
	s.POST("/email/create", s.handleSetEmailStore)
	s.GET("/email/list", s.handleListEmail)
	s.GET("/email/delete", s.handleDeleteEmail)
	s.POST("/upload/logo", s.handleUploadLogo)
	s.POST("/upload/video", s.handleUploadVideo)
	s.POST("/upload/background", s.handleUploadBackground)
	return &s
}

func (s *SettingServer) handleSetEmailStore(ctx *gin.Context) {
	var storeEmail *setting.StoreEmail
	ctx.BindJSON(&storeEmail)
	web.AssertNil(storeEmail.Create())
	s.Success(ctx)
}
func (s *SettingServer) handleListEmail(ctx *gin.Context) {
	var email = setting.GetEmails()
	s.SendData(ctx, email)
}

func (s *SettingServer) handleDeleteEmail(ctx *gin.Context) {
	err := setting.DeleteByID(ctx.Query("id"))
	web.AssertNil(err)
	s.Success(ctx)
}

func (s *SettingServer) handleTestSocket(ctx *gin.Context) {
	var res, _ = result.GetQuantityReportRealtime("store")
	s.SendData(ctx, res)
}

func (s *SettingServer) handleUploadBackground(ctx *gin.Context) {
	var file, err = ctx.FormFile("background")
	web.AssertNil(err)
	web.AssertNil(ctx.SaveUploadedFile(file, "./static/setting/background"))
	s.Success(ctx)
}

func (s *SettingServer) handleUploadLogo(ctx *gin.Context) {
	var file, err = ctx.FormFile("logo")
	web.AssertNil(err)
	web.AssertNil(ctx.SaveUploadedFile(file, "./static/setting/logo"))
	s.Success(ctx)
}

func (s *SettingServer) handleUploadVideo(ctx *gin.Context) {
	var file, err = ctx.FormFile("adv")
	web.AssertNil(err)
	web.AssertNil(ctx.SaveUploadedFile(file, "./static/adv/adv.mp4"))
	s.Success(ctx)
}

func (s *SettingServer) handleCreateSetting(ctx *gin.Context) {
	var set *setting.Setting
	ctx.BindJSON(&set)
	web.AssertNil(set.Create())
	s.Success(ctx)
}

func (s *SettingServer) handleCreateLuckySetting(ctx *gin.Context) {
	var set *setting.LuckySetting
	ctx.BindJSON(&set)
	web.AssertNil(set.Create())
	s.Success(ctx)
}

func (s *SettingServer) handleGetLuckySetting(ctx *gin.Context) {
	var set, err = setting.GetLuckySetting()
	rest.AssertNil(err)
	s.SendData(ctx, set)
}

func (s *SettingServer) handleGetSetting(ctx *gin.Context) {
	var set, err = setting.GetSetting()
	rest.AssertNil(err)
	s.SendData(ctx, set)
}
