package api

import (
	"fmt"
	"github.com/diatmpravin/gagan/configuration"
	"github.com/diatmpravin/gagan/models"
	"github.com/diatmpravin/gagan/testhelpers"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
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

var singleAppResponse = testhelpers.TestResponse{Status: http.StatusOK, Body: `
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
    }
  ]
}`}

var findAppEndpoint = testhelpers.CreateEndpoint(
	"GET",
	"/v2/spaces/my-space-guid/apps?q=name%3AApp1&inline-relations-depth=1",
	nil,
	singleAppResponse,
)

var appSummaryResponse = testhelpers.TestResponse{Status: http.StatusOK, Body: `
{
  "guid": "app1-guid",
  "name": "App1",
  "routes": [
    {
      "guid": "route-1-guid",
      "host": "app1",
      "domain": {
        "guid": "domain-1-guid",
        "name": "cfapps.io"
      }
    }
  ],
  "running_instances": 1,
  "memory": 128,
  "instances": 1
}`}

var appSummaryEndpoint = testhelpers.CreateEndpoint(
	"GET",
	"/v2/apps/app1-guid/summary",
	nil,
	appSummaryResponse,
)

var singleAppEndpoint = func(writer http.ResponseWriter, request *http.Request) {
	if strings.Contains(request.URL.Path, "summary") {
		appSummaryEndpoint(writer, request)
		return
	}

	findAppEndpoint(writer, request)
	return
}

func TestFindByName(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(singleAppEndpoint))
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
	assert.Equal(t, app.Memory, 128)
	assert.Equal(t, app.Instances, 1)

	assert.Equal(t, len(app.Urls), 1)
	assert.Equal(t, app.Urls[0], "app1.cfapps.io")

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

var deleteApplicationEndpoint = testhelpers.CreateEndpoint(
	"DELETE",
	"/v2/apps/my-cool-app-guid?recursive=true",
	nil,
	testhelpers.TestResponse{Status: http.StatusOK, Body: ""},
)

func TestDeleteApplication(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(deleteApplicationEndpoint))
	defer ts.Close()

	repo := CloudControllerApplicationRepository{}
	config := &configuration.Configuration{
		AccessToken: "BEARER my_access_token",
		Target:      ts.URL,
	}

	app := models.Application{Name: "my-cool-app", Guid: "my-cool-app-guid"}

	err := repo.Delete(config, app)
	assert.NoError(t, err)
}

var createApplicationResponse = `
{
    "metadata": {
        "guid": "my-cool-app-guid"
    },
    "entity": {
        "name": "my-cool-app"
    }
}`

var createApplicationEndpoint = testhelpers.CreateEndpoint(
	"POST",
	"/v2/apps",
	testhelpers.RequestBodyMatcher(`{"space_guid":"my-space-guid","name":"my-cool-app","instances":1,"buildpack":null,"command":null,"memory":256,"stack_guid":null}`),
	testhelpers.TestResponse{Status: http.StatusCreated, Body: createApplicationResponse},
)

var alwaysSuccessfulEndpoint = func(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(writer, "{}")
}

func TestCreateApplication(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(createApplicationEndpoint))
	defer ts.Close()

	repo := CloudControllerApplicationRepository{}
	config := &configuration.Configuration{
		AccessToken: "BEARER my_access_token",
		Target:      ts.URL,
		Space:       models.Space{Guid: "my-space-guid"},
	}

	newApp := models.Application{Name: "my-cool-app"}

	createdApp, err := repo.Create(config, newApp)
	assert.NoError(t, err)

	assert.Equal(t, createdApp, models.Application{Name: "my-cool-app", Guid: "my-cool-app-guid"})
}

func TestCreateRejectsInproperNames(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(alwaysSuccessfulEndpoint))
	defer ts.Close()

	config := &configuration.Configuration{Target: ts.URL}
	repo := CloudControllerApplicationRepository{}

	createdApp, err := repo.Create(config, models.Application{Name: "name with space"})
	assert.Equal(t, createdApp, models.Application{})
	assert.Contains(t, err.Error(), "Application name is invalid")

	_, err = repo.Create(config, models.Application{Name: "name-with-inv@lid-chars!"})
	assert.Error(t, err)

	_, err = repo.Create(config, models.Application{Name: "Valid-Name"})
	assert.NoError(t, err)

	_, err = repo.Create(config, models.Application{Name: "name_with_numbers_2"})
	assert.NoError(t, err)
}

var successfulGetInstancesEndpoint = testhelpers.CreateEndpoint(
	"GET",
	"/v2/apps/my-cool-app-guid/instances",
	nil,
	testhelpers.TestResponse{Status: http.StatusCreated, Body: `
{
  "1": {
    "state": "STARTING"
  },
  "0": {
    "state": "RUNNING"
  }
}`},
)

func TestGetInstances(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(successfulGetInstancesEndpoint))
	defer ts.Close()

	repo := CloudControllerApplicationRepository{}
	config := &configuration.Configuration{
		AccessToken: "BEARER my_access_token",
		Target:      ts.URL,
	}

	app := models.Application{Name: "my-cool-app", Guid: "my-cool-app-guid"}

	instances, err := repo.GetInstances(config, app)
	assert.NoError(t, err)
	assert.Equal(t, len(instances), 2)
	assert.Equal(t, instances[0].State, "running")
	assert.Equal(t, instances[1].State, "starting")
}
