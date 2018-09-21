package mongodb

import (
	"errors"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	ERR_NOT_EXIST = "not found"
)

func CheckExist(collection *mgo.Collection, object interface{}, filter bson.M) error {
	count, _ := collection.Find(filter).Count()
	if count > 0 {
		return errors.New("exists a unique field")
	}
	return nil
}
func ReadIfExist(collection *mgo.Collection, object interface{}, filter bson.M) {
	err := collection.Find(filter).One(&object)
	if err != nil {
		if err.Error() == ERR_NOT_EXIST {

		}
	}
}
