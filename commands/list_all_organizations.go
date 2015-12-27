package commands

import (
	"encoding/json"
	"github.com/diatmpravin/gagan/api"
	"github.com/diatmpravin/gagan/configuration"
	"github.com/diatmpravin/gagan/requirements"
	"log"
	"net/http"
	"strconv"
)

type OrganizationList struct {
	config           *configuration.Configuration
	organizationRepo api.OrganizationRepository
}

func NewOrganizationList(config *configuration.Configuration, organizationRepo api.OrganizationRepository) (o OrganizationList) {
	o.config = config
	o.organizationRepo = organizationRepo

	return
}

func (o OrganizationList) GetRequirements(reqFactory requirements.Factory, w http.ResponseWriter, r *http.Request) (reqs []Requirement, config *configuration.Configuration, err error) {
	return
}

func (o OrganizationList) Run(w http.ResponseWriter, r *http.Request) {
	render := &api.Render{r, w}

	config := configuration.GetDefaultConfig()

	session := configuration.Session{}

	if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	c := configuration.RedisConnect()
	defer c.Close()

	reply, err := c.Do("GET", "user:"+strconv.Itoa(session.SessionId))

	configuration.HandleError(err)

	if err = json.Unmarshal(reply.([]byte), &session); err != nil {
		panic(err)
	}

	config.AccessToken = session.AccessToken

	orgs, err := o.organizationRepo.FindOrganizations(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("List of organizations: %+v", orgs)
	render.JSON(orgs)
}
