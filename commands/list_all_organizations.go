package commands

import (
	"encoding/json"
	"github.com/diatmpravin/gagan/api"
	"github.com/diatmpravin/gagan/configuration"
	"github.com/diatmpravin/gagan/requirements"
	"log"
	"net/http"
)

type OrganizationList struct {
	config           *configuration.Configuration
	organizationRepo api.OrganizationRepository
}

func NewOrganizationList(config *configuration.Configuration, organizationRepo api.OrganizationRepository) (o OrganizationList) {
	o.config = config
	o.organizationRepo = organizationRepo

	return
}

func (o OrganizationList) GetRequirements(reqFactory requirements.Factory, w http.ResponseWriter, r *http.Request) (reqs []Requirement, config *configuration.Configuration, err error) {
	return
}

func (o OrganizationList) Run(w http.ResponseWriter, r *http.Request) {
	render := &api.Render{r, w}

	config := configuration.GetDefaultConfig()
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	orgs, err := o.organizationRepo.FindOrganizations(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("List of organizations: %+v", orgs)
	render.JSON(orgs)
}
