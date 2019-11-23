package system

import (
    "github.com/SimonBaeumer/goss/util"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestNewAlpinePackage(t *testing.T) {
    sys := System{}
    conf := util.Config{}

    res := NewAlpinePackage("vim", &sys, conf)

    assert.IsType(t, new(AlpinePackage), res)
    assert.Equal(t, "vim", res.Name())
}
