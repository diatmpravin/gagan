package api

import (
	"github.com/diatmpravin/gagan/configuration"
)

type RepositoryLocator struct {
	config *configuration.Configuration

	appRepo   CloudControllerApplicationRepository
	spaceRepo CloudControllerSpaceRepository
}

func NewRepositoryLocator(config *configuration.Configuration) (locator RepositoryLocator) {
	locator.config = config
	locator.appRepo = CloudControllerApplicationRepository{}
	locator.spaceRepo = CloudControllerSpaceRepository{}

	return
}

func (locator RepositoryLocator) GetConfig() *configuration.Configuration {
	return locator.config
}

func (locator RepositoryLocator) GetApplicationRepository() ApplicationRepository {
	return locator.appRepo
}

func (locator RepositoryLocator) GetSpaceRepository() SpaceRepository {
	return locator.spaceRepo
}
