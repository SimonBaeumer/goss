package outputs

import (
	"io"
	"time"

	"github.com/SimonBaeumer/goss/resource"
	"github.com/SimonBaeumer/goss/util"
)

// Silent represents the silent output type
type Silent struct{}

// Name returns the name
func (r Silent) Name() string { return "silent" }

// Output writes the actual output
func (r Silent) Output(w io.Writer, results <-chan []resource.TestResult,
	startTime time.Time, outConfig util.OutputConfig) (exitCode int) {

	var failed int
	for resultGroup := range results {
		for _, testResult := range resultGroup {
			switch testResult.Result {
			case resource.FAIL:
				failed++
			}
		}
	}

	if failed > 0 {
		return 1
	}
	return 0
}

func init() {
	RegisterOutputer("silent", &Silent{}, []string{})
}
