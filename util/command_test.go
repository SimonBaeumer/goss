package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCommand(t *testing.T) {
	cmd := NewCommand("/bin/sh")
	assert.Equal(t, "/bin/sh", cmd.name)
	assert.Equal(t, "", cmd.Stdout.String())
}

func TestCommand_Run(t *testing.T) {
	cmd := NewCommand("/bin/sh", "-c", "echo test")
	err := cmd.Run()

	assert.Nil(t, err)
	assert.Equal(t, "test\n", cmd.Stdout.String())
	assert.Equal(t, "", cmd.Stderr.String())
	assert.Equal(t, 0, cmd.Status)
}
