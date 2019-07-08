package goss

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// GossRunTime represents the global runtime configs which can be set in goss
type GossRunTime struct {
	//Gossfile which should holds the test config
	Gossfile string
	//Vars file which holds the variabesl
	Vars string
	//Package defines which package manager you want to use, i.e. yum, apt, ...
	Package string //this does not belong here imho
	//Debug on true will create a more verbose output
	Debug bool
}

// Serve serves a new health endpoint
func (g *GossRunTime) Serve(endpoint string, handler *HealthHandler) {
	handler.Serve(endpoint)
}

// Validate starts the validation process
func (g *GossRunTime) Validate(v *Validator) int {
	return v.Validate(time.Now())
}

// Render renders a template file
func (g *GossRunTime) Render() (string, error) {
	goss, err := os.Open(g.Gossfile)
	if err != nil {
		return "", fmt.Errorf("Could not open gossfile with error: %s", err.Error())
	}
	defer goss.Close()

	vars, err := os.Open(g.Vars)
	if err != nil {
		return "", fmt.Errorf("Could not open varsfile with error: %s", err.Error())
	}
	defer vars.Close()

	return RenderJSON(goss, vars), nil
}

// GetGossConfig returns the goss configuration
func (g *GossRunTime) GetGossConfig() GossConfig {
	// handle stdin
	var fh *os.File
	var path, source string
	var gossConfig GossConfig
	TemplateFilter = NewTemplateFilter(g.Vars)
	specFile := g.Gossfile
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
	} else if specFile == "testing" {
		json := []byte(`
command:
  echo hello:
    exit-status: 0
    stdout: 
      - hello
    timeout: 10000
`)
		gossConfig = ReadJSONData(json, true)
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
