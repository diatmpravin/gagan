package commands

import (
	"github.com/diatmpravin/gagan/configuration"
	"github.com/diatmpravin/gagan/requirements"
	"net/http"
)

type Runner struct {
	reqFactory requirements.Factory
}

func NewRunner(reqFactory requirements.Factory) (runner Runner) {
	runner.reqFactory = reqFactory
	return
}

type Command interface {
	GetRequirements(reqFactory requirements.Factory, w http.ResponseWriter, r *http.Request) (reqs []Requirement, config *configuration.Configuration, err error)
	Run(w http.ResponseWriter, r *http.Request)
}

type Requirement interface {
	Execute(config *configuration.Configuration) (err error)
}

func (runner Runner) Run(w http.ResponseWriter, r *http.Request, cmd Command) (err error) {
	requirements, config, err := cmd.GetRequirements(runner.reqFactory, w, r)
	if err != nil {
		return
	}

	for _, requirement := range requirements {
		err = requirement.Execute(config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	cmd.Run(w, r)
	return
}
