package outputs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNagios_Name(t *testing.T) {
	j := Nagios{}
	assert.Equal(t, "nagios", j.Name())
}

func TestNagiosOuput_Success(t *testing.T) {
	result, exitCode := runOutput(Nagios{}, getSuccessTestResult())

	expected := `GOSS OK - Count: 1, Failed: 0, Skipped: 0, Duration: 0.000s
`
	assert.Equal(t, expected, result)
	assert.Equal(t, 0, exitCode)
}

func TestNagiosOuput_Fail(t *testing.T) {
	result, exitCode := runOutput(Nagios{}, getFailTestResult())

	expected := `GOSS CRITICAL - Count: 1, Failed: 1, Skipped: 0, Duration: 0.000s
`
	assert.Equal(t, expected, result)
	assert.Equal(t, 2, exitCode)
}

func TestNagiosOuput_Skip(t *testing.T) {
	result, exitCode := runOutput(Nagios{}, getSkipTestResult())

	expected := `GOSS OK - Count: 1, Failed: 0, Skipped: 1, Duration: 0.000s
`
	assert.Equal(t, expected, result)
	assert.Equal(t, 0, exitCode)
}
