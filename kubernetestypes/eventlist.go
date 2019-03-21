package kubernetestypes

// EventList struct maps the the entire POST body that the kubernetes audit webhook sends
type EventList struct {
	Events []Event `json:"items"`
}
