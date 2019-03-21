package graphapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	grantType                     = "client_credentials"
	scope                         = "https://graph.microsoft.com/.default"
	graphAPILoginURLFormat        = "https://login.microsoftonline.com/%s/oauth2/v2.0/token"
	graphAPIUserDetailsURLFormat  = "https://graph.microsoft.com/v1.0/users/%s"
	graphAPIGroupDetailsURLFormat = "https://graph.microsoft.com/v1.0/groups/%s"
	accessTokenFieldName          = "access_token"
	displayNameFieldName          = "displayName"
	userPrincipalNameFieldName    = "userPrincipalName"
)

type GraphAPIClient struct {
	tenantID     string
	clientID     string
	clientSecret string
}

func NewGraphAPIClient(tenantID string, clientID string, clientSecret string) *GraphAPIClient {
	return &GraphAPIClient{tenantID, clientID, clientSecret}
}

func (client *GraphAPIClient) LoadUserValues(userObjectID string) (string, string) {
	accessToken := obtainAccessToken(client)
	displayName, userPrincipalName := determineUserProperties(accessToken, userObjectID)
	return displayName, userPrincipalName
}

func (client *GraphAPIClient) LoadGroupValues(groupObjectID string) string {
	accessToken := obtainAccessToken(client)
	displayName := determineGroupProperties(accessToken, groupObjectID)
	return displayName
}

func obtainAccessToken(client *GraphAPIClient) string {
	graphURL := fmt.Sprintf(graphAPILoginURLFormat, client.tenantID)

	requestBody := url.Values{
		"grant_type":    {grantType},
		"client_secret": {client.clientSecret},
		"client_id":     {client.clientID},
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

func determineGroupProperties(accessToken string, groupObjectID string) string {
	graphURL := fmt.Sprintf(graphAPIGroupDetailsURLFormat, groupObjectID)

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

	graphAPIError := responseValues["error"]
	if graphAPIError != nil {
		// there are some strange GUIDs that show up in the kubernetes audit logs that don't seem to match up with an Azure AD group
		// just log the GUID we tried to use and bail out in this case
		return fmt.Sprintf("Unknown AAD Group ID: %s", groupObjectID)
	}

	return responseValues[displayNameFieldName].(string)
}
