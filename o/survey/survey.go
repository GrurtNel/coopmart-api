package survey

import (
	"feedback/x/db/mongodb"
	"feedback/x/rest"

	"github.com/golang/glog"

	"gopkg.in/mgo.v2/bson"
)

type Survey struct {
	mongodb.BaseModel `bson:",inline"`
	Name              string     `bson:"name" json:"name" validate:"required"`
	Questions         []Question `bson:"questions" json:"questions" validate:"required"`
	I18n              bool       `bson:"i18n" json:"i18n"`
	DeviceID          []string   `bson:"device_ids" json:"device_ids"`
	Preview           bool       `bson:"preview" json:"preview"`
}
type Question struct {
	ID          string     `bson:"id" json:"id"`
	Content     string     `bson:"content" json:"content" validate:"required"`
	I18nContent string     `bson:"i18n_content" json:"i18n_content"`
	Summary     string     `bson:"summary" json:"summary"`
	Type        SurveyType `bson:"type" json:"type" validate:"required"`
	Point       int        `bson:"point" json:"point" validate:"gte=0"`
	Manded      bool       `bson:"manded" json:"manded"`
	Answers     []Answer   `bson:"answers" json:"answers" validate:"required"`
	Link        int        `bson:"link" json:"link"`
}

//Answer -
type Answer struct {
	Icon        string `bson:"icon" json:"icon"`
	Content     string `bson:"content" json:"content" validate:"required"`
	I18nContent string `bson:"i18n_content" json:"i18n_content"`
	Point       int    `bson:"point" json:"point" validate:"gte=0"`
	Link        int    `bson:"link" json:"link"`
}
type SurveyType string

const (
	SINGLE   = SurveyType("single")
	MULTIPLE = SurveyType("multiple")
	ANSWER   = SurveyType("answer")
)

var SurveyTable = mongodb.NewTable("survey", "SUR", 5)

func (q Question) Validate() error {
	if q.Type != ANSWER {
		var answerPointSum = 0
		for _, item := range q.Answers {
			if item.Point < 0 {
				return rest.BadRequest("point of answer must be > 0")
			}
			answerPointSum += item.Point
		}
		if q.Type != SINGLE && q.Point < answerPointSum {
			return rest.BadRequest("point of question must be > sum of point's answers")
		}
	}
	return nil
}
func Validate(s *Survey) error {
	for _, item := range s.Questions {
		if err := item.Validate(); err != nil {
			return err
		}
	}
	return nil
}
func ListSurvey() ([]*Survey, error) {
	var surveys []*Survey
	err := SurveyTable.FindWhere(bson.M{"preview": false}, &surveys)
	return surveys, err
}

func (s *Survey) Create() error {
	rest.AssertNil(rest.Validate(s), Validate(s))
	return SurveyTable.Create(s)
}

func AddDeviceToSurvey(deviceID string, surveyID string) error {
	if s, err := GetSurveyByDevice(deviceID); err == nil {
		if s != nil {
			return rest.BadRequest("Device exists in another survey")
		}
	}
	return SurveyTable.UpdateId(surveyID, bson.M{
		"$addToSet": bson.M{
			"device_ids": deviceID,
		},
	})
}

func GetSurveyByDevice(feedbackID string) (*Survey, error) {
	var survey *Survey
	var err = SurveyTable.FindOne(bson.M{"device_ids": feedbackID}, &survey)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	return survey, nil
}

type DeviceSurvey struct {
	Name         string `json:"survey_name" bson:"survey_name"`
	FeedbackCode string `json:"feedback_code" bson:"feedback_code"`
}

func GetListDeviceSurvey() ([]*DeviceSurvey, error) {
	var ds []*DeviceSurvey
	var err = SurveyTable.Pipe([]bson.M{
		bson.M{
			"$unwind": "$device_ids",
		},
		bson.M{
			"$project": bson.M{
				"survey_name":   "$name",
				"feedback_code": "$device_ids",
			},
		},
	}).All(&ds)
	if err != nil {
		glog.Error(err)
		return nil, err
	}
	return ds, nil
}

type QuestionValidator struct {
}

var questionError string

func (v QuestionValidator) IsSatisfied(s interface{}) bool {
	if survey, ok := s.(*Survey); ok {
		if len(survey.Questions) == 0 {
			questionError = "Question at least 1"
			return false
		}
		for _, q := range survey.Questions {
			if q.Content == "" {
				questionError = "Content is not emplty"
				return false
			}
			if q.Type != SINGLE && q.Type != MULTIPLE && q.Type != ANSWER {
				questionError = "Type is a value of {single,multiple,answer}"
				return false
			}
			if q.Point < 0 {
				questionError = "Point must be >0"
			}
			if len(q.Answers) == 0 {
				if q.Type != ANSWER {
					questionError = "Answer at least 1"
					return false
				}
			} else {
				for _, item := range q.Answers {
					if item.Content == "" {
						questionError = "Content is not emplty"
						return false
					}
					if q.Point < 0 {
						questionError = "Point must be >0"
					}
				}
			}
		}
		return true
	}
	return false
}

func (v QuestionValidator) DefaultMessage() string {
	return questionError
}

func GetPreviewSurvey(surveyID string) (*Survey, error) {
	var res *Survey
	var err = SurveyTable.FindByID(surveyID, &res)
	return res, err
}
