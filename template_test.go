package goss

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func Test_NewTemplateFilter_Variable(t *testing.T) {
	vars, err := ioutil.TempFile("", "*_vars.yaml")
	if err != nil {
		panic(err.Error())
	}
	defer os.Remove(vars.Name())

	_, err = vars.WriteString("test: testing")
	if err != nil {
		panic(err.Error())
	}

	content := []byte(`variable: {{.Vars.test}}`)

	filter := NewTemplateFilter(vars.Name())
	result := filter(content)

	assert.Equal(t, "variable: testing", string(result))
}

func Test_NewTemplateFilter_Env(t *testing.T) {
	err := os.Setenv("GOSS_TEST_ENV", "env testing")
	if err != nil {
		panic(err.Error())
	}
	defer os.Unsetenv("template_goss_test_env")

	content := []byte(`environment: {{.Env.GOSS_TEST_ENV}}`)

	filter := NewTemplateFilter("")
	result := filter(content)

	assert.Equal(t, "environment: env testing", string(result))
}

func Test_NewTemplateFilter_mkSlice(t *testing.T) {
	content := []byte(`{{- range mkSlice "test1" "test2" "test3"}}{{.}}{{end}}`)

	filter := NewTemplateFilter("")
	result := filter(content)

	assert.Equal(t, "test1test2test3", string(result))
}

func Test_NewTemplateFilter_readFile(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "read_file_temp")
	if err != nil {
		panic(err.Error())
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString("test read file from template")

	content := []byte(`{{readFile "` + tmpFile.Name() + `"}}`)

	filter := NewTemplateFilter("")
	result := filter(content)

	assert.Equal(t, "test read file from template", string(result))
}

func Test_NewTemplateFilter_regexMatch(t *testing.T) {
	content := []byte(`{{if "centos" | regexMatch "[Cc]ent(OS|os)"}}detected regex{{end}}`)

	filter := NewTemplateFilter("")
	result := filter(content)

	assert.Equal(t, "detected regex", string(result))
}

func Test_NewTemplateFilter_regexMatch_fail(t *testing.T) {
	content := []byte(`{{if "ubuntu" | regexMatch "[Cc]ent(OS|os)"}}detected regex{{else}}no match{{end}}`)

	filter := NewTemplateFilter("")
	result := filter(content)

	assert.Equal(t, "no match", string(result))
}
