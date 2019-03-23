package kubernetestypes

// ResponseStatus struct maps to the responseStatus object in the kubernetes audit record
type ResponseStatus struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
	Code   int    `json:"code"`
}
