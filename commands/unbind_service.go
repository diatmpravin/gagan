package commands

import (
	"encoding/json"
	"github.com/diatmpravin/gagan/api"
	"github.com/diatmpravin/gagan/configuration"
	"log"
	"net/http"
)

type UnbindService struct {
	config      *configuration.Configuration
	serviceRepo api.CloudControllerServiceRepository
	appRepo     api.CloudControllerApplicationRepository
}

// DeleteParticularServiceBinding unbind a service with particular app
func DeleteParticularServiceBinding(w http.ResponseWriter, r *http.Request) {
	render := &api.Render{r, w}
	config := configuration.GetDefaultConfig()
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	appName := r.URL.Query().Get("app")
	instanceName := r.URL.Query().Get("service")

	bindRepo := BindService{}
	app, err := bindRepo.appRepo.FindByName(config, appName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	instance, err := bindRepo.serviceRepo.FindInstanceByName(config, instanceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Printf("Unbinding service %s to %s...", instance.Name, app.Name)
	err = bindRepo.serviceRepo.UnbindService(config, instance, app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Printf("Service details: %+v", config)
	render.JSON(config)
}
