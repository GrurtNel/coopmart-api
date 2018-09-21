package result

import (
	"feedback/o/setting"
	"feedback/x/utils"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var StoreCache = map[string]*StoreMonitor{}

type StoreMonitor struct {
	ID       string
	Name     string
	Finish   int
	Unfinish int
	High     int `json:"high"`
	Credit   int `json:"credit"`
	Medium   int `json:"medium"`
	Low      int `json:"low"`
}

func GetQuantityReportRealtime(groupBy string) ([]*QuantityReport, error) {
	var now = utils.Now{
		Time: time.Now(),
	}
	var start = now.BeginningOfDay().Unix() * 1000
	var end = now.EndOfDay().Unix() * 1000
	var res []*QuantityReport
	var setting, _ = setting.GetSetting()
	var sumFunc = func(cond bson.M) bson.M {
		return bson.M{"$sum": cond}
	}
	countHigh := sumFunc(condQuantityReport(10.1, setting.HighRate))
	countCredit := sumFunc(condQuantityReport(setting.HighRate, setting.CreditRate))
	countMedium := sumFunc(condQuantityReport(setting.CreditRate, setting.MediumRate))
	countLow := sumFunc(condQuantityReport(setting.MediumRate, 0))
	countFinish := sumFunc(bson.M{"finished": true})
	countUnfinish := sumFunc(bson.M{"finished": false})
	count := bson.M{"$sum": 1}
	var match = matchAggregateDurationStore(int(start), int(end))
	var group = bson.M{
		"$group": bson.M{
			"_id":           "$" + groupBy,
			"store":         bson.M{"$first": "$store"},
			"device":        bson.M{"$first": "$device"},
			"high":          countHigh,
			"credit":        countCredit,
			"medium":        countMedium,
			"low":           countLow,
			"count":         count,
			"finished":      countFinish,
			"unfinished":    countUnfinish,
			"average_point": bson.M{"$avg": "$average_point"},
		},
	}
	var aggregate = []bson.M{match, group}
	if groupBy == "uname" {
		aggregate = append(aggregate, bson.M{"$sort": bson.M{"average_point": -1}})
	}
	err := ResultTable.Pipe(aggregate).All(&res)
	return res, err
}

func GetPoorFeedback(poorFilter bool) ([]*SurveyResult, error) {
	var now = utils.Now{
		Time: time.Now(),
	}
	var start = now.BeginningOfDay().Unix() * 1000
	var end = now.EndOfDay().Unix() * 1000
	var setting, _ = setting.GetSetting()
	var res []*SurveyResult
	var qualityFilter = bson.M{}
	if poorFilter {
		qualityFilter = bson.M{
			"average_point": bson.M{
				"$lte": setting.MediumRate,
			},
		}
	} else {
		qualityFilter = bson.M{
			"average_point": bson.M{
				"$gt": setting.MediumRate,
			},
		}
	}
	err := ResultTable.Find(bson.M{
		"$and": []bson.M{
			bson.M{"created_at": bson.M{"$lte": end}},
			bson.M{"created_at": bson.M{"$gte": start}},
			qualityFilter,
		},
	}).Sort("-created_at").All(&res)
	// .Limit(10)
	if res == nil {
		return []*SurveyResult{}, nil
	}
	return res, err
}
