package testhelpers

import (
	"errors"
	"github.com/diatmpravin/gagan/configuration"
	"github.com/diatmpravin/gagan/models"
)

type FakeApplicationRepository struct {
	StartedApp  models.Application
	StartAppErr bool

	StoppedApp models.Application
	StopAppErr bool

	FindAllApps []models.Application

	AppByName    models.Application
	AppName      string
	AppByNameErr bool

	CreatedApp models.Application
	DeletedApp models.Application

	GetInstancesResponses  [][]models.ApplicationInstance
	GetInstancesErrorCodes []int
}

func (repo *FakeApplicationRepository) FindAll(config *configuration.Configuration) (apps []models.Application, err error) {
	return repo.FindAllApps, err
}

func (repo *FakeApplicationRepository) FindApps(config *configuration.Configuration) (apps []models.Application, err error) {
	return repo.FindAllApps, err
}

func (repo *FakeApplicationRepository) FindByName(config *configuration.Configuration, name string) (app models.Application, err error) {
	repo.AppName = name
	if repo.AppByNameErr {
		err = errors.New("Error finding app by name.")
	}
	return repo.AppByName, err
}
func (repo *FakeApplicationRepository) Create(config *configuration.Configuration, newApp models.Application) (createdApp models.Application, err error) {
	repo.CreatedApp = newApp

	createdApp = models.Application{
		Name: newApp.Name,
		Guid: newApp.Name + "-guid",
	}

	return
}

func (repo *FakeApplicationRepository) Delete(config *configuration.Configuration, app models.Application) (err error) {
	repo.DeletedApp = app
	return
}

func (repo *FakeApplicationRepository) GetInstances(config *configuration.Configuration, app models.Application) (instances []models.ApplicationInstance, err error) {
	errorCode := repo.GetInstancesErrorCodes[0]
	repo.GetInstancesErrorCodes = repo.GetInstancesErrorCodes[1:]

	instances = repo.GetInstancesResponses[0]
	repo.GetInstancesResponses = repo.GetInstancesResponses[1:]

	if errorCode != 0 {
		err = errors.New("Error while starting app")
		return
	}

	return
}

func (repo *FakeApplicationRepository) Start(config *configuration.Configuration, app models.Application) (err error) {
	repo.StartedApp = app
	if repo.StartAppErr {
		err = errors.New("Error starting app.")
	}
	return
}

func (repo *FakeApplicationRepository) Stop(config *configuration.Configuration, app models.Application) (err error) {
	repo.StoppedApp = app
	if repo.StopAppErr {
		err = errors.New("Error stopping app.")
	}
	return
}
