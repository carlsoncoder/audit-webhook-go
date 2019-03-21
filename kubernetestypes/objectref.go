package kubernetestypes

// ObjectRef struct maps to the objectRef object in the kubernetes audit record
type ObjectRef struct {
	Resource string `json:"resource"`
	Name     string `json:"name"`
}
