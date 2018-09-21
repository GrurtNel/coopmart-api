package setting

import (
	"feedback/x/db/mongodb"

	"gopkg.in/mgo.v2/bson"
)

type StoreEmail struct {
	mongodb.BaseModel `bson:",inline"`
	Store             string `json:"store" bson:"store"`
	Email             string `json:"email" bson:"email"`
	Channel           string `json:"channel" bson:"channel"`
}

var storeEmailTable = mongodb.NewTable("store_email", "sm", 6)

func (s *StoreEmail) Create() error {
	if GetByStore(s.Store) == nil {
		return storeEmailTable.Create(s)
	}
	if s.Channel != "" {
		return storeEmailTable.Update(bson.M{"channel": s.Channel}, bson.M{"$set": bson.M{
			"email": s.Email,
		}})
	}
	return storeEmailTable.Update(bson.M{"store": s.Store}, bson.M{"$set": bson.M{
		"email": s.Email,
	}})
}

func GetByStore(store string) *StoreEmail {
	var res *StoreEmail
	storeEmailTable.Find(bson.M{"store": store}).One(&res)
	return res
}

func DeleteByID(id string) error {
	return storeEmailTable.DeleteID(id)
}

func GetEmailByStore(store, channel string) string {
	var res *StoreEmail
	if channel == "store" {
		storeEmailTable.Find(bson.M{"store": store}).One(&res)
	} else {
		storeEmailTable.Find(bson.M{"channel": channel}).One(&res)
	}
	if res != nil {
		return res.Email
	}
	return "trunglv1993@gmail.com"
}

func GetEmails() []*StoreEmail {
	var res = []*StoreEmail{}
	storeEmailTable.FindWhere(bson.M{}, &res)
	return res
}
