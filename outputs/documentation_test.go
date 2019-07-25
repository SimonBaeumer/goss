package outputs

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDocumentation_Name(t *testing.T) {
	d := Documentation{}
	assert.Equal(t, "documentation", d.Name())
}

func TestDocumentation_Output_Success(t *testing.T) {
	duration, _ := time.ParseDuration("2s")
	d := Documentation{FakeDuration: duration}
	result, exitCode := runOutput(d, getSuccessTestResult())

	expected := `Title: my title
resource type: my resource id: a property: matches expectation: [expected]


Total Duration: 2.000s
Count: 1, Failed: 0, Skipped: 0
`
	assert.Equal(t, expected, result)
	assert.Equal(t, 0, exitCode)
}

func TestDocumentation_Output_Fail(t *testing.T) {
	duration, _ := time.ParseDuration("2s")
	d := Documentation{FakeDuration: duration}
	result, exitCode := runOutput(d, getFailTestResult())

	expected := `Title: failure
resource type: my resource id: a property: doesn't match, expect: [expected] found: []


Failures/Skipped:

Title: failure
resource type: my resource id: a property: doesn't match, expect: [expected] found: []

Total Duration: 2.000s
Count: 1, Failed: 1, Skipped: 0
`
	assert.Equal(t, expected, result)
	assert.Equal(t, 1, exitCode)
}

func TestDocumentation_Output_Skip(t *testing.T) {
	duration, _ := time.ParseDuration("2s")
	d := Documentation{FakeDuration: duration}
	result, exitCode := runOutput(d, getSkipTestResult())

	expected := `Title: failure
resource type: my resource id: a property: skipped


Failures/Skipped:

Title: failure
resource type: my resource id: a property: skipped

Total Duration: 2.000s
Count: 1, Failed: 0, Skipped: 1
`
	assert.Equal(t, expected, result)
	assert.Equal(t, 0, exitCode)
}
