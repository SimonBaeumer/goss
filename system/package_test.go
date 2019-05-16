package system

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestNewPackage(t *testing.T) {
    deb := NewPackage("package", "deb")
    rpm := NewPackage("package", "rpm")
    pac := NewPackage("package", "pacman")
    apk := NewPackage("package", "apk")

    assert.Implements(t, new(Package), deb)
    assert.IsType(t, &DebPackage{}, deb)

    assert.Implements(t, new(Package), rpm)
    assert.IsType(t, &RpmPackage{}, rpm)

    assert.Implements(t, new(Package), pac)
    assert.IsType(t, &PacmanPackage{}, pac)

    assert.Implements(t, new(Package), apk)
    assert.IsType(t, &AlpinePackage{}, apk)
}
