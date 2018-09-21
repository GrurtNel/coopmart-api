package luckyuser

import (
	"feedback/x/db/mongodb"

	"gopkg.in/mgo.v2/bson"
)

type LuckyUser struct {
	mongodb.BaseModel `bson:",inline"`
	BonusContent      string `bson:"bonus_content" json:"bonus_content" validate:"required"`
	Phone             string `bson:"phone" json:"phone" validate:"required"`
	Name              string `bson:"name" json:"name" validate:"required"`
	GiftCode          int    `bson:"gift_code" json:"gift_code"`
}

var luckyUserTable = mongodb.NewTable("lucky_user", "lu", 10)

func (l *LuckyUser) Create() error {
	return luckyUserTable.Create(l)
}

func GetLuckyUsers() []*LuckyUser {
	var luckyUsers []*LuckyUser
	luckyUserTable.Find(bson.M{}).All(&luckyUsers)
	return luckyUsers
}
