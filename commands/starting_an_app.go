package commands

import (
	"encoding/json"
	"fmt"
	"github.com/diatmpravin/gagan/api"
	"github.com/diatmpravin/gagan/configuration"
	"github.com/diatmpravin/gagan/models"
	"github.com/diatmpravin/gagan/requirements"
	"log"
	"net/http"
)

type Start struct {
	config  *configuration.Configuration
	appRepo api.ApplicationRepository
	appReq  requirements.ApplicationRequirement
}

func NewStart(config *configuration.Configuration, appRepo api.ApplicationRepository) (s *Start) {
	s = new(Start)
	s.config = config
	s.appRepo = appRepo

	return
}

func (s *Start) GetRequirements(reqFactory requirements.Factory, w http.ResponseWriter, r *http.Request) (reqs []Requirement, config *configuration.Configuration, err error) {
	config = configuration.GetDefaultConfig()
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	s.config = config

	appName := r.URL.Query().Get("appname")
	s.appReq = reqFactory.NewApplicationRequirement(appName)

	reqs = []Requirement{&s.appReq}
	return
}

func (s *Start) Run(w http.ResponseWriter, r *http.Request) {
	s.ApplicationStart(s.appReq.Application, w, r)
}

func (s *Start) ApplicationStart(app models.Application, w http.ResponseWriter, r *http.Request) {
	render := &api.Render{r, w}
	if app.State == "started" {
		http.Error(w, fmt.Sprintf("Application %s is already started.", app.Name), http.StatusBadRequest)
		return
	}

	err := s.appRepo.Start(s.config, app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	instances, err := s.appRepo.GetInstances(s.config, app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Printf("App Instances: %+v", instances)

	app, err = s.appRepo.FindByName(s.config, app.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Printf("App status: %+v", app)
	render.JSON(app)
}
