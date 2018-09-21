package campaign

func DeleteByID(id string) error {
	return CampaignTable.DeleteID(id)
}
