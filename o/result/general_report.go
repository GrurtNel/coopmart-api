package result

import (
	"gopkg.in/mgo.v2/bson"
)

type ActivityFrequency struct {
	Day        string `bson:"day" json:"day"`
	Count      int    `bson:"count" json:"count"`
	Store      int    `bson:"store" json:"store"`
	Sms        int    `bson:"sms" json:"sms"`
	CallCenter int    `bson:"call_center" json:"call_center"`
	Website    int    `bson:"website" json:"website"`
}

// $group:{
// 	_id:{
// 		ctime:"$ctime"
// 	},
// 	store:{$sum:{ $cond: [ {$eq:["$channel","store"]}, 1, 0 ] }},
// 	website:{$sum:{ $cond: [ {$eq:["$channel","website"]}, 1, 0 ] }},
// 	other:{$sum:{ $cond: [ {$eq:["$channel",""]}, 1, 0 ] }},
// 	count:{$sum:1}

// }
func GetActivityFrequencyChart(q func(string) string) ([]*ActivityFrequency, error) {
	var res []*ActivityFrequency
	var match = matchAggregateQuery(q)
	var sum = func(field, value string) bson.M {
		return bson.M{"$sum": cond(bson.M{"$eq": []string{field, value}}, 1, 0)}
	}
	var group = bson.M{
		"$group": bson.M{
			"_id": bson.M{
				"ctime": "$ctime",
			},
			"day":         bson.M{"$first": "$ctime"},
			"created_at":  bson.M{"$first": "$created_at"},
			"store":       sum("$channel", "store"),
			"website":     sum("$channel", "website"),
			"call_center": sum("$channel", "call_center"),
			"sms":         sum("$channel", "sms"),
			"count":       bson.M{"$sum": 1},
		},
	}
	var sort = bson.M{"$sort": bson.M{"created_at": 1}}
	err := ResultTable.Pipe([]bson.M{match, group, sort}).All(&res)
	return res, err
}
