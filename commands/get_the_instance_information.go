package commands

import (
	"encoding/json"
	"github.com/diatmpravin/gagan/api"
	"github.com/diatmpravin/gagan/configuration"
	"github.com/diatmpravin/gagan/requirements"
	"log"
	"net/http"
)

type InstanceInformation struct {
	config  *configuration.Configuration
	appRepo api.ApplicationRepository
	appReq  requirements.ApplicationRequirement
}

func NewInstanceInformation(config *configuration.Configuration, appRepo api.ApplicationRepository) (i *InstanceInformation) {
	i = new(InstanceInformation)
	i.config = config
	i.appRepo = appRepo

	return
}

func (i *InstanceInformation) GetRequirements(reqFactory requirements.Factory, w http.ResponseWriter, r *http.Request) (reqs []Requirement, config *configuration.Configuration, err error) {
	config = configuration.GetDefaultConfig()
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	i.config = config

	appName := r.URL.Query().Get("appname")
	i.appReq = reqFactory.NewApplicationRequirement(appName)

	reqs = []Requirement{&i.appReq}
	return
}

func (i *InstanceInformation) Run(w http.ResponseWriter, r *http.Request) {
	render := &api.Render{r, w}
	app := i.appReq.Application

	app, err := i.appRepo.FindByName(i.config, app.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	instances, err := i.appRepo.GetInstances(i.config, app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Printf("Detials of a app instances: %+v", instances)
	render.JSON(instances)
}

// GetTheInstanceInformation GET instance information of particular app
// func GetTheInstanceInformation(w http.ResponseWriter, r *http.Request) {
// 	render := &api.Render{r, w}

// 	app, err := repo.FindByName(config, appName)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}

// 	instances, err := repo.GetInstances(config, app)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}

// 	log.Printf("Detials of a app instances: %+v", instances)
// 	render.JSON(instances)
// }
