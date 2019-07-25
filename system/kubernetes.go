package system

type Kubernetes interface {
	Status() string
	ResourceType() string
}

type DefKubernetes struct {
	pod string
}

func (k *DefKubernetes) setup() error {
	return nil
}

func Status() string {
	return "pending"
}

func ResourceType() string {
	return "pod"
}
