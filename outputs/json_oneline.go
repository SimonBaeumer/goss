package outputs

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/SimonBaeumer/goss/resource"
	"github.com/SimonBaeumer/goss/util"
	"github.com/fatih/color"
)

// JsonOneline represents the JsonOneline output type
type JsonOneline struct {
	duration time.Duration
}

// Name returns the name
func (r JsonOneline) Name() string { return "json_oneline" }

// Output writes the actual output
func (r JsonOneline) Output(w io.Writer, results <-chan []resource.TestResult,
	startTime time.Time, outConfig util.OutputConfig) (exitCode int) {

	color.NoColor = true
	testCount := 0
	failed := 0
	var resultsOut []map[string]interface{}
	for resultGroup := range results {
		for _, testResult := range resultGroup {
			if !testResult.Successful {
				failed++
			}
			m := struct2map(testResult)
			m["summary-line"] = humanizeResult(testResult)
			m["duration"] = int64(m["duration"].(float64))
			resultsOut = append(resultsOut, m)
			testCount++
		}
	}

	summary := make(map[string]interface{})
	duration := time.Since(startTime)
	//testing purposes
	if r.duration != 0 {
		duration = r.duration
	}
	summary["test-count"] = testCount
	summary["failed-count"] = failed
	summary["total-duration"] = duration
	summary["summary-line"] = fmt.Sprintf("Count: %d, Failed: %d, Duration: %.3fs", testCount, failed, duration.Seconds())

	out := make(map[string]interface{})
	out["results"] = resultsOut
	out["summary"] = summary

	j, _ := json.Marshal(out)
	fmt.Fprintln(w, string(j))

	if failed > 0 {
		return 1
	}

	return 0
}

func init() {
	RegisterOutputer("json_oneline", &JsonOneline{}, []string{})
}
