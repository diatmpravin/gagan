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

var multipleAppsResponse = testhelpers.TestResponse{Status: http.StatusOK, Body: `
{
  "resources": [
    {
      "metadata": {
        "guid": "app1-guid"
      },
      "entity": {
        "name": "App1",
        "memory": 256,
        "instances": 1,
        "state": "STOPPED",
        "routes": [
      	  {
      	    "metadata": {
      	      "guid": "app1-route-guid"
      	    },
      	    "entity": {
      	      "host": "app1",
      	      "domain": {
      	      	"metadata": {
      	      	  "guid": "domain1-guid"
      	      	},
      	      	"entity": {
      	      	  "name": "cfapps.io"
      	      	}
      	      }
      	    }
      	  }
        ]
      }
    },
    {
      "metadata": {
        "guid": "app2-guid"
      },
      "entity": {
        "name": "App2",
        "memory": 512,
        "instances": 2,
        "state": "STARTED",
        "routes": [
      	  {
      	    "metadata": {
      	      "guid": "app2-route-guid"
      	    },
      	    "entity": {
      	      "host": "app2",
      	      "domain": {
      	      	"metadata": {
      	      	  "guid": "domain1-guid"
      	      	},
      	      	"entity": {
      	      	  "name": "cfapps.io"
      	      	}
      	      }
      	    }
      	  }
        ]
      }
    }
  ]
}`}

var multipleAppsEndpoint = testhelpers.CreateEndpoint(
	"GET",
	"/v2/spaces/my-space-guid/apps?inline-relations-depth=2",
	nil,
	multipleAppsResponse,
)

func TestApplicationsFindAll(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(multipleAppsEndpoint))
	defer ts.Close()

	repo := CloudControllerApplicationRepository{}
	config := &configuration.Configuration{
		AccessToken: "BEARER my_access_token",
		Target:      ts.URL,
		Space:       models.Space{Name: "my-space", Guid: "my-space-guid"},
	}

	apps, err := repo.FindApps(config)
	assert.NoError(t, err)
	assert.Equal(t, len(apps), 2)

	app := apps[0]
	assert.Equal(t, app.Name, "App1")
	assert.Equal(t, app.Guid, "app1-guid")
	assert.Equal(t, app.State, "stopped")
	assert.Equal(t, app.Instances, 1)
	assert.Equal(t, app.Memory, 256)
	assert.Equal(t, len(app.Urls), 1)
	assert.Equal(t, app.Urls[0], "app1.cfapps.io")

	app = apps[1]
	assert.Equal(t, app.Guid, "app2-guid")
}

func TestFindByName(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(multipleAppsEndpoint))
	defer ts.Close()

	repo := CloudControllerApplicationRepository{}
	config := &configuration.Configuration{
		AccessToken: "BEARER my_access_token",
		Target:      ts.URL,
		Space:       models.Space{Name: "my-space", Guid: "my-space-guid"},
	}

	app, err := repo.FindByName(config, "App1")
	assert.NoError(t, err)
	assert.Equal(t, app.Name, "App1")
	assert.Equal(t, app.Guid, "app1-guid")

	app, err = repo.FindByName(config, "app1")
	assert.NoError(t, err)
	assert.Equal(t, app.Guid, "app1-guid")

	app, err = repo.FindByName(config, "app that does not exist")
	assert.Error(t, err)
}

var startApplicationEndpoint = testhelpers.CreateEndpoint(
	"PUT",
	"/v2/apps/my-cool-app-guid",
	testhelpers.RequestBodyMatcher(`{"console":true,"state":"STARTED"}`),
	testhelpers.TestResponse{Status: http.StatusCreated, Body: `
{
  "metadata": {
    "guid": "my-cool-app-guid",
  },
  "entity": {
    "name": "cli1",
    "state": "STARTED"
  }
}`},
)

func TestStartApplication(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(startApplicationEndpoint))
	defer ts.Close()

	repo := CloudControllerApplicationRepository{}
	config := &configuration.Configuration{
		AccessToken: "BEARER my_access_token",
		Target:      ts.URL,
	}

	app := models.Application{Name: "my-cool-app", Guid: "my-cool-app-guid"}

	err := repo.Start(config, app)
	assert.NoError(t, err)
}

var stopApplicationEndpoint = testhelpers.CreateEndpoint(
	"PUT",
	"/v2/apps/my-cool-app-guid",
	testhelpers.RequestBodyMatcher(`{"console":true,"state":"STOPPED"}`),
	testhelpers.TestResponse{Status: http.StatusCreated, Body: `
{
  "metadata": {
    "guid": "my-cool-app-guid",
  },
  "entity": {
    "name": "cli1",
    "state": "STOPPED"
  }
}`},
)

func TestStopApplication(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(stopApplicationEndpoint))
	defer ts.Close()

	repo := CloudControllerApplicationRepository{}
	config := &configuration.Configuration{
		AccessToken: "BEARER my_access_token",
		Target:      ts.URL,
	}

	app := models.Application{Name: "my-cool-app", Guid: "my-cool-app-guid"}

	err := repo.Stop(config, app)
	assert.NoError(t, err)
}
