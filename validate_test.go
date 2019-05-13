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
    cmdRes := &resource.Command{Title: "echo hello", Command: "echo hello", ExitStatus: 0}
    fileRes := &resource.File{Title: "/tmp", Path: "/tmp", Filetype: "directory", Exists: true}
    addrRes := &resource.Addr{Title: "tcp://google.com:443", Address: "tcp://google.com:443", Reachable: true}
    httpRes := &resource.HTTP{Title: "https://google.com", HTTP: "https://google.com", Status: 200}
    userRes := &resource.User{Title: "root", Username: "root", Exists: true}
    groupRes := &resource.Group{Title: "root", Groupname: "root", Exists: true}
    dnsRes := &resource.DNS{Title: "A:google.com", Host: "A:google.com", Resolvable: true}

    w := &bytes.Buffer{}
    v := Validator{
        GossConfig: GossConfig{
            Commands: resource.CommandMap{"echo hello": cmdRes},
            Files: resource.FileMap{"/tmp": fileRes},
            Addrs: resource.AddrMap{"127.0.0.1": addrRes},
            HTTPs: resource.HTTPMap{"https://google.com": httpRes},
            Users: resource.UserMap{"root": userRes},
            Groups: resource.GroupMap{"root": groupRes},
            DNS: resource.DNSMap{"A:https://google.com": dnsRes},
        },
        MaxConcurrent: 1,
        Outputer:      outputs.GetOutputer("documentation"),
        OutputWriter:  w,
    }

    r := v.Validate(time.Now())

    assert.Equal(t, 0, r)
    assert.Contains(t, w.String(), "Count: 8, Failed: 0, Skipped: 0")
}