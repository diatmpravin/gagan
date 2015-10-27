package commands

import (
	"encoding/json"
	"fmt"
	"github.com/diatmpravin/gagan/api"
	"github.com/diatmpravin/gagan/configuration"
	"log"
	"net/http"
)

type Login struct {
	config        *configuration.Configuration
	orgRepo       OrganizationRepository
	spaceRepo     SpaceRepository
	authenticator api.Authenticator
}

type User struct {
	Email    string `json:email`
	Password string `json:password`
}

func (u *User) IsValid() error {
	if u.Email == "" {
		return fmt.Errorf("Email can't be emtpy")
	} else if u.Password == "" {
		return fmt.Errorf("Password can't be emtpy")
	}
	return nil
}

// PutPerson get user token from CC
func (l Login) PutUser(u *User) (config *configuration.Configuration, err error) {
	config = configuration.GetDefaultConfig()

	response, err := api.Authenticate(config.AuthorizationEndpoint, u.Email, u.Password)
	if err != nil {
		return
	}

	config.AccessToken = fmt.Sprintf("%s %s", response.TokenType, response.AccessToken)
	return
}

func SessionPostCase(w http.ResponseWriter, r *http.Request) {
	render := &api.Render{r, w}

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err := u.IsValid(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	l := Login{}
	config, err := l.PutUser(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Configuration: %+v", config)
	render.JSON(config)
}
