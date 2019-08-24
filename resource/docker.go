package resource

import (
	"github.com/SimonBaeumer/goss/system"
	"github.com/SimonBaeumer/goss/util"
)

// Docker represents the docker resource
type Docker struct {
	Image string `json:"-" yaml:"-"`
	Title string `json:"title,omitempty" yaml:"title,omitempty"`
	Count int    `json:"count,omitempty" yaml:"count,omitempty"`
	Meta  meta   `json:"meta,omitempty" yaml:"meta,omitempty"`
}

// ID returns the id of the resource
func (d *Docker) ID() string {
	return d.Image
}

// SetID sets the id of the resource
func (d *Docker) SetID(id string) {
	d.Image = id
}

func (d *Docker) GetTitle() string {
	return d.Title
}

func (d *Docker) GetMeta() meta {
	return d.GetMeta()
}

// Validate validates the docker resource
func (d *Docker) Validate(sys *system.System) []TestResult {
	skip := false
	sysDocker := sys.NewDocker(d.Image, sys, util.Config{})

	var results []TestResult
	results = append(results, ValidateValue(d, "image", d.Image, sysDocker.Image, skip))
	return results
}

// TODO: implement for add feature
func NewDocker(sysDocker system.Docker) (*Docker, error) {
	return &Docker{}, nil
}
