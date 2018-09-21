package campaign

import (
	"feedback/o/survey"
	"feedback/x/db/mongodb"
)

type Campaign struct {
	mongodb.BaseModel `bson:",inline"`
	Name              string    `bson:"name" json:"name" validate:"required"`
	Surveys           []string  `bson:"surveys" json:"surveys" validate:"required"`
	Devices           []string  `bson:"devices,omitempty" json:"devices"`
	Channels          []Channel `bson:"channels" json:"channels"`
	Channel           Channel   `bson:"channel,omitempty" json:"channel"`
	Device            string    `bson:"device,omitempty" json:"device"`
	Start             int64     `bson:"start" json:"start" validate:"required"`
	End               int64     `bson:"end" json:"end" validate:"required"`
}
type CampaignView struct {
	Campaign `bson:",inline"`
	Survey   []*survey.Survey `bson:"survey" json:"survey" validate:"required"`
}
