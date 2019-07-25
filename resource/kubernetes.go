package resource

import "github.com/SimonBaeumer/goss/system"

type Kubernetes struct {
	Title        string `json:"title,omitempty" yaml:"title,omitempty"`
	Meta         meta   `json:"meta,omitempty" yaml:"meta,omitempty"`
	Name         string `json:"-" yaml:"-"`
	ResourceType string `json:"resource-type" yaml:"resource-type"`
}

func (k *Kubernetes) ID() string {
	return k.Name
}

func (k *Kubernetes) SetID(id string) {
	k.Name = id
}

func (k *Kubernetes) GetTitle() string {
	return k.Title
}

func (k *Kubernetes) GetMeta() meta {
	return k.Meta
}

func (k *Kubernetes) Validate(sys *system.System) []TestResult {
	var results []TestResult

	result := TestResult{
		Successful: true,
		Result:     SUCCESS,
	}

	results = append(results, result)
	return results
}

// TODO: add feature
func NewKubernetes() {

}
