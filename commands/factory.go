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
