package outputs

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRspecish_Name(t *testing.T) {
	j := Rspecish{}
	assert.Equal(t, "rspecish", j.Name())
}

func TestRspecish_Output_Success(t *testing.T) {
	d, _ := time.ParseDuration("2s")

	result, exitCode := runOutput(Rspecish{FakeDuration: d}, getSuccessTestResult())

	expected := `.

Total Duration: 2.000s
Count: 1, Failed: 0, Skipped: 0
`
	assert.Equal(t, expected, result)
	assert.Equal(t, 0, exitCode)
}

func TestRspecish_Output_Fail(t *testing.T) {
	d, _ := time.ParseDuration("2s")

	result, exitCode := runOutput(Rspecish{FakeDuration: d}, getFailTestResult())

	expected := `F

Failures/Skipped:

Title: failure
resource type: my resource id: a property: doesn't match, expect: [expected] found: []

Total Duration: 2.000s
Count: 1, Failed: 1, Skipped: 0
`
	assert.Equal(t, expected, result)
	assert.Equal(t, 1, exitCode)
}

func TestRspecish_Output_Skip(t *testing.T) {
	d, _ := time.ParseDuration("2s")

	result, exitCode := runOutput(Rspecish{FakeDuration: d}, getSkipTestResult())

	expected := `S

Failures/Skipped:

Title: failure
resource type: my resource id: a property: skipped

Total Duration: 2.000s
Count: 1, Failed: 0, Skipped: 1
`
	assert.Equal(t, expected, result)
	assert.Equal(t, 0, exitCode)
}
