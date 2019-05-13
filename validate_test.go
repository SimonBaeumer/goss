package goss

import (
    "bytes"
    "github.com/SimonBaeumer/goss/outputs"
    "github.com/SimonBaeumer/goss/resource"
    "github.com/stretchr/testify/assert"
    "testing"
    "time"
)

func TestValidator_Validate(t *testing.T) {
    cmdResource := &resource.Command{Title: "echo hello", Command: "echo hello", ExitStatus: 0}

    w := &bytes.Buffer{}
    v := Validator{
        GossConfig: GossConfig{
            Commands: resource.CommandMap{"echo hello": cmdResource},
        },
        MaxConcurrent: 1,
        Outputer:      outputs.GetOutputer("documentation"),
        OutputWriter:  w,
    }

    r := v.Validate(time.Now())

    assert.Equal(t, 0, r)
    assert.Contains(t, w.String(), "Title: echo hello")
    assert.Contains(t, w.String(), "Command: echo hello: exit-status: matches expectation: [0]")
    assert.Contains(t, w.String(), "Count: 1, Failed: 0, Skipped: 0")
}