package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/diatmpravin/gagan/api"
	"github.com/diatmpravin/gagan/configuration"
	"github.com/diatmpravin/gagan/models"
	"log"
	"net/http"
)

// CreatingServiceInstance create a service instance
func CreatingServiceInstance(w http.ResponseWriter, r *http.Request) {
	render := &api.Render{r, w}

	config := configuration.GetDefaultConfig()
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	name := r.URL.Query().Get("name")
	offeringName := r.URL.Query().Get("offering")
	planName := r.URL.Query().Get("plan")

	serviceRepo := api.CloudControllerServiceRepository{}
	offerings, err := serviceRepo.GetServiceOfferings(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	offering, err := findOffering(offerings, offeringName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	plan, err := findPlan(offering.Plans, planName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Printf("Creating service %s", name)
	err = serviceRepo.CreateServiceInstance(config, name, plan)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// TODO, find service by name and send json response
	log.Printf("Service details: %+v", config)
	render.JSON(config)
}

func findOffering(offerings []models.ServiceOffering, name string) (offering models.ServiceOffering, err error) {
	for _, offering := range offerings {
		if name == offering.Label {
			return offering, nil
		}
	}

	err = errors.New(fmt.Sprintf("Could not find offering with name %s", name))
	return
}

func findPlan(plans []models.ServicePlan, name string) (plan models.ServicePlan, err error) {
	for _, plan := range plans {
		if name == plan.Name {
			return plan, nil
		}
	}

	err = errors.New(fmt.Sprintf("Could not find plan with name %s", name))
	return
}
