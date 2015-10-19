package api

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
)

type AuthenticationResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func Authenticate(endpoint string, email string, password string) (response AuthenticationResponse, err error) {
	data := url.Values{
		"username":   {email},
		"password":   {password},
		"grant_type": {"password"},
		"scope":      {""},
	}

	request, err := http.NewRequest("POST", endpoint+"/oauth/token", strings.NewReader(data.Encode()))
	if err != nil {
		return
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("cf:")))

	err = PerformRequestAndParseResponse(request, &response)

	return
}
