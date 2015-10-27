package commands

import (
	"fmt"
	"github.com/diatmpravin/gagan/api"
	"github.com/diatmpravin/gagan/configuration"
	"github.com/diatmpravin/gagan/models"
	"log"
	"strings"
)

type RouteRepository interface {
	Create(config *configuration.Configuration, newRoute models.Route, domain models.Domain) (createdRoute models.Route, err error)
	Bind(config *configuration.Configuration, route models.Route, app models.Application) (err error)
}

type CloudControllerRouteRepository struct {
}

func (repo CloudControllerRouteRepository) Create(config *configuration.Configuration, newRoute models.Route, domain models.Domain) (createdRoute models.Route, err error) {
	path := fmt.Sprintf("%s/v2/routes", config.Target)
	data := fmt.Sprintf(
		`{"host":"%s","domain_guid":"%s","space_guid":"%s"}`,
		newRoute.Host, domain.Guid, config.Space.Guid,
	)
	request, err := api.NewAuthorizedRequest("POST", path, config.AccessToken, strings.NewReader(data))
	if err != nil {
		return
	}

	resource := new(api.Resource)
	err = api.PerformRequestAndParseResponse(request, resource)
	if err != nil {
		return
	}

	createdRoute.Guid = resource.Metadata.Guid
	createdRoute.Host = resource.Entity.Host

	log.Printf("App Route: %+v", createdRoute)
	return
}

func (repo CloudControllerRouteRepository) Bind(config *configuration.Configuration, route models.Route, app models.Application) (err error) {
	path := fmt.Sprintf("%s/v2/apps/%s/routes/%s", config.Target, app.Guid, route.Guid)
	request, err := api.NewAuthorizedRequest("PUT", path, config.AccessToken, nil)
	if err != nil {
		return
	}

	err = api.PerformRequest(request)

	log.Printf("App Route after Bind: %+v", app)
	return
}
