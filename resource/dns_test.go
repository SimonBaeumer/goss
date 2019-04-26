package resource

import (
    "github.com/SimonBaeumer/goss/system"
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

func TestDNS_Validate(t *testing.T) {
    addrs := convertToInterfaceSlice()

    mockDns := MockSysDNS{}
    dns, _ := NewDNS(mockDns, conf)
    dns.Resolvable = true
    dns.Host = "localhost"
    dns.Addrs = addrs
    dns.Timeout = 0

    sys := system.System{}
    sys.NewDNS = func (host string, sys *system.System, config util.Config) system.DNS {
        return &MockSysDNS{Addr: []string{"localhost:53"}}
    }

    r := dns.Validate(&sys)

    assert.Equal(t, 500, dns.Timeout, "Could not set default timeout if 0 was given")
    assert.Len(t, r, 2)

    assert.True(t, r[0].Successful)
    assert.Equal(t, "resolvable", r[0].Property)

    assert.True(t, r[1].Successful)
    assert.Equal(t, "addrs", r[1].Property)
}

func TestDNS_ValidateFail(t *testing.T) {
    addrs := convertToInterfaceSlice()

    mockDns := MockSysDNS{}
    dns, _ := NewDNS(mockDns, conf)
    dns.Timeout = 50
    dns.Resolvable = true
    dns.Host = "localhost"
    dns.Addrs = addrs

    sys := system.System{}
    sys.NewDNS = func(host string, sys *system.System, config util.Config) system.DNS {
        return &MockSysDNS{Addr: []string{"ns.localhost"}}
    }

    r := dns.Validate(&sys)

    assert.Len(t, r, 2)
    assert.True(t, r[0].Successful)
    assert.Equal(t, "resolvable", r[0].Property)

    assert.False(t, r[1].Successful)
    assert.Equal(t, "addrs", r[1].Property)
    expectedHuman := `Expected
    <[]string | len:1, cap:1>: ["ns.localhost"]
to contain element matching
    <string>: localhost:53`
    assert.Equal(t, expectedHuman, r[1].Human)
}

func convertToInterfaceSlice() []interface{} {
    // Create expected addrs as interface{} slice
    // It is necessary to allocate the memory before, because []interface{} is of an unknown size
    var expect = []string{"localhost:53"}
    var addrs = make([]interface{}, len(expect))
    for i, char := range expect {
        addrs[i] = char
    }
    return addrs
}

//MockSysDNS mocks the DNS system interface
type MockSysDNS struct {
    Addr []string
}

func (dns MockSysDNS) Addrs() ([]string, error) {
    return dns.Addr, nil
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