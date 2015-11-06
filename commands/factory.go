package commands

import (
	"github.com/diatmpravin/gagan/api"
)

type Factory struct {
	repoLocator api.RepositoryLocator
}

func NewFactory(repoLocator api.RepositoryLocator) (factory Factory) {
	return
}

func (f Factory) NewStart() *Start {
	return NewStart(
		f.repoLocator.GetConfig(),
		f.repoLocator.GetApplicationRepository(),
	)
}

func (f Factory) NewStop() *Stop {
	return NewStop(
		f.repoLocator.GetConfig(),
		f.repoLocator.GetApplicationRepository(),
	)
}

func (f Factory) NewAllApps() Apps {
	return NewAllApps(
		f.repoLocator.GetConfig(),
		f.repoLocator.GetSpaceRepository(),
	)
}

func (f Factory) NewAppSummary() *App {
	return NewAppSummary(
		f.repoLocator.GetConfig(),
		f.repoLocator.GetApplicationRepository(),
	)
}

func (f Factory) NewDelete() *Delete {
	return NewDelete(
		f.repoLocator.GetConfig(),
		f.repoLocator.GetApplicationRepository(),
	)
}

func (f Factory) NewInstanceInformation() *InstanceInformation {
	return NewInstanceInformation(
		f.repoLocator.GetConfig(),
		f.repoLocator.GetApplicationRepository(),
	)
}

func (f Factory) NewAppEvents() *AppEvents {
	return NewAppEvents(
		f.repoLocator.GetConfig(),
		f.repoLocator.GetApplicationEventRepository(),
	)
}

func (f Factory) NewSpaceList() SpaceList {
	return NewSpaceList(
		f.repoLocator.GetConfig(),
		f.repoLocator.GetSpaceRepository(),
	)
}
