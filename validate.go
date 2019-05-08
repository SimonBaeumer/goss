package goss

import (
	"fmt"
	"github.com/SimonBaeumer/goss/internal/app"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/SimonBaeumer/goss/outputs"
	"github.com/SimonBaeumer/goss/resource"
	"github.com/SimonBaeumer/goss/system"
	"github.com/SimonBaeumer/goss/util"
	"github.com/fatih/color"
)

func getGossConfig(ctx app.CliContext) GossConfig {
	// handle stdin
	var fh *os.File
	var path, source string
	var gossConfig GossConfig
	TemplateFilter = NewTemplateFilter(ctx.Vars)
	specFile := ctx.Gossfile
	if specFile == "-" {
		source = "STDIN"
		fh = os.Stdin
		data, err := ioutil.ReadAll(fh)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		OutStoreFormat = getStoreFormatFromData(data)
		gossConfig = ReadJSONData(data, true)
	} else {
		source = specFile
		path = filepath.Dir(specFile)
		OutStoreFormat = getStoreFormatFromFileName(specFile)
		gossConfig = ReadJSON(specFile)
	}

	gossConfig = mergeJSONData(gossConfig, 0, path)

	if len(gossConfig.Resources()) == 0 {
		fmt.Printf("Error: found 0 tests, source: %v\n", source)
		os.Exit(1)
	}
	return gossConfig
}

func getOutputer(ctx app.CliContext) outputs.Outputer {
	if ctx.NoColor {
		color.NoColor = true
	}
	if ctx.Color {
		color.NoColor = false
	}
	return outputs.GetOutputer(ctx.Format)
}

// Validate validation runtime
func Validate(ctx app.CliContext, startTime time.Time) {

	outputConfig := util.OutputConfig{
		FormatOptions: ctx.FormatOptions,
	}

	gossConfig := getGossConfig(ctx)
	sys := system.New(ctx.Package)
	outputer := getOutputer(ctx)

	sleep := ctx.Sleep
	retryTimeout := ctx.RetryTimeout
	i := 1
	for {
		iStartTime := time.Now()
		out := validate(sys, gossConfig, ctx.MaxConcurrent)
		exitCode := outputer.Output(os.Stdout, out, iStartTime, outputConfig)
		if retryTimeout == 0 || exitCode == 0 {
			os.Exit(exitCode)
		}
		elapsed := time.Since(startTime)
		if elapsed+sleep > retryTimeout {
			color.Red("\nERROR: Timeout of %s reached before tests entered a passing state", retryTimeout)
			os.Exit(3)
		}
		color.Red("Retrying in %s (elapsed/timeout time: %.3fs/%s)\n\n\n", sleep, elapsed.Seconds(), retryTimeout)
		// Reset cache
		sys = system.New(ctx.Package)
		time.Sleep(sleep)
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
