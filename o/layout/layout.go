package layout

import (
	"feedback/x/db/mongodb"
)

type Layout struct {
	mongodb.BaseModel `bson:",inline"`
	Name              string   `bson:"name" json:"name" validate:"required"`
	Channel           string   `bson:"channel" json:"channel" validate:"required"`
	Store             []string `bson:"store" json:"store"`
	Greeting          float32  `bson:"greeting" json:"greeting"`
	I18nGreating      float32  `bson:"i18n_greating" json:"i18n_greating"`
	Goodbye           string   `bson:"goodbye" json:"goodbye"`
	I18nGoodbye       string   `bson:"i18n_goodbye" json:"i18n_goodbye"`
}

var layoutTable = mongodb.NewTable("campaign", "layout", 5)
