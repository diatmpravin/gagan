package api

import (
	"encoding/json"
	"fmt"
	"github.com/diatmpravin/gagan/configuration"
	"github.com/diatmpravin/gagan/models"
	"log"
	"net/http"
)

type SpaceRepository interface {
	FindAllSpaces(config *configuration.Configuration) (spaces []models.Space, err error)
}

type CloudControllerSpaceRepository struct {
}

func (repo CloudControllerSpaceRepository) FindAllSpaces(config *configuration.Configuration) (spaces []models.Space, err error) {
	path := fmt.Sprintf("%s/v2/organizations/%s/spaces", config.Target, config.Organization.Guid)
	request, err := NewAuthorizedRequest("GET", path, config.AccessToken, nil)
	if err != nil {
		return
	}

	response := new(ApiResponse)

	_, err = PerformRequestAndParseResponse(request, response)

	if err != nil {
		return
	}

	for _, r := range response.Resources {
		spaces = append(spaces, models.Space{r.Entity.Name, r.Metadata.Guid})
	}

	return
}

// ListAllSpaces Get list of all spaces
func ListAllSpaces(w http.ResponseWriter, r *http.Request) {
	render := &Render{r, w}

	config := configuration.GetDefaultConfig()
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	repo := CloudControllerSpaceRepository{}
	spaces, err := repo.FindAllSpaces(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("List of spaces: %+v", spaces)
	render.JSON(spaces)
}
