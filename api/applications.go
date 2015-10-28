package api

import (
	"errors"
	"fmt"
	"github.com/diatmpravin/gagan/configuration"
	"github.com/diatmpravin/gagan/models"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type CloudControllerApplicationRepository struct {
}

type InstancesApiResponse map[string]InstanceApiResponse

type InstanceApiResponse struct {
	State string
}

type ApplicationRepository interface {
	FindApps(config *configuration.Configuration) (apps []models.Application, err error)
	FindByName(config *configuration.Configuration, name string) (app models.Application, err error)
	Stop(config *configuration.Configuration, app models.Application) (err error)
	Start(config *configuration.Configuration, app models.Application) (err error)
	Delete(config *configuration.Configuration, app models.Application) (err error)
	Create(config *configuration.Configuration, newApp models.Application) (createdApp models.Application, err error)
	Upload(config *configuration.Configuration, app models.Application) (err error)
	GetInstances(config *configuration.Configuration, app models.Application) (instance models.ApplicationInstance, err error)
}

func (repo CloudControllerApplicationRepository) FindApps(config *configuration.Configuration) (apps []models.Application, err error) {
	path := fmt.Sprintf("%s/v2/spaces/%s/apps?inline-relations-depth=2", config.Target, config.Space.Guid)
	request, err := NewAuthorizedRequest("GET", path, config.AccessToken, nil)
	if err != nil {
		return
	}

	response := new(ApplicationsApiResponse)
	_, err = PerformRequestAndParseResponse(request, response)
	if err != nil {
		return
	}

	for _, res := range response.Resources {
		urls := []string{}
		for _, routeRes := range res.Entity.Routes {
			routeEntity := routeRes.Entity
			domainEntity := routeEntity.Domain.Entity
			urls = append(urls, fmt.Sprintf("%s.%s", routeEntity.Host, domainEntity.Name))
		}

		apps = append(apps, models.Application{
			Name:      res.Entity.Name,
			Guid:      res.Metadata.Guid,
			State:     strings.ToLower(res.Entity.State),
			Instances: res.Entity.Instances,
			Memory:    res.Entity.Memory,
			Urls:      urls,
		})
	}

	return
}

func (repo CloudControllerApplicationRepository) FindByName(config *configuration.Configuration, name string) (app models.Application, err error) {
	apps, err := repo.FindApps(config)
	if err != nil {
		return
	}

	lowerName := strings.ToLower(name)
	for _, a := range apps {
		if strings.ToLower(a.Name) == lowerName {
			return a, nil
		}
	}

	err = errors.New("Application not found")
	return
}

func (repo CloudControllerApplicationRepository) Delete(config *configuration.Configuration, app models.Application) (err error) {
	path := fmt.Sprintf("%s/v2/apps/%s?recursive=true", config.Target, app.Guid)
	request, err := NewAuthorizedRequest("DELETE", path, config.AccessToken, nil)
	if err != nil {
		return
	}

	_, err = PerformRequest(request)
	return
}

func (repo CloudControllerApplicationRepository) Stop(config *configuration.Configuration, app models.Application) (err error) {
	return changeApplicationState(config, app, "STOPPED")
}

func (repo CloudControllerApplicationRepository) Start(config *configuration.Configuration, app models.Application) (err error) {
	return changeApplicationState(config, app, "STARTED")
}

func (repo CloudControllerApplicationRepository) Create(config *configuration.Configuration, newApp models.Application) (createdApp models.Application, err error) {
	err = validateApplication(newApp)
	if err != nil {
		return
	}

	path := fmt.Sprintf("%s/v2/apps", config.Target)
	data := fmt.Sprintf(
		`{"space_guid":"%s","name":"%s","instances":1,"buildpack":null,"command":null,"memory":256,"stack_guid":null}`,
		config.Space.Guid, newApp.Name,
	)
	request, err := NewAuthorizedRequest("POST", path, config.AccessToken, strings.NewReader(data))
	if err != nil {
		return
	}

	resource := new(Resource)
	_, err = PerformRequestAndParseResponse(request, resource)

	if err != nil {
		return
	}

	createdApp.Guid = resource.Metadata.Guid
	createdApp.Name = resource.Entity.Name

	log.Printf("Created App Details: %+v", createdApp)
	return
}

func (repo CloudControllerApplicationRepository) GetInstances(config *configuration.Configuration, app models.Application) (instances []models.ApplicationInstance, err error) {
	path := fmt.Sprintf("%s/v2/apps/%s/instances", config.Target, app.Guid)
	request, err := NewAuthorizedRequest("GET", path, config.AccessToken, nil)
	if err != nil {
		return
	}

	apiResponse := InstancesApiResponse{}

	_, err = PerformRequestAndParseResponse(request, &apiResponse)
	if err != nil {
		return
	}

	instances = make([]models.ApplicationInstance, len(apiResponse), len(apiResponse))
	for k, v := range apiResponse {
		index, err := strconv.Atoi(k)
		if err != nil {
			continue
		}
		instances[index] = models.ApplicationInstance{State: models.InstanceState(strings.ToLower(v.State))}
	}
	return
}

func changeApplicationState(config *configuration.Configuration, app models.Application, state string) (err error) {
	path := fmt.Sprintf("%s/v2/apps/%s", config.Target, app.Guid)
	body := fmt.Sprintf(`{"console":true,"state":"%s"}`, state)
	request, err := NewAuthorizedRequest("PUT", path, config.AccessToken, strings.NewReader(body))

	if err != nil {
		return
	}

	_, err = PerformRequest(request)
	return
}

func validateApplication(app models.Application) (err error) {
	reg := regexp.MustCompile("^[0-9a-zA-Z\\-_]*$")
	if !reg.MatchString(app.Name) {
		err = errors.New("Application name is invalid. Name can only contain letters, numbers, underscores and hyphens.")
	}

	return
}
