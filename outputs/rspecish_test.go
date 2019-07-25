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

func TestRspecish_Output(t *testing.T) {
	d, _ := time.ParseDuration("2s")

	result, exitCode := runOutput(Rspecish{FakeDuration: d}, GetExampleTestResult())

	expectedJson := `.

Total Duration: 2.000s
Count: 1, Failed: 0, Skipped: 0
`
	assert.Equal(t, expectedJson, result)
	assert.Equal(t, 0, exitCode)
}
