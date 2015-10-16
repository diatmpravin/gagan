package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthenticate(t *testing.T) {
	response, err := Authenticate("https://login.run.pivotal.io", "pravinmishra88@gmail.com", "cf@rest12")
	if err != nil {
		t.Fatalf("Unexpect error on authenticate request: %v", err)
	}

	if assert.NotNil(t, response) {
		assert.NotNil(t, response.AccessToken)
	}
}
