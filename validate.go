package goss

import (
	"fmt"
    "io"
    "os"
	"runtime"
	"sync"
	"time"

	"github.com/SimonBaeumer/goss/outputs"
	"github.com/SimonBaeumer/goss/resource"
	"github.com/SimonBaeumer/goss/system"
	"github.com/SimonBaeumer/goss/util"
	"github.com/fatih/color"
)

type Validator struct {
	GossConfig    GossConfig
	RetryTimeout  time.Duration
	Sleep         time.Duration
	FormatOptions []string
	Outputer      outputs.Outputer
	Package       string //Should be in the package resource config
	MaxConcurrent int    //Separating concurrency and validation, irritating atm...
	OutputWriter  io.Writer
}

// Validate validation runtime
func (v *Validator) Validate(startTime time.Time) int {
    if v.OutputWriter == nil {
        v.OutputWriter = os.Stdout
    }

	outputConfig := util.OutputConfig{
		FormatOptions: v.FormatOptions,
	}

	sys := system.New(v.Package)

	i := 1
	for {
		iStartTime := time.Now()

		out := validate(sys, v.GossConfig, v.MaxConcurrent)
		exitCode := v.Outputer.Output(v.OutputWriter, out, iStartTime, outputConfig)
		if v.RetryTimeout == 0 || exitCode == 0 {
			return exitCode
		}

		elapsed := time.Since(startTime)
		if elapsed + v.Sleep > v.RetryTimeout {
			color.Red("\nERROR: Timeout of %s reached before tests entered a passing state", v.RetryTimeout)
			return exitCode
		}
		color.Red("Retrying in %s (elapsed/timeout time: %.3fs/%s)\n\n\n", v.Sleep, elapsed.Seconds(), v.RetryTimeout)

		// Reset Cache
		sys = system.New(v.Package)
		time.Sleep(v.Sleep)
		i++
		fmt.Printf("Attempt #%d:\n", i)
	}
}

func validate(sys *system.System, gossConfig GossConfig, maxConcurrent int) <-chan []resource.TestResult {
    out := make(chan []resource.TestResult)
	in := make(chan resource.Resource)

	// Send resources to input channel
	go func() {
		for _, res := range gossConfig.Resources() {
			in <- res
		}
		close(in)
	}()

	// Read resources from input channel and validate
	workerCount := runtime.NumCPU() * 5
	if workerCount > maxConcurrent {
		workerCount = maxConcurrent
	}
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for res := range in {
				out <- res.Validate(sys)
			}
		}()
	}

	// Wait for the out channel to be finished, after that close it
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
