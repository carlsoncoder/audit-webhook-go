package auditextensions

import (
	"strings"

	auditinternal "k8s.io/apiserver/pkg/apis/audit"
)

// GetSourceIPAddress gets the source IP address from the object, concantenating multiple records together into one with a '|' delimiter
func GetSourceIPAddress(e *auditinternal.Event) string {
	if len(e.SourceIPs) == 0 {
		return "UNKNOWN"
	}

	return strings.Join(e.SourceIPs, "|")
}
