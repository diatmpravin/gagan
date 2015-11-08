package requirements

import (
	"github.com/diatmpravin/gagan/configuration"
	"github.com/diatmpravin/gagan/models"
	"github.com/diatmpravin/gagan/testhelpers"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApplicationReqExecute(t *testing.T) {
	app := models.Application{Name: "my-app", Guid: "my-app-guid"}
	appRepo := &testhelpers.FakeApplicationRepository{AppByName: app}
	config := &configuration.Configuration{}

	appReq := NewApplicationRequirement("foo", config, appRepo)
	err := appReq.Execute(config)

	assert.NoError(t, err)
	assert.Equal(t, appRepo.AppName, "foo")
	assert.Equal(t, appReq.Application, app)
}
