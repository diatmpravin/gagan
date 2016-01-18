package commands

import (
	"encoding/json"
	"github.com/diatmpravin/gagan/api"
	"github.com/diatmpravin/gagan/configuration"
	"github.com/diatmpravin/gagan/models"
	"github.com/diatmpravin/gagan/requirements"
	"log"
	"strconv"
	"net/http"
)

type App struct {
	config  *configuration.Configuration
	appRepo api.ApplicationRepository
	appReq  requirements.ApplicationRequirement
}

func NewAppSummary(config *configuration.Configuration, appRepo api.ApplicationRepository) (a *App) {
	a = new(App)
	a.config = config
	a.appRepo = appRepo

	return
}

func (a *App) GetRequirements(reqFactory requirements.Factory, w http.ResponseWriter, r *http.Request) (reqs []Requirement, config *configuration.Configuration, err error) {
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
	a.config = config

	appName := r.URL.Query().Get("appname")
	a.appReq = reqFactory.NewApplicationRequirement(appName)

	reqs = []Requirement{&a.appReq}
	return
}

func (a *App) Run(w http.ResponseWriter, r *http.Request) {
	a.ApplicationSummary(a.appReq.Application, w, r)
}

func (a *App) ApplicationSummary(app models.Application, w http.ResponseWriter, r *http.Request) {
	render := &api.Render{r, w}

	app, err := a.appRepo.FindByName(a.config, app.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Printf("App status: %+v", app)
	render.JSON(app)
}
