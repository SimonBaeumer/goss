package outputs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSilentDocumentation_Name(t *testing.T) {
	j := Silent{}
	assert.Equal(t, "silent", j.Name())
}

func TestSilentDocumentation_Output(t *testing.T) {
	result, exitCode := runOutput(Silent{}, GetExampleTestResult())

	assert.Equal(t, "", result)
	assert.Equal(t, 0, exitCode)
}
