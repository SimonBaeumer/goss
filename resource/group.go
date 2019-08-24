package resource

import (
	"fmt"

	"github.com/SimonBaeumer/goss/system"
	"github.com/SimonBaeumer/goss/util"
)

// Group represents the group resource
type Group struct {
	Title     string  `json:"title,omitempty" yaml:"title,omitempty"`
	Meta      meta    `json:"meta,omitempty" yaml:"meta,omitempty"`
	Groupname string  `json:"-" yaml:"-"`
	Exists    matcher `json:"exists" yaml:"exists"`
	GID       matcher `json:"gid,omitempty" yaml:"gid,omitempty"`
}

func (g *Group) ID() string      { return g.Groupname }
func (g *Group) SetID(id string) { g.Groupname = id }

func (g *Group) GetTitle() string { return g.Title }
func (g *Group) GetMeta() meta    { return g.Meta }

// Validate validates the group resource
func (g *Group) Validate(sys *system.System) []TestResult {
	skip := false
	sysgroup := sys.NewGroup(g.Groupname, sys, util.Config{})

	var results []TestResult
	results = append(results, ValidateValue(g, "exists", g.Exists, sysgroup.Exists, skip))
	if shouldSkip(results) {
		skip = true
	}
	if g.GID != nil {
		gGID := deprecateAtoI(g.GID, fmt.Sprintf("%s: group.gid", g.Groupname))
		results = append(results, ValidateValue(g, "gid", gGID, sysgroup.GID, skip))
	}
	return results
}

// NewGroup will be used to get the group by the current setting
// Will be used for the add and auto-add features
func NewGroup(sysGroup system.Group, config util.Config) (*Group, error) {
	groupname := sysGroup.Groupname()
	exists, _ := sysGroup.Exists()
	g := &Group{
		Groupname: groupname,
		Exists:    exists,
	}
	if !contains(config.IgnoreList, "stderr") {
		if gid, err := sysGroup.GID(); err == nil {
			g.GID = gid
		}
	}
	return g, nil
}
