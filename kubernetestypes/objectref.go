package kubernetestypes

type ObjectRef struct {
	Resource string `json:"resource"`
	Name     string `json:"name"`
}
