package result

import (
	"feedback/o/setting"
	"gopkg.in/mgo.v2/bson"
)

type CommonReport struct {
	AveragePoint float32 `bson:"average_point" json:"average_point"`
	High         int     `bson:"high"  json:"high" `
	Credit       int     `bson:"credit" json:"credit"`
	Medium       int     `bson:"medium" json:"medium"`
	Low          int     `bson:"low" json:"low"`
	Count        int     `bson:"count" json:"count"`
}
type GeneralChannelReport struct {
	Actor        string `bson:"channel"  json:"channel" `
	CommonReport `bson:",inline"`
}

func GetGeneralChannelReport(start, end int) ([]*GeneralChannelReport, error) {
	var res []*GeneralChannelReport
	var group = bson.M{
		"$group": groupGeneralReport("channel"),
	}
	var match = matchAggregateDuration(start, end)
	err := CampaignAggregateTable.Pipe([]bson.M{match, group}).All(&res)
	return res, err

}
func groupGeneralReport(groupBy string) bson.M {
	var setting, _ = setting.GetSetting()
	var sumFunc = func(cond bson.M) bson.M {
		return bson.M{"$sum": cond}
	}
	countHigh := sumFunc(condQuantityReport(10.1, setting.HighRate))
	countCredit := sumFunc(condQuantityReport(setting.HighRate, setting.CreditRate))
	countMedium := sumFunc(condQuantityReport(setting.CreditRate, setting.MediumRate))
	countLow := sumFunc(condQuantityReport(setting.MediumRate, -0.1))
	count := bson.M{"$sum": 1}
	return bson.M{
		"_id":           "$" + groupBy,
		groupBy:         bson.M{"$first": "$" + groupBy},
		"high":          countHigh,
		"credit":        countCredit,
		"medium":        countMedium,
		"low":           countLow,
		"count":         count,
		"average_point": bson.M{"$avg": "$average_point"},
	}
}
