package kubernetestypes

type EventList struct {
	Events []Event `json:"items"`
}
