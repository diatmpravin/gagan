package commands

import (
	"encoding/json"
	"github.com/diatmpravin/gagan/api"
	"github.com/diatmpravin/gagan/configuration"
	"github.com/diatmpravin/gagan/models"
	"net/http"
)

type OrganizationRepository interface {
	FindOrganizations(config *configuration.Configuration) (orgs []models.Organization, err error)
}

type CloudControllerOrganizationRepository struct {
}

func (repo CloudControllerOrganizationRepository) FindOrganizations(config *configuration.Configuration) (orgs []models.Organization, err error) {
	path := config.Target + "/v2/organizations"

	request, err := api.NewAuthorizedRequest("GET", path, config.AccessToken, nil)
	if err != nil {
		return
	}
	response := new(ApiResponse)

	err = api.PerformRequestForBody(request, response)
	if err != nil {
		return
	}

	for _, r := range response.Resources {
		orgs = append(orgs, models.Organization{r.Entity.Name, r.Metadata.Guid})
	}

	return
}

// ListAllOrganizations GET list of all organizations
func ListAllOrganizations(w http.ResponseWriter, r *http.Request) {
	render := &api.Render{r, w}

	config := configuration.GetDefaultConfig()
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	ccorgrep := CloudControllerOrganizationRepository{}

	orgs, err := ccorgrep.FindOrganizations(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(orgs)
}
