package configuration

import (
	"github.com/diatmpravin/gagan/models"
)

type Configuration struct {
	Target                string
	ApiVersion            string
	AuthorizationEndpoint string
	AccessToken           string
	Organization          models.Organization
	Space                 models.Space
}

func GetDefaultConfig() (c *Configuration) {
	c = new(Configuration)
	c.Target = "https://api.run.pivotal.io"
	c.ApiVersion = "2"
	c.AuthorizationEndpoint = "https://login.run.pivotal.io"
	return
}
