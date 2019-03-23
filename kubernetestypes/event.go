package kubernetestypes

import (
	"strings"
	"time"
)

// Event struct maps to the main audit object in the kubernetes audit records
type Event struct {
	StageTimestamp time.Time      `json:"stageTimestamp"`
	Level          string         `json:"level"`
	Stage          string         `json:"stage"`
	RequestURI     string         `json:"requestURI"`
	Verb           string         `json:"verb"`
	User           User           `json:"user"`
	SourceIPs      []string       `json:"sourceIPs"`
	UserAgent      string         `json:"userAgent"`
	ObjectRef      ObjectRef      `json:"objectRef"`
	ResponseStatus ResponseStatus `json:"responseStatus"`
	Annotations    Annotation     `json:"annotations"`
}

// GetSourceIPAddress gets the source IP address from the object, concantenating multiple records together into one with a ':' delimiter
func (e *Event) GetSourceIPAddress() string {
	if len(e.SourceIPs) == 0 {
		return "UNKNOWN"
	}

	return strings.Join(e.SourceIPs, "|")
}
