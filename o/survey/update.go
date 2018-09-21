package survey

import "feedback/x/rest"

func Update(newSurvey *Survey) error {
	rest.AssertNil(rest.Validate(newSurvey))
	return SurveyTable.UpdateID(newSurvey.ID, newSurvey)
}
