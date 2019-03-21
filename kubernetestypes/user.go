package kubernetestypes

import (
	"regexp"
	"strings"

	graphapi "github.com/carlsoncoder/audit-webhook-go/graphapi"
)

// User struct maps to the user object in the kubernetes audit record
type User struct {
	UserName string   `json:"username"`
	Groups   []string `json:"groups"`
}

// IsUserValidAADUser determines if the user is a valid AAD user by checking the full name against the known tenant URL prefix
func (u *User) IsUserValidAADUser(userTenantURLPrefix string) bool {
	return strings.HasPrefix(u.UserName, userTenantURLPrefix)
}

// GetUserDetails determines the user details by calling the Graph API
func (u *User) GetUserDetails(userTenantURLPrefix string, graphAPIClient *graphapi.Client) (string, string, string) {
	userObjectID := u.UserName[len(userTenantURLPrefix):]
	userDisplayName, userPrincipalName := graphAPIClient.LoadUserValues(userObjectID)
	return userObjectID, userDisplayName, userPrincipalName
}

// GetUserGroups determines the groups the user belongs to by calling the Graph API
// Note that the Graph API is only called if the group is in a GUID format
// If it's not a GUID, we assume it's just a kubernetes group name and we don't call the Graph API
func (u *User) GetUserGroups(graphAPIClient *graphapi.Client) string {
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

	return strings.Join(userGroups, "|")
}

func isGroupGUID(group string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(group)
}
