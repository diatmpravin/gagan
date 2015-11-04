package commands

import (
	"encoding/json"
	"github.com/diatmpravin/gagan/api"
	"github.com/diatmpravin/gagan/configuration"
	"github.com/diatmpravin/gagan/requirements"
	"log"
	"net/http"
)

type Apps struct {
	config    *configuration.Configuration
	spaceRepo api.SpaceRepository
}

func NewAllApps(config *configuration.Configuration, spaceRepo api.SpaceRepository) (a Apps) {
	a.config = config
	a.spaceRepo = spaceRepo

	return
}

func (a Apps) GetRequirements(reqFactory requirements.Factory, w http.ResponseWriter, r *http.Request) (reqs []Requirement, config *configuration.Configuration, err error) {
	return
}

func (a Apps) Run(w http.ResponseWriter, r *http.Request) {
	render := &api.Render{r, w}

	config := configuration.GetDefaultConfig()
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	space, err := a.spaceRepo.GetSummary(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	apps := space.Applications

	log.Printf("List of all apps: %+v", apps)
	render.JSON(apps)
}
