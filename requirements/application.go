package requirements

import (
	"errors"
	"fmt"
	"github.com/diatmpravin/gagan/api"
	"github.com/diatmpravin/gagan/configuration"
	"github.com/diatmpravin/gagan/models"
	"log"
)

type ApplicationRequirement struct {
	name    string
	config  *configuration.Configuration
	appRepo api.ApplicationRepository

	Application models.Application
}

func NewApplicationRequirement(name string, config *configuration.Configuration, aR api.ApplicationRepository) (req ApplicationRequirement) {
	req.name = name
	req.config = config
	req.appRepo = aR
	return
}

func (req *ApplicationRequirement) Execute(config *configuration.Configuration) (err error) {
	req.Application, err = req.appRepo.FindByName(config, req.name)
	if err != nil {
		log.Printf("Request failed: %+v", err)
		err = errors.New(fmt.Sprintf("Request failed: %s", err))
	}
	return
}
