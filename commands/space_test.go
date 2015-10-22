package commands

import (
	"fmt"
	"github.com/diatmpravin/gagan/configuration"
	"github.com/diatmpravin/gagan/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var multipleSpacesEndpoint = func(writer http.ResponseWriter, request *http.Request) {
	acceptHeaderMatches := request.Header.Get("accept") == "application/json"
	methodMatches := request.Method == "GET"
	pathMatches := request.URL.Path == "/v2/organizations/some-org-guid/spaces"
	authMatches := request.Header.Get("authorization") == "BEARER my_access_token"

	if !(acceptHeaderMatches && methodMatches && pathMatches && authMatches) {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonResponse := `
{
  "resources": [
    {
      "metadata": {
        "guid": "acceptance-space-guid"
      },
      "entity": {
        "name": "acceptance"
      }
    },
    {
      "metadata": {
        "guid": "staging-space-guid"
      },
      "entity": {
        "name": "staging"
      }
    }
  ]
}`
	fmt.Fprintln(writer, jsonResponse)
}

func TestFindSpaces(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(multipleSpacesEndpoint))
	defer ts.Close()

	repo := CloudControllerSpaceRepository{}
	config := &configuration.Configuration{
		AccessToken:  "BEARER my_access_token",
		Target:       ts.URL,
		Organization: models.Organization{Guid: "some-org-guid"},
	}
	spaces, err := repo.FindAllSpaces(config)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(spaces))

	firstSpace := spaces[0]
	assert.Equal(t, firstSpace.Name, "acceptance")
	assert.Equal(t, firstSpace.Guid, "acceptance-space-guid")

	secondSpace := spaces[1]
	assert.Equal(t, secondSpace.Name, "staging")
	assert.Equal(t, secondSpace.Guid, "staging-space-guid")
}

func TestFindSpacesWithIncorrectToken(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(multipleSpacesEndpoint))
	defer ts.Close()

	repo := CloudControllerSpaceRepository{}

	config := &configuration.Configuration{
		AccessToken:  "BEARER incorrect_access_token",
		Target:       ts.URL,
		Organization: models.Organization{Guid: "some-org-guid"},
	}
	spaces, err := repo.FindAllSpaces(config)

	assert.Error(t, err)
	assert.Equal(t, 0, len(spaces))
}