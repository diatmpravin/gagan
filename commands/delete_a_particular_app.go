package commands

import (
	"encoding/json"
	"github.com/diatmpravin/gagan/api"
	"github.com/diatmpravin/gagan/configuration"
	"github.com/diatmpravin/gagan/requirements"
	"log"
	"net/http"
)

type Delete struct {
	config  *configuration.Configuration
	appRepo api.ApplicationRepository
	appReq  requirements.ApplicationRequirement
}

func NewDelete(config *configuration.Configuration, appRepo api.ApplicationRepository) (d *Delete) {
	d = new(Delete)
	d.config = config
	d.appRepo = appRepo

	return
}

func (d *Delete) GetRequirements(reqFactory requirements.Factory, w http.ResponseWriter, r *http.Request) (reqs []Requirement, config *configuration.Configuration, err error) {
	session := configuration.Session{}

	if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	config = configuration.GetConfig(session)
	d.config = config

	appName := r.URL.Query().Get("appname")
	d.appReq = reqFactory.NewApplicationRequirement(appName)

	reqs = []Requirement{&d.appReq}
	return
}

func (d *Delete) Run(w http.ResponseWriter, r *http.Request) {
	render := &api.Render{r, w}
	app := d.appReq.Application

	err := d.appRepo.Delete(d.config, app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	res := map[string]string{"appname": app.Name, "message": "App deleted successfully"}
	log.Printf("App status: %+v", res)
	render.JSON(res)
}
