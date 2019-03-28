package auditextensions

import (
	"regexp"
	"strings"

	graphapi "github.com/carlsoncoder/audit-webhook-go/graphapi"
	auditinternal "k8s.io/apiserver/pkg/apis/audit"
)

// IsUserValidAADUser determines if the user is a valid AAD user by checking the full name against the known tenant URL prefix
func IsUserValidAADUser(u *auditinternal.UserInfo, userTenantURLPrefix string) bool {
	return strings.HasPrefix(u.Username, userTenantURLPrefix)
}

// GetUserDetails determines the user details by calling the Graph API
func GetUserDetails(u *auditinternal.UserInfo, userTenantURLPrefix string, graphAPIClient graphapi.Client) (string, string) {
	userObjectID := u.Username[len(userTenantURLPrefix):]
	userDisplayName, userPrincipalName := graphAPIClient.LoadUserValues(userObjectID)
	return userDisplayName, userPrincipalName
}

// GetUserGroups determines the groups the user belongs to by calling the Graph API
// Note that the Graph API is only called if the group is in a GUID format
// If it's not a GUID, we assume it's just a kubernetes group name and we don't call the Graph API
func GetUserGroups(u *auditinternal.UserInfo, graphAPIClient graphapi.Client) string {
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
