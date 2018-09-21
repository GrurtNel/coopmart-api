package campaign

import (
	config "feedback/init"
	"feedback/x/rest"
	"gopkg.in/mgo.v2/bson"
)

const (
	WEBSITE_CHANNEL = "website"
	STORE_CHANNEL   = "store"
	SMS_CHANNEL     = "sms"
	CALL_CENTER     = "call_center"
	BOOKING_ONLINE  = "booking_online"
)

type Channel string
type ChannelLink struct {
	Channel Channel `json:"channel"`
	Link    string  `json:"link"`
}

func GetAllChannelLink() []ChannelLink {
	return []ChannelLink{
		{Channel: WEBSITE_CHANNEL, Link: config.DeviceConfig + WEBSITE_CHANNEL},
		{Channel: STORE_CHANNEL, Link: config.DeviceConfig + STORE_CHANNEL},
		{Channel: SMS_CHANNEL, Link: config.DeviceConfig + SMS_CHANNEL},
		{Channel: BOOKING_ONLINE, Link: config.DeviceConfig + BOOKING_ONLINE},
	}
}
func GetAllChannels() []string {
	return []string{WEBSITE_CHANNEL, STORE_CHANNEL, SMS_CHANNEL}
}
func getChannelsExceptStore() []string {
	return []string{WEBSITE_CHANNEL, SMS_CHANNEL, CALL_CENTER}
}

func GetCampaignByDeviceOrChannel(deviceID string, channel string, at int64) (*CampaignView, error) {
	if deviceID == "" && channel == "" {
		return nil, rest.BadRequest("Missing device or channel param")
	}
	var campaign *CampaignView
	var match = bson.M{}
	var join = bson.M{}
	// var subMatch = bson.M{
	// 	"$and": []bson.M{
	// 		bson.M{
	// 			"start": bson.M{
	// 				"$lte": at,
	// 			},
	// 		},
	// 		bson.M{
	// 			"end": bson.M{
	// 				"$gte": at,
	// 			},
	// 		},
	// 		bson.M{
	// 			"updated_at": bson.M{
	// 				"$ne": 0,
	// 			},
	// 		},
	// 		bson.M{
	// 			"channels": bson.M{
	// 				"$eq": channel,
	// 			},
	// 		},
	// 	},
	// }
	var matchChanel = bson.M{
		"$and": []bson.M{
			bson.M{
				"start": bson.M{
					"$lte": at,
				},
			},
			bson.M{
				"end": bson.M{
					"$gte": at,
				},
			},
			bson.M{
				"updated_at": bson.M{
					"$ne": 0,
				},
			},
			bson.M{
				"channels": bson.M{
					"$eq": channel,
				},
			},
		},
	}
	var matchDevice = bson.M{
		"$and": []bson.M{
			bson.M{
				"start": bson.M{
					"$lte": at,
				},
			},
			bson.M{
				"end": bson.M{
					"$gte": at,
				},
			},
			bson.M{
				"updated_at": bson.M{
					"$ne": 0,
				},
			},
			bson.M{
				"devices": bson.M{
					"$eq": deviceID,
				},
			},
		},
	}
	if deviceID != "" {
		// subMatch["devices"] = deviceID
		match["$match"] = matchDevice
	} else {
		// subMatch["channels"] = channel
		match["$match"] = matchChanel
	}
	// match["$match"] = subMatch

	join["$lookup"] = bson.M{
		"from":         "survey",
		"localField":   "surveys",
		"foreignField": "_id",
		"as":           "survey",
	}
	err := CampaignTable.Pipe([]bson.M{
		match,
		join,
	}).One(&campaign)
	if err != nil {
		if rest.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return campaign, err
}
func GetCampaigns() ([]*Campaign, error) {
	var campaigns []*Campaign
	err := CampaignTable.FindWhere(bson.M{}, &campaigns)
	return campaigns, err
}

func GetCampaignByChannels(channels []string, start int64) ([]*Campaign, error) {
	var campaigns []*Campaign
	err := CampaignTable.Pipe([]bson.M{
		bson.M{
			"$match": bson.M{
				"$and": []bson.M{
					bson.M{"channels": bson.M{"$in": channels}},
					// bson.M{"channels": bson.M{"$in": getChannelsExceptStore()}},
					bson.M{"end": bson.M{"$gte": start}},
					bson.M{"start": bson.M{"$lte": start}},
					bson.M{"updated_at": bson.M{"$ne": 0}},
				},
			},
		},
	}).All(&campaigns)
	return campaigns, err
}

func GetCampaignByDevices(deviceArr []string, start int64) ([]*Campaign, error) {
	var campaigns []*Campaign
	err := CampaignTable.Pipe([]bson.M{
		bson.M{
			"$match": bson.M{
				"$and": []bson.M{
					bson.M{"devices": bson.M{"$in": deviceArr}},
					bson.M{"end": bson.M{"$gte": start}},
					bson.M{"start": bson.M{"$lte": start}},
					bson.M{"updated_at": bson.M{"$ne": 0}},
				},
			},
		},
		bson.M{
			"$unwind": "$devices",
		},
		bson.M{
			"$addFields": bson.M{
				"device": "$devices",
			},
		},
	}).All(&campaigns)
	return campaigns, err
}
