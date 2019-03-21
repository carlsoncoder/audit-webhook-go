package kubernetestypes

// Annotation struct maps to the annotation object in the kubernetes audit record
type Annotation struct {
	Decision string `json:"authorization.k8s.io/decision"`
	Reason   string `json:"authorization.k8s.io/reason"`
}
