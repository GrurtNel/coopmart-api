package survey

import (
	"feedback/o/luckyuser"
	"feedback/o/result"
	"feedback/o/setting"
	"feedback/o/survey"
	"feedback/room"
	"feedback/x/rest"
	"fmt"
	"g/x/web"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

var hub = room.NewHub()

type SurveyServer struct {
	*gin.RouterGroup
	rest.JsonRender
}

func NewSurveyServer(parent *gin.RouterGroup, name string) *SurveyServer {
	var s = SurveyServer{
		RouterGroup: parent.Group(name),
	}
	s.GET("/list", s.GetAll)
	s.GET("/test", s.handleTest)
	s.POST("/create", s.AddSurvey)
	s.POST("/update", s.handleUpdate)
	s.GET("/get", s.handleGetSurvey)
	s.GET("/delete", s.handleDelete)
	s.POST("/result/create", s.handleCreateResult)
	s.POST("/unfinish_result/create", s.handleCreateUnfinishResult)
	s.POST("/lucky_user/create", s.handleCreateLuckyUser)
	s.POST("/icon/upload", s.UploadIcon)
	s.GET("/icon/list", s.ListIcon)
	s.POST("/device/add", s.AddSurveyDevice)
	s.GET("/list-device", s.handleListDeviceService)
	s.GET("/device/:id", s.handleFeedbackDevice)
	go hub.Loop()
	return &s
}

func (s *SurveyServer) handleGetSurvey(ctx *gin.Context) {
	var id = ctx.Query("id")
	var srv, err = survey.GetPreviewSurvey(id)
	web.AssertNil(err)
	s.SendData(ctx, srv)
}

func (s *SurveyServer) handleCreateLuckyUser(ctx *gin.Context) {
	var luckyUser *luckyuser.LuckyUser
	web.AssertNil(ctx.BindJSON(&luckyUser))
	web.AssertNil(luckyUser.Create())
	s.Success(ctx)
}

func (s *SurveyServer) handleTest(ctx *gin.Context) {
	var count = result.CountFeedbackToday()
	var luckySetting, _ = setting.GetLuckySetting()
	if luckySetting != nil {
		if count == luckySetting.LuckyNumber-1 {
			s.SendData(ctx, luckySetting)
		}
	}
	s.Success(ctx)
}

func (s *SurveyServer) AddSurvey(ctx *gin.Context) {
	var srv *survey.Survey
	web.AssertNil(ctx.BindJSON(&srv))
	web.AssertNil(srv.Create())
	s.SendData(ctx, srv)
}

type LuckyResponse struct {
	*setting.LuckySetting
	GiftCode int `json:"gift_code"`
}

func getLuckyProgram() *LuckyResponse {
	var count = result.CountFeedbackToday()
	var luckySetting, _ = setting.GetLuckySetting()
	if luckySetting != nil {
		if (count%luckySetting.LuckyNumber == 0) && luckySetting.Activated {
			return &LuckyResponse{
				LuckySetting: luckySetting,
				GiftCode:     count,
			}
		}
	}
	return nil
}

func (s *SurveyServer) handleCreateResult(ctx *gin.Context) {
	var res *result.SurveyResult
	rest.AssertNil(ctx.BindJSON(&res), res.Create())
	hub.Transporter <- res
	result.ConvertToSurveyAggregate(res)
	result.ConvertToCampaignAggregate(res)
	if luckySetting := getLuckyProgram(); luckySetting != nil {
		s.SendData(ctx, luckySetting)
	} else {
		s.SendData(ctx, res)
	}
}

func (s *SurveyServer) handleCreateUnfinishResult(ctx *gin.Context) {
	var res *result.SurveyResult
	rest.AssertNil(ctx.BindJSON(&res), res.CreateUnfinish())
	s.SendData(ctx, res)
}

func (s *SurveyServer) handleUpdate(ctx *gin.Context) {
	var srv *survey.Survey
	web.AssertNil(ctx.BindJSON(&srv))
	web.AssertNil(survey.Update(srv))
	s.SendData(ctx, srv)
}
func (s *SurveyServer) handleDelete(ctx *gin.Context) {
	var id = ctx.Query("id")
	rest.AssertNil(survey.DeleteByID(id))
	s.Success(ctx)
}

func (s *SurveyServer) handleListDeviceService(ctx *gin.Context) {
	var srv, err = survey.GetListDeviceSurvey()
	rest.AssertNil(err)
	s.SendData(ctx, srv)
}

func (s *SurveyServer) handleFeedbackDevice(ctx *gin.Context) {
	var id = ctx.Param("id")
	var srv, err = survey.GetSurveyByDevice(id)
	rest.AssertNil(err)
	s.SendData(ctx, srv)
}

func (s *SurveyServer) ListIcon(ctx *gin.Context) {
	var listIcon = []string{}
	files, err := ioutil.ReadDir("./static/smiley")
	rest.AssertNil(err)
	fmt.Println(ctx.Request.Host)
	for _, item := range files {
		listIcon = append(listIcon, "http://"+ctx.Request.Host+"/static/smiley/"+item.Name())
	}
	s.SendData(ctx, listIcon)
}

func (s *SurveyServer) AddSurveyDevice(ctx *gin.Context) {
	var surveyDevice = struct {
		SurveyID string `json:"survey_id"`
		DeviceID string `json:"device_id"`
	}{}
	web.AssertNil(ctx.BindJSON(&surveyDevice))
	web.AssertNil(survey.AddDeviceToSurvey(surveyDevice.DeviceID, surveyDevice.SurveyID))
	s.Success(ctx)
}

func (s *SurveyServer) GetAll(ctx *gin.Context) {
	surveys, _ := survey.ListSurvey()
	s.SendData(ctx, surveys)
}

func (s *SurveyServer) UploadIcon(ctx *gin.Context) {
	var file, err = ctx.FormFile("icon")
	web.AssertNil(err)
	web.AssertNil(ctx.SaveUploadedFile(file, "./static/smiley/abc"))
	s.Success(ctx)
}
