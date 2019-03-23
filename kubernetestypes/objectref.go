package kubernetestypes

// ObjectRef struct maps to the objectRef object in the kubernetes audit record
type ObjectRef struct {
	Resource  string `json:"resource"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}
