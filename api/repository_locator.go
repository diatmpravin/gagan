package api

import (
	"github.com/diatmpravin/gagan/configuration"
)

type RepositoryLocator struct {
	config *configuration.Configuration

	appRepo          CloudControllerApplicationRepository
	organizationRepo CloudControllerOrganizationRepository
	spaceRepo        CloudControllerSpaceRepository
	appEventRepo     CloudControllerAppEventsRepository
}

func NewRepositoryLocator(config *configuration.Configuration) (locator RepositoryLocator) {
	locator.config = config

	locator.appRepo = CloudControllerApplicationRepository{}
	locator.organizationRepo = CloudControllerOrganizationRepository{}
	locator.spaceRepo = CloudControllerSpaceRepository{}
	locator.appEventRepo = CloudControllerAppEventsRepository{}

	return
}

func (locator RepositoryLocator) GetConfig() *configuration.Configuration {
	return locator.config
}

func (locator RepositoryLocator) GetApplicationRepository() ApplicationRepository {
	return locator.appRepo
}

func (locator RepositoryLocator) GetOrganizationRepository() OrganizationRepository {
	return locator.organizationRepo
}

func (locator RepositoryLocator) GetSpaceRepository() SpaceRepository {
	return locator.spaceRepo
}

func (locator RepositoryLocator) GetApplicationEventRepository() AppEventsRepository {
	return locator.appEventRepo
}
