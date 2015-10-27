package api

import (
	"fmt"
	"github.com/diatmpravin/gagan/configuration"
	"github.com/diatmpravin/gagan/models"
	"log"
)

type DomainRepository interface {
	FindAll(config *configuration.Configuration) (domains []models.Domain, err error)
}

type CloudControllerDomainRepository struct {
}

func (repo CloudControllerDomainRepository) FindAll(config *configuration.Configuration) (domains []models.Domain, err error) {
	path := fmt.Sprintf("%s/v2/spaces/%s/domains", config.Target, config.Space.Guid)
	request, err := NewAuthorizedRequest("GET", path, config.AccessToken, nil)
	if err != nil {
		return
	}

	response := new(ApiResponse)
	err = PerformRequestAndParseResponse(request, response)
	if err != nil {
		return
	}

	for _, r := range response.Resources {
		domains = append(domains, models.Domain{r.Entity.Name, r.Metadata.Guid})
	}

	log.Printf("App Domain: %+v", domains)
	return
}
