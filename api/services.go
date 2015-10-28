package api

import (
	"fmt"
	"github.com/diatmpravin/gagan/configuration"
	"github.com/diatmpravin/gagan/models"
	"strings"
)

type ServiceRepository interface {
	GetServiceOfferings(config *configuration.Configuration) (offerings []models.ServiceOffering, err error)
	CreateServiceInstance(config *configuration.Configuration, name string, plan models.ServicePlan) (err error)
}

type CloudControllerServiceRepository struct {
}

func (repo CloudControllerServiceRepository) GetServiceOfferings(config *configuration.Configuration) (offerings []models.ServiceOffering, err error) {
	path := fmt.Sprintf("%s/v2/services?inline-relations-depth=1", config.Target)
	request, err := NewAuthorizedRequest("GET", path, config.AccessToken, nil)
	if err != nil {
		return
	}

	response := new(ServiceOfferingsApiResponse)

	_, err = PerformRequestAndParseResponse(request, response)

	if err != nil {
		return
	}

	for _, r := range response.Resources {
		plans := []models.ServicePlan{}
		for _, p := range r.Entity.ServicePlans {
			plans = append(plans, models.ServicePlan{Name: p.Entity.Name, Guid: p.Metadata.Guid})
		}
		offerings = append(offerings, models.ServiceOffering{Label: r.Entity.Label, Guid: r.Metadata.Guid, Plans: plans})
	}

	return
}

func (repo CloudControllerServiceRepository) CreateServiceInstance(config *configuration.Configuration, name string, plan models.ServicePlan) (err error) {
	path := fmt.Sprintf("%s/v2/service_instances", config.Target)

	data := fmt.Sprintf(
		`{"name":"%s","service_plan_guid":"%s","space_guid":"%s"}`,
		name, plan.Guid, config.Space.Guid,
	)
	request, err := NewAuthorizedRequest("POST", path, config.AccessToken, strings.NewReader(data))
	if err != nil {
		return
	}

	_, err = PerformRequest(request)
	return
}
