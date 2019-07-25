package outputs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJson_Name(t *testing.T) {
	j := Json{}
	assert.Equal(t, "json", j.Name())
}

func TestJson_Output(t *testing.T) {
	result, exitCode := runOutput(
		Json{FakeDuration: 1000},
		getSuccessTestResult(),
	)

	expectedJson := `{
    "results": [
        {
            "duration": 500,
            "err": null,
            "expected": [
                "expected"
            ],
            "found": null,
            "human": "",
            "meta": null,
            "property": "a property",
            "resource-id": "my resource id",
            "resource-type": "resource type",
            "result": 0,
            "successful": true,
            "summary-line": "resource type: my resource id: a property: matches expectation: [expected]",
            "test-type": 0,
            "title": "my title"
        }
    ],
    "summary": {
        "failed-count": 0,
        "summary-line": "Count: 1, Failed: 0, Duration: 0.000s",
        "test-count": 1,
        "total-duration": 1000
    }
}
`
	assert.Equal(t, expectedJson, result)
	assert.Equal(t, 0, exitCode)
}

func TestJson_Output_FAIL(t *testing.T) {
	result, exitCode := runOutput(
		Json{FakeDuration: 1000},
		getFailTestResult(),
	)

	expected := `{
    "results": [
        {
            "duration": 500,
            "err": null,
            "expected": [
                "expected"
            ],
            "found": null,
            "human": "",
            "meta": null,
            "property": "a property",
            "resource-id": "my resource id",
            "resource-type": "resource type",
            "result": 1,
            "successful": false,
            "summary-line": "resource type: my resource id: a property: doesn't match, expect: [expected] found: []",
            "test-type": 0,
            "title": "failure"
        }
    ],
    "summary": {
        "failed-count": 1,
        "summary-line": "Count: 1, Failed: 1, Duration: 0.000s",
        "test-count": 1,
        "total-duration": 1000
    }
}
`
	assert.Equal(t, expected, result)
	assert.Equal(t, 1, exitCode)
}
