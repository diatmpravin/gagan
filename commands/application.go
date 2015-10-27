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
	request, err := api.NewAuthorizedRequest("POST", path, config.AccessToken, strings.NewReader(data))
	if err != nil {
		return
	}

	resource := new(Resource)
	err = api.PerformRequestAndParseResponse(request, resource)

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
	request, err := api.NewAuthorizedRequest("GET", path, config.AccessToken, nil)
	if err != nil {
		return
	}

	apiResponse := InstancesApiResponse{}

	err = api.PerformRequestAndParseResponse(request, &apiResponse)
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
	request, err := api.NewAuthorizedRequest("PUT", path, config.AccessToken, strings.NewReader(body))

	if err != nil {
		return
	}

	return api.PerformRequest(request)
}

func validateApplication(app models.Application) (err error) {
	reg := regexp.MustCompile("^[0-9a-zA-Z\\-_]*$")
	if !reg.MatchString(app.Name) {
		err = errors.New("Application name is invalid. Name can only contain letters, numbers, underscores and hyphens.")
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

// GetAppSummary GET details of particulat app
func GetAppSummary(w http.ResponseWriter, r *http.Request) {
	render := &api.Render{r, w}

	config := configuration.GetDefaultConfig()
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	appName := r.URL.Query().Get("appname")
	repo := CloudControllerApplicationRepository{}
	app, err := repo.FindByName(config, appName)
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

	appName := r.URL.Query().Get("appname")
	repo := CloudControllerApplicationRepository{}

	app, err := repo.FindByName(config, appName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = repo.Stop(config, app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	app, err = repo.FindByName(config, appName)
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

	appName := r.URL.Query().Get("appname")
	repo := CloudControllerApplicationRepository{}

	app, err := repo.FindByName(config, appName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if app.State == "started" {
		http.Error(w, fmt.Sprintf("Application %s is already started.", appName), http.StatusBadRequest)
	}

	err = repo.Start(config, app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	instances, err := repo.GetInstances(config, app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Printf("App Instances: %+v", instances)

	app, err = repo.FindByName(config, appName)
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

	appName := r.URL.Query().Get("appname")
	repo := CloudControllerApplicationRepository{}

	app, err := repo.FindByName(config, appName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = repo.Delete(config, app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	res := map[string]string{"appname": appName, "message": "App deleted successfully"}

	log.Printf("App status: %+v", res)
	render.JSON(res)
}

// CreatingAnApp will create app
func CreatingAnApp(w http.ResponseWriter, r *http.Request) {
	render := &api.Render{r, w}

	config := configuration.GetDefaultConfig()
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	appName := r.URL.Query().Get("appname")
	repo := CloudControllerApplicationRepository{}

	app, err := repo.FindByName(config, appName)

	if err != nil {
		app, err = createApp(config, appName)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// TODO, need to implement
	// err = repo.Upload(config, app)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	log.Printf("App details: %+v", app)
	render.JSON(app)
}

func createApp(config *configuration.Configuration, appName string) (app models.Application, err error) {
	newApp := models.Application{Name: appName}
	repo := CloudControllerApplicationRepository{}

	log.Printf("Creating %s...", appName)
	app, err = repo.Create(config, newApp)
	if err != nil {
		err = errors.New("Error creating application.")
		return
	}

	domainRepo := CloudControllerDomainRepository{}
	domains, err := domainRepo.FindAll(config)

	if err != nil {
		err = errors.New("Error loading domains")
		return
	}

	domain := domains[0]
	newRoute := models.Route{Host: app.Name}

	routeRepo := CloudControllerRouteRepository{}
	log.Printf("Creating route %s.%s...", app.Name, domain.Name)
	createdRoute, err := routeRepo.Create(config, newRoute, domain)
	if err != nil {
		err = errors.New("Error creating route")
		return
	}

	log.Printf("Binding %s.%s to %s...", app.Name, domain.Name, app.Name)
	err = routeRepo.Bind(config, createdRoute, app)
	if err != nil {
		err = errors.New("Error binding route")
		return
	}

	return
}

// GetTheInstanceInformation GET instance information of particular app
func GetTheInstanceInformation(w http.ResponseWriter, r *http.Request) {
	render := &api.Render{r, w}

	config := configuration.GetDefaultConfig()
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	appName := r.URL.Query().Get("appname")
	repo := CloudControllerApplicationRepository{}
	app, err := repo.FindByName(config, appName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	instances, err := repo.GetInstances(config, app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Printf("Detials of a app instances: %+v", instances)
	render.JSON(instances)
}
