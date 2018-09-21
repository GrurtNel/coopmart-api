package survey

func DeleteByID(id string) error {
	return SurveyTable.DeleteID(id)
}
