package api

import (
	"fmt"
	"github.com/diatmpravin/gagan/configuration"
	"github.com/diatmpravin/gagan/models"
)

type CloudControllerAppEventsRepository struct {
}

type AppEventsRepository interface {
	RecentEvents(config *configuration.Configuration, app models.Application) ([]models.EventFields, error)
}

func (repo CloudControllerAppEventsRepository) RecentEvents(config *configuration.Configuration, app models.Application) (events []models.EventFields, err error) {
	path := fmt.Sprintf("%s/v2/app_usage_events/%s", config.Target, app.Guid)
	request, err := NewAuthorizedRequest("GET", path, config.AccessToken, nil)
	if err != nil {
		return
	}

	apiResponse := make([]interface{}, 100)

	_, err = PerformRequestAndParseResponse(request, &apiResponse)
	if err != nil {
		return
	}

	// TODO, need to work on response hander

	return
}
