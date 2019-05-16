package system

import (
	"errors"
	"strings"

	"github.com/SimonBaeumer/goss/util"
)

//PackmanPackage represents a package inside the pacman manager
type PacmanPackage struct {
	name      string
	versions  []string
	loaded    bool
	installed bool
}

//NewPacmanPackage creates a new pacman manager
func NewPacmanPackage(name string) Package {
	return &PacmanPackage{name: name}
}

func (p *PacmanPackage) setup() {
	if p.loaded {
		return
	}
	p.loaded = true
	// TODO: extract versions
	cmd := util.NewCommand("pacman", "-Q", "--color", "never", "--noconfirm", p.name)
	if err := cmd.Run(); err != nil {
		return
	}
	p.installed = true
	// the output format is "pkgname version\n", so if we split the string on
	// whitespace, the version is the second item.
	p.versions = []string{strings.Fields(cmd.Stdout.String())[1]}
}

// Name returns the name of the package
func (p *PacmanPackage) Name() string {
	return p.name
}

// Exists returns if the package is installed
func (p *PacmanPackage) Exists() (bool, error) { return p.Installed() }

// Installed will check and returns if the package is installed
func (p *PacmanPackage) Installed() (bool, error) {
	p.setup()

	return p.installed, nil
}

// Versions returns all installed versions of the package
func (p *PacmanPackage) Versions() ([]string, error) {
	p.setup()
	if len(p.versions) == 0 {
		return p.versions, errors.New("Package version not found")
	}
	return p.versions, nil
}
