package testhelpers

import (
	"github.com/diatmpravin/gagan/configuration"
	"github.com/diatmpravin/gagan/models"
)

type FakeApplicationRepository struct {
	FindAllApps []models.Application
}

func (repo *FakeApplicationRepository) FindAll(config *configuration.Configuration) (apps []models.Application, err error) {
	return repo.FindAllApps, err
}
