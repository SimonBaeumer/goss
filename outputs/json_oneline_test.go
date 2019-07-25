package outputs

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestJsonOneline_Name(t *testing.T) {
	j := JsonOneline{}
	assert.Equal(t, "json_oneline", j.Name())
}

func TestJsonOneline_Output_FAIL(t *testing.T) {
	duration, _ := time.ParseDuration("2s")
	result, exitCode := runOutput(JsonOneline{duration: duration}, getFailTestResult())

	expected := `{"results":[{"duration":500,"err":null,"expected":["expected"],"found":null,"human":"","meta":null,"property":"a property","resource-id":"my resource id","resource-type":"resource type","result":1,"successful":false,"summary-line":"resource type: my resource id: a property: doesn't match, expect: [expected] found: []","test-type":0,"title":"failure"}],"summary":{"failed-count":1,"summary-line":"Count: 1, Failed: 1, Duration: 2.000s","test-count":1,"total-duration":2000000000}}
`
	assert.Equal(t, expected, result)
	assert.Equal(t, 1, exitCode)
}

func TestJsonOneline_Output_Success(t *testing.T) {
	duration, _ := time.ParseDuration("2s")
	result, exitCode := runOutput(JsonOneline{duration: duration}, getSuccessTestResult())

	expected := `{"results":[{"duration":500,"err":null,"expected":["expected"],"found":null,"human":"","meta":null,"property":"a property","resource-id":"my resource id","resource-type":"resource type","result":0,"successful":true,"summary-line":"resource type: my resource id: a property: matches expectation: [expected]","test-type":0,"title":"my title"}],"summary":{"failed-count":0,"summary-line":"Count: 1, Failed: 0, Duration: 2.000s","test-count":1,"total-duration":2000000000}}
`
	assert.Equal(t, expected, result)
	assert.Equal(t, 0, exitCode)
}

func TestJsonOneline_Output_Skip(t *testing.T) {
	duration, _ := time.ParseDuration("2s")
	result, exitCode := runOutput(JsonOneline{duration: duration}, getSkipTestResult())

	expected := `{"results":[{"duration":500,"err":null,"expected":["expected"],"found":null,"human":"","meta":null,"property":"a property","resource-id":"my resource id","resource-type":"resource type","result":2,"successful":true,"summary-line":"resource type: my resource id: a property: skipped","test-type":0,"title":"failure"}],"summary":{"failed-count":0,"summary-line":"Count: 1, Failed: 0, Duration: 2.000s","test-count":1,"total-duration":2000000000}}
`
	assert.Equal(t, expected, result)
	assert.Equal(t, 0, exitCode)
}
