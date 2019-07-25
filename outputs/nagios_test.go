package outputs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNagiosDocumentation_Name(t *testing.T) {
	j := Nagios{}
	assert.Equal(t, "nagios", j.Name())
}

func TestNagiosDocumentation_Output(t *testing.T) {
	result, exitCode := runOutput(Nagios{}, GetExampleTestResult())

	expected := `GOSS OK - Count: 1, Failed: 0, Skipped: 0, Duration: 0.000s
`
	assert.Equal(t, expected, result)
	assert.Equal(t, 0, exitCode)
}
