package result

import (
	"feedback/o/setting"
	"feedback/x/db/mongodb"
	"feedback/x/utils"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

type SurveyResult struct {
	mongodb.BaseModel `bson:",inline"`
	Uname             string          `bson:"uname" json:"uname"`
	Device            string          `bson:"device" json:"device"`
	Store             string          `bson:"store" json:"store"`
	StoreCode         string          `bson:"store_code" json:"store_code"`
	Counter           string          `bson:"counter" json:"counter"`
	Campaign          string          `bson:"campaign" json:"campaign"`
	CampaignID        string          `bson:"campaign_id" json:"campaign_id"`
	Question          string          `bson:"question" json:"question"`
	SurveyDetails     []*SurveyDetail `bson:"survey_detail" json:"survey_detail"`
	Point             int             `bson:"point" json:"point"`
	MaxPoint          int             `bson:"max_point" json:"max_point"`
	AveragePoint      float32         `bson:"average_point" json:"average_point"`
	DayCtime          string          `bson:"day_ctime" json:"day_ctime"`
	Location          string          `bson:"location" json:"location"`
	HourCtime         string          `bson:"hour_ctime" json:"hour_ctime"`
	CTime             string          `bson:"ctime" json:"ctime"`
	Channel           string          `bson:"channel" json:"channel"`
	Finished          bool            `bson:"finished" json:"finished"`
}

type SurveyDetail struct {
	SurveyID        string            `bson:"survey_id" json:"survey_id"`
	SurveyName      string            `bson:"survey_name" json:"survey_name"`
	Point           int               `bson:"point" json:"point"`
	MaxPoint        int               `bson:"max_point" json:"max_point"`
	FeedbackDetails []*FeedbackDetail `bson:"feedback_detail" json:"feedback_detail"`
}

type FeedbackDetail struct {
	Content  string `bson:"content" json:"content"`
	Answer   string `bson:"answer" json:"answer"`
	Point    int    `bson:"point" json:"point"`
	MaxPoint int    `bson:"max_point" json:"max_point"`
	Type     string `bson:"type" json:"type"`
}

// ResultTable đasadsa
var ResultTable = mongodb.NewTable("result", "RES", 5)
var UnfinishResultTable = mongodb.NewTable("unfinish_result", "RES", 5)

// Create create
func (s *SurveyResult) Create() error {
	var temp = strings.Split(utils.UnixToDate(time.Now().Unix()*1000), " ")
	s.CTime = temp[0]
	s.StoreCode = strings.Split(s.Device, "_")[0]
	for _, survey := range s.SurveyDetails {
		survey.caculatePoint()
	}
	s.caculatePoint()
	go func() {
		var email = setting.GetEmailByStore(s.StoreCode, s.Channel)
		glog.Infof("[email]send email to %s store %s", email, s.StoreCode)
		s.sendLowResult(email)
	}()
	return ResultTable.Create(s)
}

// CreateUnfinish s
func (s *SurveyResult) CreateUnfinish() error {
	var temp = strings.Split(utils.UnixToDate(time.Now().Unix()*1000), " ")
	s.CTime = temp[0]
	s.StoreCode = strings.Split(s.Device, "_")[0]
	for _, survey := range s.SurveyDetails {
		survey.caculatePoint()
	}
	s.caculatePoint()
	return UnfinishResultTable.Create(s)
}

const (
	HIGH_RATE   = 0.85
	CREDIT_RATE = 0.65
	MEDIUM_RATE = 0.5
)

func (s *SurveyResult) caculatePoint() {
	var point = 0
	var maxPoint = 0
	for _, survey := range s.SurveyDetails {
		point += survey.Point
		maxPoint += survey.MaxPoint
	}
	s.Point = point
	s.MaxPoint = maxPoint
	s.AveragePoint = float32(s.Point) / float32(s.MaxPoint) * 10
}

func (s *SurveyDetail) caculatePoint() {
	var point = 0
	var maxPoint = 0
	for _, feedback := range s.FeedbackDetails {
		point += feedback.Point
		maxPoint += feedback.MaxPoint
	}
	s.Point = point
	s.MaxPoint = maxPoint
}

func CountFeedbackToday() int {
	var now = utils.Now{
		Time: time.Now(),
	}
	var startOfDay = now.BeginningOfDay().Unix()
	var count = 0
	count, _ = ResultTable.Find(bson.M{"created_at": bson.M{
		"$gte": startOfDay * 1000,
	}}).Count()
	return count
}
