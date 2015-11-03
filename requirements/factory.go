package requirements

import (
	"github.com/diatmpravin/gagan/api"
)

type Factory interface {
	NewApplicationRequirement(name string) ApplicationRequirement
}

type ApiRequirementFactory struct {
	repoLocator api.RepositoryLocator
}

func NewFactory(repoLocator api.RepositoryLocator) (factory ApiRequirementFactory) {
	return ApiRequirementFactory{repoLocator}
}

func (f ApiRequirementFactory) NewApplicationRequirement(name string) ApplicationRequirement {
	return NewApplicationRequirement(
		name,
		f.repoLocator.GetConfig(),
		f.repoLocator.GetApplicationRepository(),
	)
}
