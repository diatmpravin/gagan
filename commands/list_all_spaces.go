package commands

import (
	"encoding/json"
	"github.com/diatmpravin/gagan/api"
	"github.com/diatmpravin/gagan/configuration"
	"github.com/diatmpravin/gagan/requirements"
	"log"
	"net/http"
)

type SpaceList struct {
	config    *configuration.Configuration
	spaceRepo api.SpaceRepository
}

func NewSpaceList(config *configuration.Configuration, spaceRepo api.SpaceRepository) (s SpaceList) {
	s.config = config
	s.spaceRepo = spaceRepo

	return
}

func (s SpaceList) GetRequirements(reqFactory requirements.Factory, w http.ResponseWriter, r *http.Request) (reqs []Requirement, config *configuration.Configuration, err error) {
	return
}

func (s SpaceList) Run(w http.ResponseWriter, r *http.Request) {
	render := &api.Render{r, w}

	config := configuration.GetDefaultConfig()
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	spaces, err := s.spaceRepo.FindAllSpaces(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("List of spaces: %+v", spaces)
	render.JSON(spaces)
}
