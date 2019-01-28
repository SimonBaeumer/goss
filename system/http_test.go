package system

import (
	"github.com/SimonBaeumer/goss/util"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

const SuccessStatusCode = 200

func newTestingDefHTTP() HTTP {
	return NewDefHTTP(
		"http://goss.rocks",
		&System{},
		util.Config{},
	)
}

func Test_NewDefHTTP(t *testing.T) {
	system := &System{}
	conf := util.Config{}

	got := NewDefHTTP(
		"http://goss.rocks",
		system,
		conf,
	)

	assert.Implements(t, (*HTTP)(nil), got, "DefHTTP does not implement HTTP")
	assert.Equal(t, got.HTTP(), "http://goss.rocks")
}

func TestDefHTTP_Body(t *testing.T) {
	http := newTestingDefHTTP()

	got, err := http.Body()

	assert.Nil(t, err)
	assert.Implements(t, new(io.Reader), got)
}

func TestDefHTTP_Status(t *testing.T) {
	http := newTestingDefHTTP()

	got, err := http.Status()

	assert.Nil(t, err)
	assert.Equal(t, SuccessStatusCode, got)
}

func TestDefHTTP_Headers(t *testing.T) {
	defHTTP := newTestingDefHTTP()

	got, err := defHTTP.Headers()

	assert.Nil(t, err)
	assert.IsType(t, make(Header), got)
}
