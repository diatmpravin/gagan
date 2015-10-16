package commands

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPutUser(t *testing.T) {
	u := User{
		Email:    "pravinmishra88@gmail.com",
		Password: "cf@rest12",
	}

	config, err := PutUser(&u)
	if err != nil {
		t.Fatalf("Unexpect error on authenticating user: %v", err)
	}

	if assert.NotNil(t, config) {
		assert.Equal(t, config.Target, "https://api.run.pivotal.io")
		assert.Equal(t, config.ApiVersion, "2")
		assert.Equal(t, config.AuthorizationEndpoint, "https://login.run.pivotal.io")
		assert.NotNil(t, config.AccessToken)
		assert.Contains(t, config.AccessToken, "bearer")
	}
}
