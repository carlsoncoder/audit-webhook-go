package graphapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
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

// Client interface is used to make calls to the Microsoft Graph API to load user and group details out of Azure AD
type Client interface {
	LoadUserValues(string) (string, string)
	LoadGroupValues(string) string
}

type client struct {
	tenantID         string
	clientID         string
	clientSecret     string
	userLookupTable  map[string]string
	groupLookupTable map[string]string
}

// NewClient instantiates a new instance of the Client struct
func NewClient(tenantID string, clientID string, clientSecret string) Client {
	return &client{tenantID, clientID, clientSecret, make(map[string]string), make(map[string]string)}
}

// LoadUserValues loads the user displayName and user userPrincipalName out of Azure AD via the Graph API
func (client *client) LoadUserValues(userObjectID string) (string, string) {
	// check the cache first - we don't bother invalidating cache ever since Azure User ID's aren't going to change
	user, ok := client.userLookupTable[userObjectID]
	if ok {
		// user object is in the format DisplayName:UserPrincipalName
		userValues := strings.Split(user, ":")
		return userValues[0], userValues[1]
	}

	accessToken := obtainAccessToken(client)
	displayName, userPrincipalName := determineUserProperties(accessToken, userObjectID)

	// add it to the cache
	client.userLookupTable[userObjectID] = fmt.Sprintf("%s:%s", displayName, userPrincipalName)
	return displayName, userPrincipalName
}

// LoadGroupValues loads the group displayName out of Azure AD via the Graph API
func (client *client) LoadGroupValues(groupObjectID string) string {
	// check the cache first - we don't bother invalidating cache ever since Azure Group ID's aren't going to change
	group, ok := client.groupLookupTable[groupObjectID]
	if ok {
		return group
	}

	accessToken := obtainAccessToken(client)
	displayName := determineGroupProperties(accessToken, groupObjectID)

	// add it to the cache
	client.groupLookupTable[groupObjectID] = displayName
	return displayName
}

func obtainAccessToken(client *client) string {
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
	responseValues := callGraphAPIGet(accessToken, graphAPIUserDetailsURLFormat, userObjectID)
	displayName := responseValues[displayNameFieldName].(string)
	userPrincipalName := responseValues[userPrincipalNameFieldName].(string)
	return displayName, userPrincipalName
}

func determineGroupProperties(accessToken string, groupObjectID string) string {
	responseValues := callGraphAPIGet(accessToken, graphAPIGroupDetailsURLFormat, groupObjectID)
	graphAPIError := responseValues["error"]
	if graphAPIError != nil {
		// there are some strange GUIDs that show up in the kubernetes audit logs that don't seem to match up with an Azure AD group
		// just log the GUID we tried to use and bail out in this case
		return fmt.Sprintf("Unknown AAD Group ID: %s", groupObjectID)
	}

	return responseValues[displayNameFieldName].(string)
}

func callGraphAPIGet(accessToken string, urlFormat string, objectID string) map[string]interface{} {
	graphURL := fmt.Sprintf(urlFormat, objectID)

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

	return responseValues
}
