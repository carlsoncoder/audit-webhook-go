package graphapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	grantType                    = "client_credentials"
	scope                        = "https://graph.microsoft.com/.default"
	graphAPILoginURLFormat       = "https://login.microsoftonline.com/%s/oauth2/v2.0/token"
	graphAPIUserDetailsURLFormat = "https://graph.microsoft.com/v1.0/users/%s"
	accessTokenFieldName         = "access_token"
	displayNameFieldName         = "displayName"
	userPrincipalNameFieldName   = "userPrincipalName"
)

type UserFinder struct {
	tenantID     string
	clientID     string
	clientSecret string
}

func NewUserFinder(tenantID string, clientID string, clientSecret string) *UserFinder {
	return &UserFinder{tenantID, clientID, clientSecret}
}

func (uf *UserFinder) LoadUserValues(userObjectID string) (string, string) {
	accessToken := obtainAccessToken(uf)
	displayName, userPrincipalName := determineUserProperties(accessToken, userObjectID)
	return displayName, userPrincipalName
}

func obtainAccessToken(uf *UserFinder) string {
	graphURL := fmt.Sprintf(graphAPILoginURLFormat, uf.tenantID)

	requestBody := url.Values{
		"grant_type":    {grantType},
		"client_secret": {uf.clientSecret},
		"client_id":     {uf.clientID},
		"scope":         {scope}}

	response, err := http.PostForm(graphURL, requestBody)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		panic(err)
	}

	var responseValues map[string]interface{}
	err = json.Unmarshal([]byte(body), &responseValues)

	if err != nil {
		panic(err)
	}

	return responseValues[accessTokenFieldName].(string)
}

func determineUserProperties(accessToken string, userObjectID string) (string, string) {
	graphURL := fmt.Sprintf(graphAPIUserDetailsURLFormat, userObjectID)

	client := &http.Client{}
	request, err := http.NewRequest("GET", graphURL, nil)
	if err != nil {
		panic(err)
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		panic(err)
	}

	var responseValues map[string]interface{}
	err = json.Unmarshal([]byte(body), &responseValues)

	if err != nil {
		panic(err)
	}

	return responseValues[displayNameFieldName].(string), responseValues[userPrincipalNameFieldName].(string)
}
