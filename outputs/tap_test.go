package outputs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTapDocumentation_Name(t *testing.T) {
	j := Tap{}
	assert.Equal(t, "tap", j.Name())
}

func TestTapDocumentation_Output(t *testing.T) {
	result, exitCode := runOutput(
		Tap{},
		GetExampleTestResult(),
		GetExampleTestResult(),
	)

	expected := `1..2
ok 1 - resource type: my resource id: a property: matches expectation: [expected]
ok 2 - resource type: my resource id: a property: matches expectation: [expected]
`
	assert.Equal(t, expected, result)
	assert.Equal(t, 0, exitCode)
}
