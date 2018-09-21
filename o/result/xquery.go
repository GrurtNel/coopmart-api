package result

import (
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

const (
	MATCH_STATE   = "$match"
	GROUP_STATE   = "$match"
	PROJECT_STATE = "$project"
)

func getState(state string, query []interface{}) bson.M {
	var pipe = bson.M{}
	for i := 0; i < len(query)/2; i++ {
		pipe[query[i*2].(string)] = query[i*2+1]
	}
	return bson.M{state: pipe}
}

func cond(condition interface{}, trueR, falseR interface{}) bson.M {
	return bson.M{
		"$cond": []interface{}{condition, trueR, falseR},
	}
}
func matchAggregateDuration(start, end int) bson.M {
	return bson.M{
		"$match": bson.M{
			"$and": []bson.M{
				bson.M{
					"created_at": bson.M{
						"$lte": end,
					},
				},
				bson.M{
					"created_at": bson.M{
						"$gte": start,
					},
				},
			},
		},
	}
}

func matchAggregateDurationStore(start, end int) bson.M {
	return bson.M{
		"$match": bson.M{
			"$and": []bson.M{
				bson.M{
					"created_at": bson.M{
						"$lte": end,
					},
				},
				bson.M{
					"created_at": bson.M{
						"$gte": start,
					},
				},
				bson.M{
					"channel": bson.M{
						"$eq": "store",
					},
				},
			},
		},
	}
}

func matchAggregateQuery(q func(string) string) bson.M {
	var start, _ = strconv.Atoi(q("start"))
	var end, _ = strconv.Atoi(q("end"))
	return matchAggregateDuration(start, end)
}

func groupByLevel(field string) bson.M {
	return bson.M{
		"$sum": bson.M{
			"$cond": []interface{}{
				bson.M{"$eq": []interface{}{field, true}}, 1, 0,
			},
		},
	}
}

func unwind(field string) bson.M {
	return bson.M{"$unwind": field}
}
