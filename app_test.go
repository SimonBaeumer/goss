package goss

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func Test_Render_WithoutGossfile(t *testing.T) {
	runtime := GossRunTime{}
	result, err := runtime.Render()

	assert.NotNil(t, err)
	assert.Equal(t, "Could not open gossfile with error: open : no such file or directory", err.Error())
	assert.Empty(t, result)
}

func Test_Render_WithoutVarsfile(t *testing.T) {
	file, err := ioutil.TempFile("", "tmp_gossfile_*.yaml")
	defer os.Remove(file.Name())

	runtime := GossRunTime{
		Gossfile: file.Name(),
		Vars:     "/invalidpath",
	}
	result, err := runtime.Render()

	assert.NotNil(t, err)
	assert.Equal(t, "Could not open varsfile with error: open /invalidpath: no such file or directory", err.Error())
	assert.Empty(t, result)
}
