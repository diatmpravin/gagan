package testhelpers

import (
	"github.com/diatmpravin/gagan/configuration"
	"github.com/diatmpravin/gagan/models"
)

type FakeServiceRepo struct {
	ServiceOfferings          []models.ServiceOffering
	CreateServiceInstanceName string
	CreateServiceInstancePlan models.ServicePlan
}

func (repo *FakeServiceRepo) GetServiceOfferings(config *configuration.Configuration) (offerings []models.ServiceOffering, err error) {
	offerings = repo.ServiceOfferings
	return
}

func (repo *FakeServiceRepo) CreateServiceInstance(config *configuration.Configuration, name string, plan models.ServicePlan) (err error) {
	repo.CreateServiceInstanceName = name
	repo.CreateServiceInstancePlan = plan
	return
}
