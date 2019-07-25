package outputs

import (
	"bytes"
	"github.com/SimonBaeumer/goss/resource"
	"github.com/SimonBaeumer/goss/util"
	"sync"
	"time"
)

// runOutput runs the output on the given outputer
func runOutput(outputer Outputer, results ...resource.TestResult) (string, int) {
	var wg sync.WaitGroup
	buffer := &bytes.Buffer{}
	//d, _ := time.ParseDuration("2s")
	out := make(chan []resource.TestResult)
	codeChan := make(chan int)

	go func(o Outputer, b *bytes.Buffer, e chan int) {
		defer wg.Done()
		wg.Add(1)
		exit := o.Output(b, out, time.Now(), util.OutputConfig{})
		codeChan <- exit
	}(outputer, buffer, codeChan)

	//Write results to channel
	out <- results
	close(out)

	exitCode := <-codeChan
	close(codeChan)
	wg.Wait()

	return buffer.String(), exitCode
}

//getSuccessTestResult returns an example test result
func getSuccessTestResult() resource.TestResult {
	return resource.TestResult{
		Title:        "my title",
		Duration:     time.Duration(500),
		Successful:   true,
		Result:       resource.SUCCESS,
		ResourceType: "resource type",
		ResourceId:   "my resource id",
		Property:     "a property",
		Expected:     []string{"expected"},
	}
}

func getFailTestResult() resource.TestResult {
	return resource.TestResult{
		Title:        "failure",
		Duration:     time.Duration(500),
		Successful:   false,
		Result:       resource.FAIL,
		ResourceType: "resource type",
		ResourceId:   "my resource id",
		Property:     "a property",
		Expected:     []string{"expected"},
	}
}

func getSkipTestResult() resource.TestResult {
	return resource.TestResult{
		Title:        "failure",
		Duration:     time.Duration(500),
		Successful:   true,
		Result:       resource.SKIP,
		ResourceType: "resource type",
		ResourceId:   "my resource id",
		Property:     "a property",
		Expected:     []string{"expected"},
	}
}
