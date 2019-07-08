package outputs

import (
	"fmt"
	"io"
	"time"

	"github.com/SimonBaeumer/goss/resource"
	"github.com/SimonBaeumer/goss/util"
)

// Rspecish represents the rspecish output type
type Rspecish struct {
	//FakeDuration will only be needed for testing purposes
	FakeDuration time.Duration
}

// Name returns the name
func (r Rspecish) Name() string { return "rspecish" }

// Output writes the actual output
func (r Rspecish) Output(w io.Writer, results <-chan []resource.TestResult,
	startTime time.Time, outConfig util.OutputConfig) (exitCode int) {

	testCount := 0
	var failedOrSkipped [][]resource.TestResult
	var skipped, failed int
	for resultGroup := range results {
		failedOrSkippedGroup := []resource.TestResult{}
		for _, testResult := range resultGroup {
			switch testResult.Result {
			case resource.SUCCESS:
				fmt.Fprintf(w, green("."))
			case resource.SKIP:
				fmt.Fprintf(w, yellow("S"))
				failedOrSkippedGroup = append(failedOrSkippedGroup, testResult)
				skipped++
			case resource.FAIL:
				fmt.Fprintf(w, red("F"))
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
	RegisterOutputer("rspecish", &Rspecish{}, []string{})
}
