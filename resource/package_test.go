package resource

import (
	"github.com/SimonBaeumer/goss/system"
	"github.com/SimonBaeumer/goss/util"
	"github.com/SimonBaeumer/goss/util/goss_testing"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPackage(t *testing.T) {
	pkg, _ := NewPackage(TestPackage{}, util.Config{})

	assert.Equal(t, "test-pkg-manager", pkg.Name)
	assert.True(t, pkg.Installed.(bool))
}

func TestPackage_Validate(t *testing.T) {
	p := Package{
		Title:          "vim",
		Name:           "vim",
		Installed:      false,
		Versions:       goss_testing.ConvertStringSliceToInterfaceSlice([]string{"1.0.0"}),
		PackageManager: "deb",
	}

	sys := &system.System{NewPackage: func(name string, pkg string) system.Package {
		return TestPackage{}
	}}

	r := p.Validate(sys)

	assert.False(t, r[0].Successful)
	assert.Equal(t, "installed", r[0].Property)

	assert.True(t, r[1].Successful)
	assert.Equal(t, "version", r[1].Property)
}

type TestPackage struct{}

func (p TestPackage) Name() string { return "test-pkg-manager" }

func (p TestPackage) Exists() (bool, error) { return true, nil }

func (p TestPackage) Installed() (bool, error) { return true, nil }

func (p TestPackage) Versions() ([]string, error) { return []string{"1.0.0"}, nil }
