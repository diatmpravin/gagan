package commands

import (
	"encoding/json"
	"errors"
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
	FindByName(config *configuration.Configuration, name string) (app models.Application, err error)
	Stop(config *configuration.Configuration, app models.Application) (err error)
	Start(config *configuration.Configuration, app models.Application) (err error)
	Delete(config *configuration.Configuration, app models.Application) (err error)
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
	request, err := api.NewAuthorizedRequest("DELETE", path, config.AccessToken, nil)
	if err != nil {
		return
	}

	err = api.PerformRequest(request)
	return
}

func (repo CloudControllerApplicationRepository) Stop(config *configuration.Configuration, app models.Application) (err error) {
	return changeApplicationState(config, app, "STOPPED")
}

func (repo CloudControllerApplicationRepository) Start(config *configuration.Configuration, app models.Application) (err error) {
	return changeApplicationState(config, app, "STARTED")
}

func changeApplicationState(config *configuration.Configuration, app models.Application, state string) (err error) {
	path := fmt.Sprintf("%s/v2/apps/%s", config.Target, app.Guid)
	body := fmt.Sprintf(`{"console":true,"state":"%s"}`, state)
	request, err := api.NewAuthorizedRequest("PUT", path, config.AccessToken, strings.NewReader(body))

	if err != nil {
		return
	}

	return api.PerformRequest(request)
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

// GetAppSummary GET details of particulat app
func GetAppSummary(w http.ResponseWriter, r *http.Request) {
	render := &api.Render{r, w}

	config := configuration.GetDefaultConfig()
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	appname := r.URL.Query().Get("appname")
	repo := CloudControllerApplicationRepository{}
	app, err := repo.FindByName(config, appname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Printf("Detials of a app: %+v", app)
	render.JSON(app)
}

// StopAnApp will stop an app
func StopAnApp(w http.ResponseWriter, r *http.Request) {
	render := &api.Render{r, w}

	config := configuration.GetDefaultConfig()
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	appname := r.URL.Query().Get("appname")
	repo := CloudControllerApplicationRepository{}

	app, err := repo.FindByName(config, appname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = repo.Stop(config, app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	app, err = repo.FindByName(config, appname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Printf("App status: %+v", app)
	render.JSON(app)
}

// StartingAnApp will stop an app
func StartingAnApp(w http.ResponseWriter, r *http.Request) {
	render := &api.Render{r, w}

	config := configuration.GetDefaultConfig()
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	appname := r.URL.Query().Get("appname")
	repo := CloudControllerApplicationRepository{}

	app, err := repo.FindByName(config, appname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = repo.Start(config, app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	app, err = repo.FindByName(config, appname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Printf("App status: %+v", app)
	render.JSON(app)
}

// DeleteAPraticularApp will delete a partucular app
func DeleteAPraticularApp(w http.ResponseWriter, r *http.Request) {
	render := &api.Render{r, w}

	config := configuration.GetDefaultConfig()
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	appname := r.URL.Query().Get("appname")
	repo := CloudControllerApplicationRepository{}

	app, err := repo.FindByName(config, appname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = repo.Delete(config, app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	res := map[string]string{"appname": appname, "message": "App deleted successfully"}

	log.Printf("App status: %+v", res)
	render.JSON(res)
}
