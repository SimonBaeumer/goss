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

func TestDocumentation_Output(t *testing.T) {
	duration, _ := time.ParseDuration("2s")
	d := Documentation{FakeDuration: duration}
	result, exitCode := runOutput(d, GetExampleTestResult())

	expected := `Title: my title
resource type: my resource id: a property: matches expectation: [expected]


Total Duration: 2.000s
Count: 1, Failed: 0, Skipped: 0
`
	assert.Equal(t, expected, result)
	assert.Equal(t, 0, exitCode)
}
