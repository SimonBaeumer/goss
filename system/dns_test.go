package system

import (
	"github.com/SimonBaeumer/goss/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDefDNS(t *testing.T) {
	dns := NewDefDNS("google.com", &System{}, util.Config{Timeout: 50, Server: "8.8.8.8"})

	assert.Implements(t, new(DNS), dns)
	assert.Equal(t, "8.8.8.8", dns.Server())
}

func TestNewDefDNS_WithQueryType(t *testing.T) {
	dns := NewDefDNS("TXT:google.com", &System{}, util.Config{})

	assert.Implements(t, new(DNS), dns)
	assert.Equal(t, "TXT", dns.Qtype())
	assert.Equal(t, "google.com", dns.Host())
}

func TestAddr(t *testing.T) {
	dns := NewDefDNS("localhost", &System{}, util.Config{Timeout: 200})

	r, err := dns.Resolvable()
	assert.Nil(t, err)
	assert.True(t, r)
}

func TestAddr_WithServer(t *testing.T) {
	dns := NewDefDNS("google.com", &System{}, util.Config{Timeout: 150, Server: "8.8.8.8"})

	r, err := dns.Resolvable()
	assert.Nil(t, err)
	assert.True(t, r)
}

func TestDNSLookup(t *testing.T) {
	dns := NewDefDNS("localhost", &System{}, util.Config{Timeout: 500})

	r, err := dns.Resolvable()
	addr, _ := dns.Addrs()

	assert.Nil(t, err)
	assert.True(t, r)
	assert.Equal(t, []string{"127.0.0.1"}, addr)
}

func TestDefDNS_Exists(t *testing.T) {
	dns := NewDefDNS("", &System{}, util.Config{})
	r, _ := dns.Exists()
	assert.False(t, r)
}

func TestDNSLookup_ShouldFail(t *testing.T) {
	dns := NewDefDNS("SRV:s-baeumer.de", &System{}, util.Config{Timeout: 500})

	r, err := dns.Resolvable()
	assert.Nil(t, err)
	assert.False(t, r)
}

func TestDNSLook_WithInvalidAddress(t *testing.T) {
	dns := NewDefDNS("thisdomaindoesnotexist123.com", &System{}, util.Config{Timeout: 50})
	r, _ := dns.Resolvable()
	assert.False(t, r)
}

func TestDNSLookupWithTimeout(t *testing.T) {
	dns := NewDefDNS("SRV:s-baeumer.de", &System{}, util.Config{Timeout: 0})
	r, err := dns.Resolvable()
	assert.Equal(t, "DNS lookup timed out (0s)", err.Error())
	assert.False(t, r)
}
