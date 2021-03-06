package api

import (
	"errors"
	"fmt"
	"github.com/diatmpravin/gagan/configuration"
	"github.com/diatmpravin/gagan/models"
	"strings"
)

type ServiceRepository interface {
	GetServiceOfferings(config *configuration.Configuration) (offerings []models.ServiceOffering, err error)
	CreateServiceInstance(config *configuration.Configuration, name string, plan models.ServicePlan) (err error)
	FindInstanceByName(config *configuration.Configuration, name string) (instance models.ServiceInstance, err error)
	BindService(config *configuration.Configuration, instance models.ServiceInstance, app models.Application) (err error)
	UnbindService(config *configuration.Configuration, instance models.ServiceInstance, app models.Application) (err error)
	DeleteService(config *configuration.Configuration, instance models.ServiceInstance) (err error)
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

func (repo CloudControllerServiceRepository) FindInstanceByName(config *configuration.Configuration, name string) (instance models.ServiceInstance, err error) {
	path := fmt.Sprintf("%s/v2/spaces/%s/service_instances?q=name%s&inline-relations-depth=1", config.Target, config.Space.Guid, "%3A"+name)
	request, err := NewAuthorizedRequest("GET", path, config.AccessToken, nil)
	if err != nil {
		return
	}

	response := new(ApiResponse)
	_, err = PerformRequestAndParseResponse(request, response)
	if err != nil {
		return
	}

	if len(response.Resources) == 0 {
		err = errors.New(fmt.Sprintf("Service %s not found", name))
		return
	}

	resource := response.Resources[0]
	instance.Guid = resource.Metadata.Guid
	instance.Name = resource.Entity.Name
	return
}

func (repo CloudControllerServiceRepository) BindService(config *configuration.Configuration, instance models.ServiceInstance, app models.Application) (err error) {
	path := fmt.Sprintf("%s/v2/service_bindings", config.Target)
	body := fmt.Sprintf(
		`{"app_guid":"%s","service_instance_guid":"%s"}`,
		app.Guid, instance.Guid,
	)
	request, err := NewAuthorizedRequest("POST", path, config.AccessToken, strings.NewReader(body))
	if err != nil {
		return
	}

	_, err = PerformRequest(request)
	return
}

func (repo CloudControllerServiceRepository) UnbindService(config *configuration.Configuration, instance models.ServiceInstance, app models.Application) (err error) {
	path := fmt.Sprintf("%s/v2/service_instances/%s?inline-relations-depth=1", config.Target, instance.Guid)
	request, err := NewAuthorizedRequest("GET", path, config.AccessToken, nil)
	if err != nil {
		return
	}

	response := new(ServiceInstanceApiResponse)

	_, err = PerformRequestAndParseResponse(request, response)
	if err != nil {
		return
	}

	path = ""
	for _, bindingResource := range response.Entity.ServiceBindings {
		if bindingResource.Entity.AppGuid == app.Guid {
			path = config.Target + bindingResource.Metadata.Url
			break
		}
	}

	if path == "" {
		err = errors.New("Error finding service binding")
		return
	}

	request, err = NewAuthorizedRequest("DELETE", path, config.AccessToken, nil)
	if err != nil {
		return
	}

	_, err = PerformRequest(request)
	return
}

func (repo CloudControllerServiceRepository) DeleteService(config *configuration.Configuration, instance models.ServiceInstance) (err error) {
	if len(instance.ServiceBindings) > 0 {
		return errors.New("Cannot delete service instance, apps are still bound to it")
	}

	path := fmt.Sprintf("%s/v2/service_instances/%s", config.Target, instance.Guid)
	request, err := NewAuthorizedRequest("DELETE", path, config.AccessToken, nil)
	if err != nil {
		return
	}

	_, err = PerformRequest(request)
	return
}
