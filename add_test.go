package goss

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ExtractHeaderArgument(t *testing.T) {
	expected := make(map[string][]string)
	expected["Set-Cookie"] = []string{"something"}

	headerArgs := "Set-Cookie: something"
	got := extractHeaderArgument(headerArgs)

	assert.Equal(t, expected, got)
}

func Test_ExtractHeaderArgumentWithoutEmptyArg(t *testing.T) {
	got := extractHeaderArgument("")
	assert.Equal(t, make(map[string][]string), got)
}
