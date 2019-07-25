package outputs

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestJunit_Name(t *testing.T) {
	j := JUnit{}
	assert.Equal(t, "junit", j.Name())
}

func getDate() string {
	return time.Date(
		2019, 01, 01, 10, 30, 30, 3000, time.UTC).Format(time.RFC3339)
}

func TestJunit_Output_Fail(t *testing.T) {
	result, exitCode := runOutput(JUnit{testingTimestamp: getDate()}, getFailTestResult())

	expected := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="goss" errors="0" tests="1" failures="1" skipped="0" time="0.000" testingTimestamp="2019-01-01T10:30:30Z">
<testcase name="resource type my resource id a property" time="0.000">
<system-err>resource type: my resource id: a property: doesn&#39;t match, expect: [expected] found: []</system-err>
<failure>resource type: my resource id: a property: doesn&#39;t match, expect: [expected] found: []</failure>
</testcase>
</testsuite>
`
	assert.Equal(t, expected, result)
	assert.Equal(t, 1, exitCode)
}

func TestJunit_Output_Success(t *testing.T) {
	result, exitCode := runOutput(JUnit{testingTimestamp: getDate()}, getSuccessTestResult())

	expected := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="goss" errors="0" tests="1" failures="0" skipped="0" time="0.000" testingTimestamp="2019-01-01T10:30:30Z">
<testcase name="resource type my resource id a property" time="0.000">
<system-out>resource type: my resource id: a property: matches expectation: [expected]</system-out>
</testcase>
</testsuite>
`
	assert.Equal(t, expected, result)
	assert.Equal(t, 0, exitCode)
}

func TestJunit_Output_Skip(t *testing.T) {
	result, exitCode := runOutput(JUnit{testingTimestamp: getDate()}, getSkipTestResult())

	expected := `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="goss" errors="0" tests="1" failures="0" skipped="1" time="0.000" testingTimestamp="2019-01-01T10:30:30Z">
<testcase name="resource type my resource id a property" time="0.000">
<skipped/><system-out>resource type: my resource id: a property: skipped</system-out>
</testcase>
</testsuite>
`
	assert.Equal(t, expected, result)
	assert.Equal(t, 0, exitCode)
}
