package commands

import (
	"encoding/json"
	"github.com/diatmpravin/gagan/api"
	"github.com/diatmpravin/gagan/configuration"
	"log"
	"net/http"
)

// StopAnApp will stop an app
func StopAnApp(w http.ResponseWriter, r *http.Request) {
	render := &api.Render{r, w}

	config := configuration.GetDefaultConfig()
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	appName := r.URL.Query().Get("appname")
	repo := api.CloudControllerApplicationRepository{}

	app, err := repo.FindByName(config, appName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = repo.Stop(config, app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	app, err = repo.FindByName(config, appName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Printf("App status: %+v", app)
	render.JSON(app)
}
