package resource

import (
	"fmt"
	"github.com/SimonBaeumer/goss/system"
	"path/filepath"
	"strconv"
)

// A Resource defines a type on which tests can be executed, i.e. http or file
type Resource interface {
	Validate(*system.System) []TestResult
	SetID(string)
}

type ResourceRead interface {
	ID() string
	GetTitle() string
	GetMeta() meta
}

type matcher interface{}
type meta map[string]interface{}

func contains(a []string, s string) bool {
	for _, e := range a {
		if m, _ := filepath.Match(e, s); m {
			return true
		}
	}
	return false
}

func deprecateAtoI(depr interface{}, desc string) interface{} {
	s, ok := depr.(string)
	if !ok {
		return depr
	}
	fmt.Printf("DEPRECATION WARNING: %s should be an integer not a string\n", desc)
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return float64(i)
}

func shouldSkip(results []TestResult) bool {
	if results[0].Err != nil || results[0].Found[0] == "false" {
		return true
	}
	return false
}
