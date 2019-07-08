package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
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

func TestApp_Add(t *testing.T) {
	b := &bytes.Buffer{}
	app := createApp()
	app.Writer = b

	file, err := ioutil.TempFile("/tmp", "testing_goss_*.yaml")
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	r := app.Run([]string{"", "--gossfile", file.Name(), "add", "http", "http://google.com"})

	assert.Nil(t, r)
	assert.Contains(t, b.String(), getAddResult())
	assert.Contains(t, b.String(), "Adding HTTP to '/tmp/testing_goss_")
}

func getAddResult() string {
	return `http://google.com:
  status: 200
  allow-insecure: false
  no-follow-redirects: false
  timeout: 5000
  body: []`
}
