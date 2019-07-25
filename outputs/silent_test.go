package outputs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSilentDocumentation_Name(t *testing.T) {
	j := Silent{}
	assert.Equal(t, "silent", j.Name())
}

func TestSilent_Output_Success(t *testing.T) {
	result, exitCode := runOutput(Silent{}, getSuccessTestResult())

	assert.Equal(t, "", result)
	assert.Equal(t, 0, exitCode)
}

func TestSilent_Output_Fail(t *testing.T) {
	result, exitCode := runOutput(Silent{}, getFailTestResult())

	assert.Equal(t, "", result)
	assert.Equal(t, 1, exitCode)
}

func TestSilent_Output_Skip(t *testing.T) {
	result, exitCode := runOutput(Silent{}, getSkipTestResult())

	assert.Equal(t, "", result)
	assert.Equal(t, 0, exitCode)
}
