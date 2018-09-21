package result

import (
	"feedback/o/setting"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

type SurveyDetailReport struct {
	Content  string                 `json:"content" bson:"content"`
	Type     string                 `json:"type" bson:"type"`
	Answers  []Answer               `json:"answers" bson:"answers"`
	Results  map[string]interface{} `json:"results" bson:"-"`
	Point    int                    `json:"point" bson:"point"`
	MaxPoint int                    `json:"max_point" bson:"max_point"`
	Count    int                    `json:"count" bson:"count"`
	High     int                    `bson:"high"  json:"high" `
	Credit   int                    `bson:"credit" json:"credit"`
	Medium   int                    `bson:"medium" json:"medium"`
	Low      int                    `bson:"low" json:"low"`
}
type Answer struct {
	Content  string `json:"content" bson:"content"`
	Store    string `json:"store" bson:"store"`
	Channel  string `json:"channel" bson:"channel"`
	Location string `json:"location" bson:"location"`
	Ctime    int64  `json:"ctime" bson:"ctime"`
}

// db.getCollection("result").aggregate([
// 	{
// 		$unwind:"$survey_detail"
// 	},
// 	{
// 		$match:{
// 			"survey_detail.survey_id":"SUR_vlfTu"
// 		}
// 	},
// 	{
// 		$unwind:"$survey_detail.feedback_detail"
// 	},
// 	{
// 		$group:{
// 			_id:{
// 			   content:"$survey_detail.feedback_detail.content",
// 			},
// 			content:{$first:"$survey_detail.feedback_detail.content"},
// 			type:{$first:"$survey_detail.feedback_detail.type"},
// 			answers:{$push:"$survey_detail.feedback_detail.answer"},
// 			answers:{$push:{answer:"$survey_detail.feedback_detail.answer",store:"$store",channel:"$channel",point:"$survey_detail.feedback_detail.point",max_point:"$survey_detail.feedback_detail.max_point"}}
// 		}
// 	},
// 	])
func GetSurveyAnalyst(surveyID string) ([]*SurveyDetailReport, error) {
	var result []*SurveyDetailReport
	var unwindSurvey = unwind("$survey_detail")
	var unwindFeedback = unwind("$survey_detail.feedback_detail")
	var match = getState(MATCH_STATE, []interface{}{"survey_detail.survey_id", surveyID})

	var setting, _ = setting.GetSetting()
	var sumFunc = func(cond bson.M) bson.M {
		return bson.M{"$sum": cond}
	}
	countHigh := sumFunc(condQuantityReport(10.1, setting.HighRate))
	countCredit := sumFunc(condQuantityReport(setting.HighRate, setting.CreditRate))
	countMedium := sumFunc(condQuantityReport(setting.CreditRate, setting.MediumRate))
	countLow := sumFunc(condQuantityReport(setting.MediumRate, 0))

	var group = bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"content": "$survey_detail.feedback_detail.content",
			},
			"content": bson.M{"$first": "$survey_detail.feedback_detail.content"},
			"type":    bson.M{"$first": "$survey_detail.feedback_detail.type"},
			"answers": bson.M{
				"$push": bson.M{
					"content":  "$survey_detail.feedback_detail.answer",
					"store":    "$store",
					"channel":  "$channel",
					"location": "$location",
					"ctime":    "$created_at",
				},
			},
			"point":         bson.M{"$sum": "$survey_detail.feedback_detail.point"},
			"max_point":     bson.M{"$sum": "$survey_detail.feedback_detail.max_point"},
			"count":         bson.M{"$sum": 1},
			"high":          countHigh,
			"credit":        countCredit,
			"medium":        countMedium,
			"low":           countLow,
			"average_point": bson.M{"$avg": "$average_point"},
		},
	}

	err := ResultTable.Pipe([]bson.M{unwindSurvey, match, unwindFeedback, group}).All(&result)
	if len(result) > 0 {
		for _, item := range result {
			item.TransformSurveyDetailReport()
		}
	}
	return result, err
}

func (s *SurveyDetailReport) TransformSurveyDetailReport() {
	s.Results = map[string]interface{}{}
	if s.Type == "single" {
		for _, item := range s.Answers {
			if s.Results[item.Content] != nil {
				s.Results[item.Content] = s.Results[item.Content].(int) + 1
			} else {
				s.Results[item.Content] = 1
			}
		}
		s.Answers = []Answer{}
	} else if s.Type == "multiple" {
		for _, item := range s.Answers {
			for _, item := range strings.Split(item.Content, ",") {
				if s.Results[item] != nil {
					s.Results[item] = s.Results[item].(int) + 1
				} else {
					s.Results[item] = 1
				}
			}
		}
		s.Answers = []Answer{}
	}
}

func GetRecentFeedback() ([]*SurveyResult, error) {
	var res []*SurveyResult
	err := ResultTable.Find(bson.M{}).Limit(10).Sort("-ctime").All(&res)
	return res, err
}
