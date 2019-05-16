package outputs

import (
	"fmt"
	"io"
	"time"

	"github.com/SimonBaeumer/goss/resource"
	"github.com/SimonBaeumer/goss/util"
)

// Documentation represents the documentation output type
type Documentation struct{
	//FakeDuration will only be used for testing purposes
	FakeDuration time.Duration
}

// Name returns the name
func (r Documentation) Name() string { return "documentation" }

func (r Documentation) Output(w io.Writer, results <-chan []resource.TestResult,
	startTime time.Time, outConfig util.OutputConfig) (exitCode int) {

	testCount := 0
	var failedOrSkipped [][]resource.TestResult
	var skipped, failed int
	for resultGroup := range results {
		failedOrSkippedGroup := []resource.TestResult{}
		first := resultGroup[0]
		header := header(first)
		if header != "" {
			fmt.Fprint(w, header)
		}
		for _, testResult := range resultGroup {
			switch testResult.Result {
			case resource.SUCCESS:
				fmt.Fprintln(w, humanizeResult(testResult))
			case resource.SKIP:
				fmt.Fprintln(w, humanizeResult(testResult))
				failedOrSkippedGroup = append(failedOrSkippedGroup, testResult)
				skipped++
			case resource.FAIL:
				fmt.Fprintln(w, humanizeResult(testResult))
				failedOrSkippedGroup = append(failedOrSkippedGroup, testResult)
				failed++
			}
			testCount++
		}
		if len(failedOrSkippedGroup) > 0 {
			failedOrSkipped = append(failedOrSkipped, failedOrSkippedGroup)
		}
	}

	fmt.Fprint(w, "\n\n")
	fmt.Fprint(w, failedOrSkippedSummary(failedOrSkipped))

	duration := time.Since(startTime)
	if r.FakeDuration != 0 {
		duration = r.FakeDuration
	}
	fmt.Fprint(w, summary(duration.Seconds(), testCount, failed, skipped))
	if failed > 0 {
		return 1
	}
	return 0
}

func init() {
	RegisterOutputer("documentation", &Documentation{}, []string{})
}
