package system

import (
	"errors"

	"github.com/SimonBaeumer/goss/util"
)

type Package interface {
	Name() string
	Exists() (bool, error)
	Installed() (bool, error)
	Versions() ([]string, error)
}

var ErrNullPackage = errors.New("Could not detect Package type on this system, please use --package flag to explicity set it")

type NullPackage struct {
	name string
}

func NewNullPackage(name string, system *System, config util.Config) Package {
	return &NullPackage{name: name}
}

func (p *NullPackage) Name() string { return p.name }

func (p *NullPackage) Exists() (bool, error) { return p.Installed() }

func (p *NullPackage) Installed() (bool, error) {
	return false, ErrNullPackage
}

func (p *NullPackage) Versions() ([]string, error) {
	return nil, ErrNullPackage
}

// DetectPackageManager attempts to detect whether or not the system is using
// "deb", "rpm", "apk", or "pacman" package managers. It first attempts to
// detect the distro. If that fails, it falls back to finding package manager
// executables. If that fails, it returns the empty string.
func DetectPackageManager() string {
	switch DetectDistro() {
	case "ubuntu":
		return "deb"
	case "redhat":
		return "rpm"
	case "alpine":
		return "apk"
	case "arch":
		return "pacman"
	case "debian":
		return "deb"
	}
	for _, manager := range []string{"deb", "rpm", "apk", "pacman"} {
		if HasCommand(manager) {
			return manager
		}
	}
	return ""
}

// NewPackage is the constructor method which creates the correct package manager
// If pkgManager is empty the package manager will be automatically detected
func NewPackage(name string, pkgManager string) Package {
	if pkgManager != "deb" && pkgManager != "apk" && pkgManager != "pacman" && pkgManager != "rpm" {
		pkgManager = DetectPackageManager()
	}
	switch pkgManager {
	case "deb":
		return NewDebPackage(name)
	case "apk":
		return NewAlpinePackage(name)
	case "pacman":
		return NewPacmanPackage(name)
	default:
		return NewRpmPackage(name)
	}
}
