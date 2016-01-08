package commands

import (
	"encoding/json"
	"fmt"
	"github.com/diatmpravin/gagan/api"
	"github.com/diatmpravin/gagan/configuration"
	"github.com/diatmpravin/gagan/requirements"
	"log"
	"net/http"
	"strconv"
)

type Stop struct {
	config  *configuration.Configuration
	appRepo api.ApplicationRepository
	appReq  requirements.ApplicationRequirement
}

func NewStop(config *configuration.Configuration, appRepo api.ApplicationRepository) (s *Stop) {
	s = new(Stop)
	s.config = config
	s.appRepo = appRepo

	return
}

func (s *Stop) GetRequirements(reqFactory requirements.Factory, w http.ResponseWriter, r *http.Request) (reqs []Requirement, config *configuration.Configuration, err error) {
	config = configuration.GetDefaultConfig()
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
	s.config = config

	appName := r.URL.Query().Get("appname")
	s.appReq = reqFactory.NewApplicationRequirement(appName)

	reqs = []Requirement{&s.appReq}
	return
}

func (s *Stop) Run(w http.ResponseWriter, r *http.Request) {
	render := &api.Render{r, w}

	app := s.appReq.Application

	if app.State == "stopped" {
		http.Error(w, fmt.Sprintf("Application %s is already stopped.", app.Name), http.StatusBadRequest)
		return
	}

	err := s.appRepo.Stop(s.config, app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	app, err = s.appRepo.FindByName(s.config, app.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Printf("App status: %+v", app)
	render.JSON(app)
}
