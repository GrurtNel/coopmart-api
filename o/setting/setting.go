package setting

import (
	"feedback/x/db/mongodb"

	"gopkg.in/mgo.v2/bson"
)

type Setting struct {
	mongodb.BaseModel `bson:",inline"`
	Heading           string  `bson:"heading" json:"heading" validate:"required"`
	Logo              string  `bson:"logo" json:"logo" validate:"required"`
	Background        string  `bson:"background" json:"background" validate:"required"`
	MediumRate        float32 `bson:"medium_rate" json:"medium_rate"`
	CreditRate        float32 `bson:"credit_rate" json:"credit_rate"`
	HighRate          float32 `bson:"high_rate" json:"high_rate"`
}
type LuckySetting struct {
	mongodb.BaseModel `bson:",inline"`
	LuckyNumber       int    `bson:"lucky_number" json:"lucky_number" validate:"required"`
	BonusContent      string `bson:"bonus_content" json:"bonus_content"`
	Activated         bool   `bson:"activated" json:"activated"`
}

var SettingTable = mongodb.NewTable("setting", "SET", 5)

// SYSTEMSETTING system setting
const SYSTEMSETTING = "system_setting"
const LUCKYSETTING = "lucky_setting"

func (s *Setting) Create() error {
	s.ID = SYSTEMSETTING
	_, err := SettingTable.Upsert(bson.M{"_id": SYSTEMSETTING}, s)
	return err
}

func (s *LuckySetting) Create() error {
	s.ID = LUCKYSETTING
	_, err := SettingTable.Upsert(bson.M{"_id": LUCKYSETTING}, s)
	return err
}

func GetSetting() (*Setting, error) {
	var setting *Setting
	err := SettingTable.FindId(SYSTEMSETTING).One(&setting)
	if err != nil {
		if err.Error() == "not found" {
			return nil, nil
		}
	}
	return setting, err
}
func GetLuckySetting() (*LuckySetting, error) {
	var setting *LuckySetting
	err := SettingTable.FindId(LUCKYSETTING).One(&setting)
	if err != nil {
		if err.Error() == "not found" {
			return nil, nil
		}
	}
	return setting, err
}
