package result

import (
	"feedback/x/db/mongodb"

	"gopkg.in/mgo.v2/bson"
)

type SurveyAggregate struct {
	mongodb.BaseModel `bson:",inline"`
	StoreID           string  `bson:"store_id" json:"store_id"`
	StoreName         string  `bson:"store_name" json:"store_name"`
	DeviceID          string  `bson:"device_id" json:"device_id"`
	DeviceName        string  `bson:"device_name" json:"device_name"`
	UName             string  `bson:"uname" json:"uname"`
	SurveyID          string  `bson:"survey_id" json:"survey_id"`
	SurveyName        string  `bson:"survey_name" json:"survey_name"`
	Campaign          string  `bson:"campaign" json:"campaign"`
	Location          string  `bson:"location" json:"location"`
	High              bool    `bson:"high" json:"high"`
	Credit            bool    `bson:"credit" json:"credit"`
	Medium            bool    `bson:"medium" json:"medium"`
	Low               bool    `bson:"low" json:"low"`
	Channel           string  `bson:"channel" json:"channel"`
	AveragePoint      float32 `bson:"average_point" json:"average_point"`
	Finished          bool    `bson:"finished" json:"finished"`
}

var SurveyAggregateTable = mongodb.NewTable("survey_aggregate", "SAG", 5)

func GetSurveyResults() ([]*SurveyAggregate, error) {
	var a []*SurveyAggregate
	err := ResultTable.Pipe([]bson.M{
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"day": "$created_at",
				},
			},
		},
	}).All(&a)
	return a, err
}

func ConvertToSurveyAggregate(s *SurveyResult) {
	var a = SurveyAggregate{
		StoreID:    s.Store,
		StoreName:  s.Store,
		DeviceID:   s.Device,
		DeviceName: s.Device,
		UName:      s.Uname,
		Campaign:   s.Campaign,
		Location:   s.Location,
	}
	for _, survey := range s.SurveyDetails {
		a.SurveyID = survey.SurveyID
		a.SurveyName = survey.SurveyName
		SurveyAggregateTable.Create(&a)
	}
}
