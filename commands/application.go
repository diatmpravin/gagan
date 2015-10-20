package commands

import (
	"encoding/json"
	"fmt"
	"github.com/diatmpravin/gagan/api"
	"github.com/diatmpravin/gagan/configuration"
	"github.com/diatmpravin/gagan/models"
	"log"
	"net/http"
	"strings"
)

type CloudControllerApplicationRepository struct {
}

type ApplicationRepository interface {
	FindApps(config *configuration.Configuration) (apps []models.Application, err error)
}

func (repo CloudControllerApplicationRepository) FindApps(config *configuration.Configuration) (apps []models.Application, err error) {
	path := fmt.Sprintf("%s/v2/spaces/%s/apps?inline-relations-depth=2", config.Target, config.Space.Guid)
	request, err := api.NewAuthorizedRequest("GET", path, config.AccessToken, nil)
	if err != nil {
		return
	}

	response := new(ApplicationsApiResponse)
	err = api.PerformRequestAndParseResponse(request, response)
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

// ListAllApps GET list of all apps
func ListAllApps(w http.ResponseWriter, r *http.Request) {
	render := &api.Render{r, w}

	config := configuration.GetDefaultConfig()
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	repo := CloudControllerApplicationRepository{}
	apps, err := repo.FindApps(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Printf("List of all apps: %+v", apps)
	render.JSON(apps)
}
