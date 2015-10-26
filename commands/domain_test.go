package commands

import (
	"github.com/diatmpravin/gagan/configuration"
	"github.com/diatmpravin/gagan/models"
	"github.com/diatmpravin/gagan/testhelpers"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var multipleDomainsResponse = testhelpers.TestResponse{Status: http.StatusOK, Body: `
{
  "total_results": 2,
  "total_pages": 1,
  "prev_url": null,
  "next_url": null,
  "resources": [
    {
      "metadata": {
        "guid": "domain1-guid"
      },
      "entity": {
        "name": "domain1.cf-app.com"
      }
    },
    {
      "metadata": {
        "guid": "domain2-guid"
      },
      "entity": {
        "name": "domain2.cf-app.com"
      }
    }
  ]
}`}

var multipleDomainsEndpoint = testhelpers.CreateEndpoint(
	"GET",
	"/v2/spaces/my-space-guid/domains",
	nil,
	multipleDomainsResponse,
)

func TestFindAll(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(multipleDomainsEndpoint))
	defer ts.Close()

	repo := CloudControllerDomainRepository{}

	config := &configuration.Configuration{
		AccessToken: "BEARER my_access_token",
		Target:      ts.URL,
		Space:       models.Space{Guid: "my-space-guid"},
	}

	domains, err := repo.FindAll(config)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(domains))

	first := domains[0]
	assert.Equal(t, first.Name, "domain1.cf-app.com")
	assert.Equal(t, first.Guid, "domain1-guid")
	second := domains[1]
	assert.Equal(t, second.Name, "domain2.cf-app.com")
	assert.Equal(t, second.Guid, "domain2-guid")
}
