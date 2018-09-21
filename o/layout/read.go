package layout

import (
	"gopkg.in/mgo.v2/bson"
)

func GetLayouts() ([]*Layout, error) {
	var layouts []*Layout
	var err = layoutTable.FindWhere(bson.M{}, &layouts)
	return layouts, err
}
