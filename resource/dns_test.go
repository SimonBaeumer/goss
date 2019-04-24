package resource

import (
    "github.com/SimonBaeumer/goss/util"
    "github.com/stretchr/testify/assert"
    "testing"
)

var (
    QType = "TXT"
    conf = util.Config{Timeout: 50}
)


func TestNewDNS(t *testing.T) {
    mockDns := MockSysDNS{}
    dns, err := NewDNS(mockDns, conf)

    assert.Nil(t, err)
    assert.Implements(t, new(Resource), dns)
    assert.Equal(t, "TXT:google.com", dns.Host)
}

func TestNewDNS_WithoutQType(t *testing.T) {
    QType = ""

    mockDns := MockSysDNS{}
    dns, err := NewDNS(mockDns, conf)

    assert.Nil(t, err)
    assert.Implements(t, new(Resource), dns)
    assert.Equal(t, "google.com", dns.Host)
    assert.Equal(t, 50, dns.Timeout)
    assert.True(t, dns.Resolvable.(bool))
}

//MockSysDNS mocks the DNS system interface
type MockSysDNS struct {
}

func (dns MockSysDNS) Addrs() ([]string, error) {
    return []string{}, nil
}

func (dns MockSysDNS) Resolvable() (bool, error) {
    return true, nil
}

func (dns MockSysDNS) Exists() (bool, error) {
    panic("implement me")
}

func (dns MockSysDNS) Server() string {
    return "8.8.8.8"
}

func (dns MockSysDNS) Qtype() string {
    return QType
}

func (dns MockSysDNS) Host() string {
    return "google.com"
}