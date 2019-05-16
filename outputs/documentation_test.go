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

func TestDocumentation_Name(t *testing.T) {
    j := Documentation{}
    assert.Equal(t, "documentation", j.Name())
}

func TestDocumentation_Output(t *testing.T) {
    var wg sync.WaitGroup
    b := &bytes.Buffer{}
    d, _ := time.ParseDuration("2s")
    j := Documentation{FakeDuration: d}
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
    expectedJson := `Title: my title
resource type: my resource id: a property: matches expectation: [expected]


Total Duration: 2.000s
Count: 1, Failed: 0, Skipped: 0
`
    assert.Equal(t, expectedJson, b.String())
    assert.Equal(t, 0, r)
}
