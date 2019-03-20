package kubernetestypes

type Annotation struct {
	Decision string `json:"authorization.k8s.io/decision"`
	Reason   string `json:"authorization.k8s.io/reason"`
}
