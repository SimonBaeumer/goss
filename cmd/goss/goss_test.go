package main

import (
    "bytes"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestApp_Validate(t *testing.T) {
    b := &bytes.Buffer{}
    app := createApp()
    app.Writer = b

    r := app.Run([]string{"", "--gossfile", "testing", "validate"})

    assert.Nil(t, r)
    assert.Contains(t, b.String(), "Count: 2, Failed: 0, Skipped: 0")
}
