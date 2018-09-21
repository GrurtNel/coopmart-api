package api

import (
	"feedback/api/campaign"
	"feedback/api/report"
	"feedback/api/setting"
	"feedback/api/survey"
	"feedback/api/vendor"
	"github.com/gin-gonic/gin"
)

func InitApi(root *gin.RouterGroup) {
	survey.NewSurveyServer(root, "survey")
	campaign.NewCampaignServer(root, "campaign")
	report.NewReportServer(root, "report")
	setting.NewSettingServer(root, "setting")
	vendor.NewCyberServer(root, "cyber")
}
