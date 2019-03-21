package kubernetestypes

import (
	"regexp"
	"strings"

	graphapi "github.com/carlsoncoder/audit-webhook-go/graphapi"
)

type User struct {
	UserName string   `json:"username"`
	Groups   []string `json:"groups"`
}

func (u *User) GetUserDetails(userTenantURLPrefix string, graphAPIClient *graphapi.GraphAPIClient) (string, string, string) {
	userObjectID := u.UserName[len(userTenantURLPrefix):]
	userDisplayName, userPrincipalName := graphAPIClient.LoadUserValues(userObjectID)
	return userObjectID, userDisplayName, userPrincipalName
}

func (u *User) GetUserGroups(graphAPIClient *graphapi.GraphAPIClient) string {
	arrayLength := len(u.Groups)
	if arrayLength == 0 {
		return "UNKNOWN"
	}

	var userGroups []string
	userGroups = make([]string, arrayLength, arrayLength)

	for i, group := range u.Groups {
		if isGroupGUID(group) {
			displayName := graphAPIClient.LoadGroupValues(group)
			userGroups[i] = displayName
		} else {
			// just a kubernetes group, not an AAD group
			userGroups[i] = group
		}
	}

	return strings.Join(userGroups, ",")
}

func (u *User) IsUserValidAADUser(userTenantURLPrefix string) bool {
	return strings.HasPrefix(u.UserName, userTenantURLPrefix)
}

func isGroupGUID(group string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(group)
}
