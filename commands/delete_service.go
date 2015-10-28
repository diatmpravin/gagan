package commands

import (
	"encoding/json"
	"github.com/diatmpravin/gagan/api"
	"github.com/diatmpravin/gagan/configuration"
	"log"
	"net/http"
)

type DeleteService struct {
	config      *configuration.Configuration
	serviceRepo api.CloudControllerServiceRepository
}

// DeleteParticularService delete a particular service
func DeleteParticularService(w http.ResponseWriter, r *http.Request) {
	render := &api.Render{r, w}
	config := configuration.GetDefaultConfig()
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	instanceName := r.URL.Query().Get("service")

	bindRepo := DeleteService{}
	instance, err := bindRepo.serviceRepo.FindInstanceByName(config, instanceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Printf("Deleting service %s...", instance.Name)
	err = bindRepo.serviceRepo.DeleteService(config, instance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Printf("Service details: %+v", config)
	render.JSON(config)
}
