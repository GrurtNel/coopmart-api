package result

import (
	"feedback/x/db/mongodb"
)

type CampaignAggregate struct {
	mongodb.BaseModel `bson:",inline"`
	StoreID           string  `bson:"store_id" json:"store_id"`
	StoreName         string  `bson:"store_name" json:"store_name"`
	DeviceID          string  `bson:"device_id" json:"device_id"`
	DeviceName        string  `bson:"device_name" json:"device_name"`
	UName             string  `bson:"uname" json:"uname"`
	Campaign          string  `bson:"campaign" json:"campaign"`
	Channel           string  `bson:"channel" json:"channel"`
	Location          string  `bson:"location" json:"location"`
	AveragePoint      float32 `bson:"average_point" json:"average_point"`
	Finished          bool    `bson:"finished" json:"finished"`
}

var CampaignAggregateTable = mongodb.NewTable("campaign_aggregate", "CAG", 5)

func ConvertToCampaignAggregate(s *SurveyResult) {
	var a = CampaignAggregate{
		StoreID:    s.Store,
		StoreName:  s.Store,
		DeviceID:   s.Device,
		DeviceName: s.Device,
		UName:      s.Uname,
		Campaign:   s.Campaign,
		Channel:    s.Channel,
		Location:   s.Location,
	}
	var avgPoint float32
	var point int
	var maxPoint int
	for _, survey := range s.SurveyDetails {
		point += survey.Point
		maxPoint += survey.MaxPoint
	}
	if maxPoint != 0 {
		avgPoint += float32(point) / float32(maxPoint) * 10
	} else {
		a.AveragePoint = 0
	}
	a.AveragePoint = avgPoint
	a.Channel = s.Channel
	CampaignAggregateTable.Create(&a)
}
