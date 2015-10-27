package commands

import (
	"encoding/json"
	"errors"
	"github.com/diatmpravin/gagan/api"
	"github.com/diatmpravin/gagan/configuration"
	"github.com/diatmpravin/gagan/models"
	"log"
	"net/http"
)

// CreatingAnApp will create app
func CreatingAnApp(w http.ResponseWriter, r *http.Request) {
	render := &api.Render{r, w}

	config := configuration.GetDefaultConfig()
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	appName := r.URL.Query().Get("appname")
	repo := api.CloudControllerApplicationRepository{}

	app, err := repo.FindByName(config, appName)

	if err != nil {
		app, err = createApp(config, appName)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// TODO, need to implement
	// err = repo.Upload(config, app)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	log.Printf("App details: %+v", app)
	render.JSON(app)
}

func createApp(config *configuration.Configuration, appName string) (app models.Application, err error) {
	newApp := models.Application{Name: appName}
	repo := api.CloudControllerApplicationRepository{}

	log.Printf("Creating %s...", appName)
	app, err = repo.Create(config, newApp)
	if err != nil {
		err = errors.New("Error creating application.")
		return
	}

	domainRepo := api.CloudControllerDomainRepository{}
	domains, err := domainRepo.FindAll(config)

	if err != nil {
		err = errors.New("Error loading domains")
		return
	}

	domain := domains[0]
	newRoute := models.Route{Host: app.Name}

	routeRepo := api.CloudControllerRouteRepository{}
	log.Printf("Creating route %s.%s...", app.Name, domain.Name)
	createdRoute, err := routeRepo.Create(config, newRoute, domain)
	if err != nil {
		err = errors.New("Error creating route")
		return
	}

	log.Printf("Binding %s.%s to %s...", app.Name, domain.Name, app.Name)
	err = routeRepo.Bind(config, createdRoute, app)
	if err != nil {
		err = errors.New("Error binding route")
		return
	}

	return
}
