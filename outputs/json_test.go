package outputs

import (
	"bytes"
	"github.com/SimonBaeumer/goss/resource"
	"github.com/SimonBaeumer/goss/util"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestJson_Name(t *testing.T) {
	j := Json{}
	assert.Equal(t, "json", j.Name())
}

func TestJson_Output(t *testing.T) {
	var wg sync.WaitGroup
	b := &bytes.Buffer{}
	j := Json{FakeDuration: 1000}
	out := make(chan []resource.TestResult)
	r := 1

	go func() {
		defer wg.Done()
		wg.Add(1)
		r = j.Output(b, out, time.Now(), util.OutputConfig{})
	}()

	out <- GetExampleTestResult()

	close(out)
	wg.Wait()
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
	assert.Equal(t, expectedJson, b.String())
	assert.Equal(t, 0, r)
}
