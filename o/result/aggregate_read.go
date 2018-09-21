package result

import (
	"feedback/x/utils"
	"fmt"
	"github.com/golang/glog"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

type Report struct {
	ID            string  `bson:"_id" json:"actor_name"`
	AveragePoint  float32 `bson:"average_point" json:"average_point"`
	High          int     `bson:"high" json:"high"`
	Credit        int     `bson:"credit" json:"credit"`
	Medium        int     `bson:"medium" json:"medium"`
	Low           int     `bson:"low" json:"low"`
	Count         int     `bson:"count" json:"count"`
	HighPercent   string  `bson:"-" json:"high_percent"`
	CreditPercent string  `bson:"-" json:"credit_percent"`
	MediumPercent string  `bson:"-" json:"medium_percent"`
	LowPercent    string  `bson:"-" json:"low_percent"`
	Channel       string  `bson:"channel" json:"channel"`
}
type GeneralReport struct {
	Report `bson:",inline"`
}

func GetGeneralReport(start, end int) (*GeneralReport, error) {
	var gReports *GeneralReport
	var match = matchAggregateDuration(start, end)
	var group = bson.M{
		"$group": bson.M{
			"_id": nil,
			"average_point": bson.M{
				"$avg": "$average_point",
			},
			"high":   groupByLevel("$high"),
			"credit": groupByLevel("$credit"),
			"medium": groupByLevel("$medium"),
			"low":    groupByLevel("$low"),
			"count":  bson.M{"$sum": 1},
		},
	}
	err := SurveyAggregateTable.Pipe([]bson.M{
		match,
		group,
	}).One(&gReports)
	if gReports != nil {
		gReports.TransformData()
	}
	return gReports, err
}

func GetGeneralReportBy(by string, start, end int) ([]*GeneralReport, error) {
	var reports []*GeneralReport
	var match = matchAggregateDuration(start, end)
	var group = bson.M{
		"$group": bson.M{
			"_id": "$" + by,
			"average_point": bson.M{
				"$avg": "$average_point",
			},
			"high":   groupByLevel("$high"),
			"credit": groupByLevel("$credit"),
			"medium": groupByLevel("$medium"),
			"low":    groupByLevel("$low"),
			"count":  bson.M{"$sum": 1},
		},
	}
	err := SurveyAggregateTable.Pipe([]bson.M{
		match,
		group,
	}).All(&reports)
	if reports != nil {
		for _, item := range reports {
			item.TransformData()
		}
	}
	return reports, err
}

type CampaignReport struct {
	CampaignName string   `bson:"campaign_name"  json:"campaign_name" `
	Channels     []string `bson:"channels"  json:"channels" `
	Count        int      `bson:"count" json:"count"`
	AveragePoint float32  `bson:"average_point" json:"average_point"`
	Start        int64    `bson:"start"  json:"-" `
	End          int64    `bson:"end" json:"-"`
	Duration     string   `bson:"-" json:"duration"`
	Surveys      []struct {
		ID   string `bson:"_id" json:"id"`
		Name string `bson:"name" json:"name"`
	} `bson:"surveys" json:"surveys"`
}

func GetCampaignReport(start, end int) ([]*CampaignReport, error) {
	var match = matchAggregateDuration(start, end)
	var result []*CampaignReport
	var group = bson.M{
		"$group": bson.M{
			"_id": "$campaign_id",
			"campaign_name": bson.M{
				"$first": "$campaign",
			},
			"count": bson.M{
				"$sum": 1,
			},
			"average_point": bson.M{
				"$avg": "$average_point",
			},
		},
	}
	var joinCampaign = bson.M{
		"$lookup": bson.M{
			"from":         "campaign",
			"localField":   "_id",
			"foreignField": "_id",
			"as":           "campaign",
		},
	}
	var joinSurvey = bson.M{
		"$lookup": bson.M{
			"from":         "survey",
			"localField":   "campaign.surveys",
			"foreignField": "_id",
			"as":           "surveys",
		},
	}
	var project = bson.M{
		"$project": bson.M{
			"campaign_name": "$campaign_name",
			"count":         "$count",
			"average_point": "$average_point",
			"start":         bson.M{"$arrayElemAt": []interface{}{"$campaign.start", 0}},
			"end":           bson.M{"$arrayElemAt": []interface{}{"$campaign.end", 0}},
			"surveys":       "$surveys",
			"channels":      bson.M{"$arrayElemAt": []interface{}{"$campaign.channels", 0}},
		},
	}
	err := ResultTable.Pipe([]bson.M{match, group, joinCampaign, joinSurvey, project}).All(&result)
	if result != nil {
		for _, item := range result {
			item.TransformData()
		}
	}
	return result, err
}

type History struct {
	CreatedAt string
}

func GetHistory(start, end int) ([]*SurveyResult, error) {
	var histories []*SurveyResult
	var match = matchAggregateDuration(start, end)
	var unwind = bson.M{
		"$unwind": "$survey_detail",
	}

	var project = bson.M{
		"$project": bson.M{
			"created_at":    "$created_at",
			"uname":         "$uname",
			"store":         "$store",
			"device":        "$device",
			"channel":       "$channel",
			"location":      "$location",
			"point":         bson.M{"$sum": "$survey_detail.feedback_detail.point"},
			"max_point":     bson.M{"$sum": "$survey_detail.feedback_detail.max_point"},
			"survey_detail": "$survey_detail",
		},
	}

	var group = bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"created_at": "$created_at",
				"uname":      "$uname",
				"store":      "$store",
				"device":     "$device",
				"channel":    "$channel",
				"location":   "$location",
			},
			"survey_detail": bson.M{"$addToSet": "$survey_detail"},
			"point":         bson.M{"$sum": "$point"},
			"max_point":     bson.M{"$sum": "$max_point"},
		},
	}
	var cond = bson.M{
		"$cond": []interface{}{
			bson.M{"$eq": []interface{}{"$max_point", 0}},
			"0",
			bson.M{"$divide": []string{"$point", "$max_point"}},
		},
	}
	var project2 = bson.M{
		"$project": bson.M{
			"_id":           0,
			"created_at":    "$_id.created_at",
			"uname":         "$_id.uname",
			"store":         "$_id.store",
			"device":        "$_id.device",
			"channel":       "$_id.channel",
			"location":      "$_id.location",
			"point":         bson.M{"$sum": "$point"},
			"max_point":     bson.M{"$sum": "$max_point"},
			"average_point": cond,
			"survey_detail": "$survey_detail",
		},
	}
	var orderBy = bson.M{
		"$sort": bson.M{
			"_id.created_at": -1,
		},
	}
	err := ResultTable.Pipe([]bson.M{match, unwind, project, group, orderBy, project2}).All(&histories)
	if err != nil {
		glog.Error(err)
	}
	fmt.Println(len(histories))
	if histories != nil {
		for _, item := range histories {
			item.TransformData()
		}
	}
	return histories, err
}

func (r *GeneralReport) TransformData() {
	r.HighPercent = utils.Float32ToPercentString(float32(r.High) / float32(r.Count))
	r.CreditPercent = utils.Float32ToPercentString(float32(r.Credit) / float32(r.Count))
	r.MediumPercent = utils.Float32ToPercentString(float32(r.Medium) / float32(r.Count))
	r.LowPercent = utils.Float32ToPercentString(float32(r.Low) / float32(r.Count))
	r.AveragePoint = utils.FormatFloatPoint(r.AveragePoint)
}

func (r *CampaignReport) TransformData() {
	r.AveragePoint = utils.FormatFloatPoint(r.AveragePoint)
	var split = func(input int64) string {
		return strings.Split(utils.UnixToDate(input), " ")[0]
	}
	var start = split(r.Start)
	var end = split(r.End)
	r.Duration = start + " - " + end
}

func (r *SurveyResult) TransformData() {
	r.AveragePoint = utils.FormatDecimalPoint(r.AveragePoint)
	var temp = strings.Split(utils.UnixToDate(r.CreatedAt), " ")
	r.DayCtime = temp[0]
	r.HourCtime = temp[1]
}
