package commands

import (
	"encoding/json"
	"github.com/diatmpravin/gagan/api"
	"github.com/diatmpravin/gagan/configuration"
	"log"
	"net/http"
)

type BindService struct {
	config      *configuration.Configuration
	serviceRepo api.CloudControllerServiceRepository
	appRepo     api.CloudControllerApplicationRepository
}

// CreatingServiceBinding bind a service with particular app
func CreatingServiceBinding(w http.ResponseWriter, r *http.Request) {
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

	log.Printf("Binding service %s to %s...", instance.Name, app.Name)
	err = bindRepo.serviceRepo.BindService(config, instance, app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Printf("Service details: %+v", config)
	render.JSON(config)
}
