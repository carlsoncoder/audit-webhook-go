package kubernetestypes

import (
	"strings"
	"time"
)

type Event struct {
	StageTimestamp time.Time  `json:"stageTimestamp"`
	Level          string     `json:"level"`
	Stage          string     `json:"stage"`
	RequestURI     string     `json:"requestURI"`
	Verb           string     `json:"verb"`
	User           User       `json:"user"`
	SourceIPs      []string   `json:"sourceIPs"`
	ObjectRef      ObjectRef  `json:"objectRef"`
	Annotations    Annotation `json:"annotations"`
}

func (e *Event) GetSourceIPAddress() string {
	if len(e.SourceIPs) == 0 {
		return "UNKNOWN"
	} else if len(e.SourceIPs) == 1 {
		return e.SourceIPs[0]
	} else {
		return strings.Join(e.SourceIPs, ",")
	}
}
