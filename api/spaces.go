package api

import (
	"encoding/json"
	"fmt"
	"github.com/diatmpravin/gagan/configuration"
	"github.com/diatmpravin/gagan/models"
	"log"
	"net/http"
	"strings"
)

type SpaceRepository interface {
	FindAllSpaces(config *configuration.Configuration) (spaces []models.Space, err error)
	GetSummary(config *configuration.Configuration) (space models.Space, err error)
}

type CloudControllerSpaceRepository struct {
}

func (repo CloudControllerSpaceRepository) FindAllSpaces(config *configuration.Configuration) (spaces []models.Space, err error) {
	path := fmt.Sprintf("%s/v2/organizations/%s/spaces", config.Target, config.Organization.Guid)
	request, err := NewAuthorizedRequest("GET", path, config.AccessToken, nil)
	if err != nil {
		return
	}

	response := new(ApiResponse)

	_, err = PerformRequestAndParseResponse(request, response)

	if err != nil {
		return
	}

	for _, r := range response.Resources {
		spaces = append(spaces, models.Space{Name: r.Entity.Name, Guid: r.Metadata.Guid})
	}

	return
}

func (repo CloudControllerSpaceRepository) GetSummary(config *configuration.Configuration) (space models.Space, err error) {
	path := fmt.Sprintf("%s/v2/spaces/%s/summary", config.Target, config.Space.Guid)
	request, err := NewAuthorizedRequest("GET", path, config.AccessToken, nil)
	if err != nil {
		return
	}

	response := new(SpaceSummary) // but not an ApiResponse
	_, err = PerformRequestAndParseResponse(request, response)

	if err != nil {
		return
	}

	applications := extractApplicationsFromSummary(response.Apps)
	serviceInstances := extractServiceInstancesFromSummary(response.ServiceInstances, response.Apps)

	space = models.Space{Name: response.Name, Guid: response.Guid, Applications: applications, ServiceInstances: serviceInstances}

	return
}

func extractApplicationsFromSummary(appSummaries []ApplicationSummary) (applications []models.Application) {
	for _, appSummary := range appSummaries {
		app := models.Application{
			Name:             appSummary.Name,
			Guid:             appSummary.Guid,
			Urls:             appSummary.Urls,
			State:            strings.ToLower(appSummary.State),
			Instances:        appSummary.Instances,
			RunningInstances: appSummary.RunningInstances,
			Memory:           appSummary.Memory,
		}
		applications = append(applications, app)
	}

	return
}

func extractServiceInstancesFromSummary(instanceSummaries []ServiceInstanceSummary, appSummaries []ApplicationSummary) (instances []models.ServiceInstance) {
	for _, instanceSummary := range instanceSummaries {
		applicationNames := findApplicationNamesForInstance(instanceSummary.Name, appSummaries)

		planSummary := instanceSummary.ServicePlan
		offeringSummary := planSummary.ServiceOffering

		serviceOffering := models.ServiceOffering{
			Label:    offeringSummary.Label,
			Provider: offeringSummary.Provider,
			Version:  offeringSummary.Version,
		}

		servicePlan := models.ServicePlan{
			Name:            planSummary.Name,
			ServiceOffering: serviceOffering,
		}

		instance := models.ServiceInstance{
			Name:             instanceSummary.Name,
			ServicePlan:      servicePlan,
			ApplicationNames: applicationNames,
		}

		instances = append(instances, instance)
	}

	return
}

// ListAllSpaces Get list of all spaces
func ListAllSpaces(w http.ResponseWriter, r *http.Request) {
	render := &Render{r, w}

	config := configuration.GetDefaultConfig()
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	repo := CloudControllerSpaceRepository{}
	spaces, err := repo.FindAllSpaces(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("List of spaces: %+v", spaces)
	render.JSON(spaces)
}

func findApplicationNamesForInstance(instanceName string, appSummaries []ApplicationSummary) (applicationNames []string) {
	for _, appSummary := range appSummaries {
		for _, name := range appSummary.ServiceNames {
			if name == instanceName {
				applicationNames = append(applicationNames, appSummary.Name)
			}
		}
	}

	return
}
