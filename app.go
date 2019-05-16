package goss

import (
    "fmt"
    "github.com/SimonBaeumer/goss/internal/app"
    "io/ioutil"
    "os"
    "path/filepath"
)

// GossRunTime represents the global runtime configs which can be set in goss
type GossRunTime struct {
    //Gossfile which should holds the test config
    Gossfile    string
    //Vars file which holds the variabesl
    Vars       string
    //Package defines which package manager you want to use, i.e. yum, apt, ...
    Package    string //this does not belong here imho
}

func NewGossRunTime(ctx app.CliContext) *GossRunTime {
    return &GossRunTime{}
}

func (g *GossRunTime) Serve() {
    //Serve()
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