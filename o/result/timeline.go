package result

import (
	"feedback/o/setting"
	"fmt"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const CHANNELWEB = "website"
const CHANNELSTORE = "store"

type Timeline struct {
	Time   string `bson:"time"  json:"time" `
	High   int    `bson:"high"  json:"high" `
	Credit int    `bson:"credit" json:"credit"`
	Medium int    `bson:"medium" json:"medium"`
	Low    int    `bson:"low" json:"low"`
	Count  int    `bson:"count" json:"count"`
}

func containArray(arr []string, search string) bool {
	for _, item := range arr {
		if item == search {
			return true
		}
	}
	return false
}
func processTimeline(start, end int, stores, unames, devices, channels []string) TimelineReport {
	var surveyResult []*SurveyResult
	var queryFilter = bson.M{}

	if len(unames) > 0 {
		queryFilter = bson.M{
			"uname": bson.M{
				"$in": unames,
			},
		}
	} else if len(devices) > 0 {
		queryFilter = bson.M{
			"device": bson.M{
				"$in": devices,
			},
		}
	} else if len(stores) > 0 {
		queryFilter = bson.M{
			"store_code": bson.M{
				"$in": stores,
			},
		}
	} else if len(channels) > 0 {
		if containArray(channels, CHANNELWEB) {
			queryFilter = bson.M{
				"channel": bson.M{
					"$in": channels,
				},
			}
		}
		if containArray(channels, CHANNELSTORE) && len(stores) > 0 {
			queryFilter["store_code"] = bson.M{
				"$in": stores,
			}
		}
		// if len(channels) == 1 && channels[0] == "store" && len(stores) == 0 {
		// 	queryFilter = bson.M{
		// 		"channel": bson.M{
		// 			"$in": []string{"xxx"},
		// 		},
		// 	}
		// } else {
		// 	queryFilter = bson.M{
		// 		"channel": bson.M{
		// 			"$in": channels,
		// 		},
		// 	}
		// }

	} else if len(channels) == 0 {
		queryFilter = bson.M{
			"channel": bson.M{
				"$in": []string{"xxx"},
			},
		}
	}
	var query = bson.M{
		"$and": []bson.M{
			bson.M{
				"created_at": bson.M{
					"$lte": end,
				},
			},
			bson.M{
				"created_at": bson.M{
					"$gte": start,
				},
			},
			queryFilter,
		},
	}
	ResultTable.FindWhere(query, &surveyResult)
	var timelineReport = initTimelineReport()
	var set, _ = setting.GetSetting()
	glog.Info(set)
	if set != nil {
		if surveyResult != nil {
			for _, item := range surveyResult {
				var hour = time.Unix(int64(item.CreatedAt/1000), 0).Hour()
				if item.AveragePoint >= set.HighRate {
					timelineReport[hour].High++
				} else if item.AveragePoint >= set.CreditRate {
					timelineReport[hour].Credit++
				} else if item.AveragePoint >= set.MediumRate {
					timelineReport[hour].Medium++
				} else {
					timelineReport[hour].Low++
				}
			}
			return timelineReport
		}
	}
	return nil
}

func GetTimelineReport(start, end int, stores, unames, devices, channels []string) []*Timeline {
	var timeLineReport = processTimeline(start, end, stores, unames, devices, channels)
	var timeLines = make([]*Timeline, 0)
	if timeLineReport != nil {
		for k := 7; k <= 22; k++ {
			// if k%2 != 0 {
			// 	continue
			// }
			timeLines = append(timeLines, &Timeline{
				Time:   fmt.Sprintf("%d.00 - %d.00", k, k+1),
				High:   timeLineReport[k].High,
				Credit: timeLineReport[k].Credit,
				Medium: timeLineReport[k].Medium,
				Low:    timeLineReport[k].Low,
			})
		}
	} else {
		for k := 7; k <= 22; k++ {
			timeLines = append(timeLines, &Timeline{
				Time:   fmt.Sprintf("%d.00 - %d.00", k, k+1),
				High:   0,
				Credit: 0,
				Medium: 0,
				Low:    0,
			})
		}
	}
	return timeLines
}

type TimelineReport map[int]*Timeline

func initTimelineReport() TimelineReport {
	var timelineReport = TimelineReport{}
	for i := 0; i < 24; i++ {
		timelineReport[i] = &Timeline{}
	}
	return timelineReport
}
