package system

import "github.com/SimonBaeumer/goss/util"

type SnapPackage struct {
    name string
    versions []string
    installed bool
}

func NewSnapPackage(name string, system *System, config util.Config) *SnapPackage {
    return &SnapPackage{name: name}
}

func (p *SnapPackage) setup() {

}

func (p *SnapPackage) Name() string {
    return p.name
}

func (p *SnapPackage) Exists() (bool, error) {
    return false, nil
}

func (p *SnapPackage) Installed() (bool, error) {
    p.setup()

    return p.installed, nil
}