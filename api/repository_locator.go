package api

import (
	"github.com/diatmpravin/gagan/configuration"
)

type RepositoryLocator struct {
	config *configuration.Configuration

	appRepo CloudControllerApplicationRepository
}

func NewRepositoryLocator(config *configuration.Configuration) (locator RepositoryLocator) {
	locator.config = config
	locator.appRepo = CloudControllerApplicationRepository{}

	return
}

func (locator RepositoryLocator) GetConfig() *configuration.Configuration {
	return locator.config
}

func (locator RepositoryLocator) GetApplicationRepository() ApplicationRepository {
	return locator.appRepo
}
