package survey

import (
	"feedback/x/rest"

	"gopkg.in/mgo.v2/bson"
)

func CheckExistSurveys(surveys []string) error {
	var count int
	count, err := SurveyTable.Find(bson.M{"_id": bson.M{"$in": surveys}}).Count()
	if len(surveys) != count {
		return rest.BadRequest("not exists any survey")
	}
	return err
}
