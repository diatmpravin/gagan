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
	session := configuration.Session{}

	if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	config = configuration.GetConfig(session)
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
