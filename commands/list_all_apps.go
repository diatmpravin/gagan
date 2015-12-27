package commands

import (
	"encoding/json"
	"github.com/diatmpravin/gagan/api"
	"github.com/diatmpravin/gagan/configuration"
	"github.com/diatmpravin/gagan/requirements"
	"log"
	"net/http"
	"strconv"
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

	session := configuration.Session{}

	if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	config.Organization.Name = session.Organization.Name
	config.Organization.Guid = session.Organization.Guid
	config.Space.Name = session.Space.Name
	config.Space.Guid = session.Space.Guid

	c := configuration.RedisConnect()
	defer c.Close()

	reply, err := c.Do("GET", "user:"+strconv.Itoa(session.SessionId))

	configuration.HandleError(err)

	if err = json.Unmarshal(reply.([]byte), &session); err != nil {
		panic(err)
	}

	config.AccessToken = session.AccessToken

	space, err := a.spaceRepo.GetSummary(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	apps := space.Applications

	log.Printf("List of all apps: %+v", apps)
	render.JSON(apps)
}
