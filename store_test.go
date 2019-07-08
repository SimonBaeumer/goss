package goss

import (
    "github.com/stretchr/testify/assert"
    "io/ioutil"
    "os"
    "testing"
)

const GossTestingEnvOS = "GOSS_TESTING_OS"

func Test_RenderJSON(t *testing.T) {
    err := os.Setenv(GossTestingEnvOS, "centos")
    if err != nil {
        panic(err.Error())
    }
    defer os.Unsetenv(GossTestingEnvOS)

    tmpVars, err := ioutil.TempFile("", "example_tmp_vars_*.yaml")
    if err != nil {
        panic(err.Error())
    }
    defer os.Remove(tmpVars.Name())

    _, err = tmpVars.WriteString(getExampleVars())
    if err != nil {
        panic(err.Error())
    }

    tmpGossfile, err := ioutil.TempFile("", "example_tmp_gossfile_*.yaml")
    if err != nil {
        panic(err.Error())
    }
    defer os.Remove(tmpGossfile.Name())

    _, err = tmpGossfile.WriteString(getExampleTemplate())
    if err != nil {
        panic(err.Error())
    }

    result := RenderJSON(tmpGossfile, tmpVars)

    assert.Equal(t, getExpecetd(), result)
}

func getExampleVars() string {
    return `
    centos:
      packages:
        kernel:
          - "4.9.11-centos"
          - "4.9.11-centos2"
    debian:
      packages:
        kernel:
          - "4.9.11-debian"
          - "4.9.11-debian2"
    users:
      - user1
      - user2
    `
}

func getExampleTemplate() string {
    return `
    package:
    # Looping over a variables defined in a vars.yaml using $OS environment variable as a lookup key
    {{range $name, $vers := index .Vars .Env.GOSS_TESTING_OS "packages"}}
      {{$name}}:
        installed: true
        versions:
        {{range $vers}}
          - {{.}}
        {{end}}
    {{end}}
    
    # This test is only when the OS environment variable matches the pattern
    {{if .Env.GOSS_TESTING_OS | regexMatch "[Cc]ent(OS|os)"}}
      libselinux:
        installed: true
    {{end}}
    
    # Loop over users
    user:
    {{range .Vars.users}}
      {{.}}:
        exists: true
        groups:
        - {{.}}
        home: /home/{{.}}
        shell: /bin/bash
    {{end}}
    
    
    package:
    {{if eq .Env.GOSS_TESTING_OS "centos"}}
      # This test is only when $OS environment variable is set to "centos"
      libselinux:
        installed: true
    {{end}}
    `
}

func getExpecetd() string {
    expected := `package:
  libselinux:
    installed: true
user:
  user1:
    exists: true
    groups:
    - user1
    home: /home/user1
    shell: /bin/bash
  user2:
    exists: true
    groups:
    - user2
    home: /home/user2
    shell: /bin/bash
`
    return expected
}
