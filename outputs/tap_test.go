package outputs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTapDocumentation_Name(t *testing.T) {
	j := Tap{}
	assert.Equal(t, "tap", j.Name())
}

func TestTap_Output_Success(t *testing.T) {
	result, exitCode := runOutput(
		Tap{},
		getSuccessTestResult(),
		getSuccessTestResult(),
	)

	expected := `1..2
ok 1 - resource type: my resource id: a property: matches expectation: [expected]
ok 2 - resource type: my resource id: a property: matches expectation: [expected]
`
	assert.Equal(t, expected, result)
	assert.Equal(t, 0, exitCode)
}

func TestTap_Output_Fail(t *testing.T) {
	result, exitCode := runOutput(
		Tap{},
		getFailTestResult(),
	)

	expected := `1..1
not ok 1 - resource type: my resource id: a property: doesn't match, expect: [expected] found: []
`
	assert.Equal(t, expected, result)
	assert.Equal(t, 1, exitCode)
}

func TestTap_Output_Skip(t *testing.T) {
	result, exitCode := runOutput(
		Tap{},
		getSkipTestResult(),
	)

	expected := `1..1
ok 1 - # SKIP resource type: my resource id: a property: skipped
`
	assert.Equal(t, expected, result)
	assert.Equal(t, 0, exitCode)
}
