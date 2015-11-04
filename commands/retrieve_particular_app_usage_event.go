package commands

import (
	"encoding/json"
	"github.com/diatmpravin/gagan/api"
	"github.com/diatmpravin/gagan/configuration"
	"github.com/diatmpravin/gagan/models"
	"github.com/diatmpravin/gagan/requirements"
	"log"
	"net/http"
)

type AppEvents struct {
	config       *configuration.Configuration
	appEventRepo api.AppEventsRepository
	appReq       requirements.ApplicationRequirement
}

func NewAppEvents(config *configuration.Configuration, appEventRepo api.AppEventsRepository) (a *AppEvents) {
	a = new(AppEvents)
	a.config = config
	a.appEventRepo = appEventRepo

	return
}

func (a *AppEvents) GetRequirements(reqFactory requirements.Factory, w http.ResponseWriter, r *http.Request) (reqs []Requirement, config *configuration.Configuration, err error) {
	config = configuration.GetDefaultConfig()
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	a.config = config

	appName := r.URL.Query().Get("appname")
	a.appReq = reqFactory.NewApplicationRequirement(appName)

	reqs = []Requirement{&a.appReq}
	return
}

func (a *AppEvents) Run(w http.ResponseWriter, r *http.Request) {
	a.ListAppEvents(a.appReq.Application, w, r)
}

func (a *AppEvents) ListAppEvents(app models.Application, w http.ResponseWriter, r *http.Request) {
	render := &api.Render{r, w}

	events, err := a.appEventRepo.RecentEvents(a.config, app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Printf("App status: %+v", events)
	render.JSON(events)
}
