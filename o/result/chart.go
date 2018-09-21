package result

import (
	"feedback/o/setting"
	"feedback/x/rest"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

type QuantityChart struct {
	Day      string `bson:"day" json:"day"`
	Quantity int    `bson:"quantity" json:"quantity"`
}

func GetQuantityChart(start, end int, campaignID string) ([]*QuantityChart, error) {
	var res []*QuantityChart
	var match = matchAggregateDuration(start, end)
	var match2 = bson.M{"$match": bson.M{}}
	if campaignID != "" {
		match2["$match"] = bson.M{
			"campaign_id": campaignID,
		}
	}
	var group = bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"ctime": "$ctime",
			},
			"quantity": bson.M{
				"$sum": 1,
			},
		},
	}
	var project = bson.M{
		"$project": bson.M{
			"day":      "$_id.ctime",
			"quantity": "$quantity",
		},
	}
	var orderBy = bson.M{"$sort": bson.M{"_id.ctime": -1}}
	err := ResultTable.Pipe([]bson.M{match, match2, group, orderBy, project}).All(&res)
	if err != nil {
		if rest.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	sortQuantityChartByTime(res)
	return res, nil
}

func mmDDYYYYToUnix(date string) int64 {
	var temp = strings.Split(date, "-")
	t, _ := time.Parse(time.RFC3339, temp[2]+"-"+temp[1]+"-"+temp[0]+"T00:00:00Z")
	return t.Unix()
}
func sortQuantityChartByTime(quantityChart []*QuantityChart) {
	for i := 0; i < len(quantityChart)-1; i++ {
		for j := i + 1; j < len(quantityChart); j++ {
			if mmDDYYYYToUnix(quantityChart[i].Day) > mmDDYYYYToUnix(quantityChart[j].Day) {
				var mid = quantityChart[i]
				quantityChart[i] = quantityChart[j]
				quantityChart[j] = mid
			}
		}
	}
}

type QuantityReport struct {
	Actor        string  `bson:"_id"  json:"actor" `
	Store        string  `bson:"store"  json:"store" `
	AveragePoint float32 `bson:"average_point" json:"average_point"`
	High         int     `bson:"high"  json:"high" `
	Credit       int     `bson:"credit" json:"credit"`
	Medium       int     `bson:"medium" json:"medium"`
	Low          int     `bson:"low" json:"low"`
	Count        int     `bson:"count" json:"count"`
	Finished     int     `bson:"finished" json:"finished"`
	Unfinished   int     `bson:"unfinished" json:"unfinished"`
}

func GetQuantityReport(start, end int, groupBy string) ([]*QuantityReport, error) {
	var res []*QuantityReport
	var setting, _ = setting.GetSetting()
	var sumFunc = func(cond bson.M) bson.M {
		return bson.M{"$sum": cond}
	}
	countHigh := sumFunc(condQuantityReport(10.1, setting.HighRate))
	countCredit := sumFunc(condQuantityReport(setting.HighRate, setting.CreditRate))
	countMedium := sumFunc(condQuantityReport(setting.CreditRate, setting.MediumRate))
	countLow := sumFunc(condQuantityReport(setting.MediumRate, 0))
	count := bson.M{"$sum": 1}
	var match = matchAggregateDurationStore(start, end)
	var group = bson.M{
		"$group": bson.M{
			"_id":           "$" + groupBy,
			"store":         bson.M{"$first": "$store"},
			"high":          countHigh,
			"credit":        countCredit,
			"medium":        countMedium,
			"low":           countLow,
			"count":         count,
			"average_point": bson.M{"$avg": "$average_point"},
		},
	}
	err := ResultTable.Pipe([]bson.M{match, group}).All(&res)
	return res, err
}

func condQuantityReport(above, under float32) bson.M {
	var andCond = bson.M{
		"$and": []bson.M{
			bson.M{
				"$gte": []interface{}{"$average_point", under},
			},
			bson.M{
				"$lt": []interface{}{"$average_point", above},
			},
		},
	}
	return bson.M{
		"$cond": []interface{}{andCond, 1, 0},
	}
}
