package campaign

import (
	"feedback/o/survey"
	"feedback/x/db/mongodb"
	"feedback/x/rest"
	"fmt"
)

var CampaignTable = mongodb.NewTable("campaign", "CAM", 5)

func (c *Campaign) Create() error {
	rest.AssertNil(rest.Validate(c), survey.CheckExistSurveys(c.Surveys))
	rest.AssertNil(c.CheckExistDeviceChannel(nil, nil))
	return CampaignTable.Create(c)
}

func (c *Campaign) CheckExistDeviceChannel(devices []string, channels []string) error {
	var checkingDevices = c.Devices
	if devices != nil {
		checkingDevices = devices
	}
	var checkingChannels = channelToString(c.Channels)
	if channels != nil {
		checkingChannels = channels
	}
	var campaigns, err = GetCampaignByDevices(checkingDevices, c.Start)
	rest.AssertNil(err)
	for _, campaign := range campaigns {
		for _, device := range devices {
			if device == campaign.Device {
				errStr := fmt.Sprintf("Thiết bị %s không thể áp dụng vì đang được áp dụng trong chiến dịch %s", device, campaign.Name)
				return rest.BadRequest(errStr)
			}
		}
	}
	fmt.Println(checkingChannels)
	campaignsByChannel, err := GetCampaignByChannels(checkingChannels, c.Start)
	rest.AssertNil(err)
	for _, campaign := range campaignsByChannel {
		for _, channel := range campaign.Channels {
			for _, chn := range c.Channels {
				if channel == chn {
					errStr := fmt.Sprintf("Kênh %s không thể áp dụng vì đang được áp dụng trong chiến dịch %s", chn, campaign.Name)
					return rest.BadRequest(errStr)
				}
			}
		}
	}
	return nil
}

func channelToString(channels []Channel) []string {
	var result []string
	for _, item := range channels {
		if item != STORE_CHANNEL {
			result = append(result, string(item))
		}
	}
	return result
}

func (c *Campaign) ValidateUpdate(exceptDevice []string, exceptChannel []string) error {
	return c.CheckExistDeviceChannel(getExceptSlice(c.Devices, exceptDevice), getExceptSlice(channelToString(c.Channels), exceptChannel))
}

// getExceptSlice find element slice 1 not in slice 2
func getExceptSlice(slice1 []string, slice2 []string) []string {
	var devicesCond = []string{}
	var found bool
	for _, s1 := range slice1 {
		for _, s2 := range slice2 {
			found = false
			if s1 == s2 {
				found = true
				break
			}
		}
		if !found {
			devicesCond = append(devicesCond, s1)
		}
	}
	return devicesCond
}

func UpdateByID(newCampaign *Campaign) error {
	var campaign *Campaign
	err := CampaignTable.FindByID(newCampaign.ID, &campaign)
	rest.AssertNil(err)
	rest.AssertNil(rest.Validate(newCampaign))
	rest.AssertNil(newCampaign.ValidateUpdate(campaign.Devices, channelToString(campaign.Channels)))
	return CampaignTable.UpdateID(newCampaign.ID, newCampaign)
}
