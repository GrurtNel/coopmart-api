package report

import (
	"feedback/o/luckyuser"
	"feedback/o/result"
	"feedback/x/rest"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ReportServer struct {
	*gin.RouterGroup
	rest.JsonRender
}

func NewReportServer(parent *gin.RouterGroup, name string) *ReportServer {
	var s = ReportServer{
		RouterGroup: parent.Group(name),
	}
	s.GET("/general", s.handleGeneralReport)
	s.GET("/activity-frequency", s.handleActivityFrequency)
	s.GET("/general/channel", s.handleGeneralChannelReport)
	s.GET("/aggregate", s.handleAggregate)
	s.GET("/campaign", s.handleCampaignReport)
	s.GET("/history", s.handleGetHistory)
	s.GET("/chart/feedback-quantity", s.handleFeedbackQuantity)
	s.GET("/survey-analyst", s.handleSurveyAnalyst)
	s.GET("/recent-feedback", s.handleRecentFeedback)
	s.GET("/poor-feedback", s.handleGetPoorFeedback)
	s.GET("/test", s.handleTest)
	s.GET("/lucky-user/list", s.handleGetLuckyUsers)
	s.GET("/timeline", s.handleTimeline)
	return &s
}

func (s *ReportServer) handleTimeline(ctx *gin.Context) {
	var splitComma = func(param string) []string {
		var result = strings.Split(param, ",")
		if len(result) == 1 && result[0] == "" {
			return []string{}
		}
		return result
	}
	// var stores = ctx.QueryArray("stores")
	var storesParam = ctx.Query("stores")
	var unamesParam = ctx.Query("uname")
	var deviceParam = ctx.Query("devices")
	var channelParam = ctx.Query("channel")
	var start, _ = strconv.Atoi(ctx.Query("start"))
	var end, _ = strconv.Atoi(ctx.Query("end"))
	res := result.GetTimelineReport(start, end, splitComma(storesParam), splitComma(unamesParam), splitComma(deviceParam), splitComma(channelParam))
	s.SendData(ctx, res)
}

func (s *ReportServer) handleTest(ctx *gin.Context) {
	res, _ := result.GetPoorFeedback(false)
	s.SendData(ctx, res)
}

func (s *ReportServer) handleGetLuckyUsers(ctx *gin.Context) {
	res := luckyuser.GetLuckyUsers()
	s.SendData(ctx, res)
}

func (s *ReportServer) handleSurveyAnalyst(ctx *gin.Context) {
	var surveyID = ctx.Query("survey_id")
	res, err := result.GetSurveyAnalyst(surveyID)
	rest.AssertNil(err)
	s.SendData(ctx, res)
}

func (s *ReportServer) handleActivityFrequency(ctx *gin.Context) {
	res, err := result.GetActivityFrequencyChart(ctx.Query)
	rest.AssertNil(err)
	s.SendData(ctx, res)
}

func (s *ReportServer) handleGeneralChannelReport(ctx *gin.Context) {
	var start, _ = strconv.Atoi(ctx.Query("start"))
	var end, _ = strconv.Atoi(ctx.Query("end"))
	res, err := result.GetGeneralChannelReport(start, end)
	rest.AssertNil(err)
	s.SendData(ctx, res)
}

func (s *ReportServer) handleAggregate(ctx *gin.Context) {
	var start, _ = strconv.Atoi(ctx.Query("start"))
	var end, _ = strconv.Atoi(ctx.Query("end"))
	var group = ctx.Query("group")
	res, err := result.GetQuantityReport(start, end, group)
	rest.AssertNil(err)
	s.SendData(ctx, res)
}

func (s *ReportServer) handleFeedbackQuantity(ctx *gin.Context) {
	var start, _ = strconv.Atoi(ctx.Query("start"))
	var end, _ = strconv.Atoi(ctx.Query("end"))
	var campaignID = ctx.Query("campaign_id")
	res, err := result.GetQuantityChart(start, end, campaignID)
	rest.AssertNil(err)
	s.SendData(ctx, res)
}

func (s *ReportServer) handleGetHistory(ctx *gin.Context) {
	var start, _ = strconv.Atoi(ctx.Query("start"))
	var end, _ = strconv.Atoi(ctx.Query("end"))
	res, err := result.GetHistory(start, end)
	rest.AssertNil(err)
	s.SendData(ctx, res)
}

func (s *ReportServer) handleCampaignReport(ctx *gin.Context) {
	var start, _ = strconv.Atoi(ctx.Query("start"))
	var end, _ = strconv.Atoi(ctx.Query("end"))
	res, err := result.GetCampaignReport(start, end)
	rest.AssertNil(err)
	s.SendData(ctx, res)
}

func (s *ReportServer) handleGeneralReport(ctx *gin.Context) {
	var start, _ = strconv.Atoi(ctx.Query("start"))
	var end, _ = strconv.Atoi(ctx.Query("end"))
	var by = ctx.Query("by")
	if by == "" {
		res, err := result.GetGeneralReport(start, end)
		rest.AssertNil(err)
		s.SendData(ctx, res)
	} else {
		res, err := result.GetGeneralReportBy(by, start, end)
		rest.AssertNil(err)
		s.SendData(ctx, res)
	}
}

func (s *ReportServer) handleRecentFeedback(ctx *gin.Context) {
	res, err := result.GetRecentFeedback()
	rest.AssertNil(err)
	s.SendData(ctx, res)
}

func (s *ReportServer) handleGetPoorFeedback(ctx *gin.Context) {
	res, err := result.GetPoorFeedback(true)
	rest.AssertNil(err)
	s.SendData(ctx, res)
}
