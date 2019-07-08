package system

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"

	util2 "github.com/SimonBaeumer/goss/util"
	"github.com/aelsabbahy/GOnetstat"
	// This needs a better name
	"github.com/aelsabbahy/go-ps"
)

type Resource interface {
	Exists() (bool, error)
}

// System holds all constructor functions for each
type System struct {
	NewPackage     func(string, string) Package
	NewFile        func(string, *System, util2.Config) File
	NewAddr        func(string, *System, util2.Config) Addr
	NewPort        func(string, *System, util2.Config) Port
	NewService     func(string, *System, util2.Config) Service
	NewUser        func(string, *System, util2.Config) User
	NewGroup       func(string, *System, util2.Config) Group
	NewCommand     func(string, *System, util2.Config) Command
	NewDNS         func(string, *System, util2.Config) DNS
	NewProcess     func(string, *System, util2.Config) Process
	NewGossfile    func(string, *System, util2.Config) Gossfile
	NewKernelParam func(string, *System, util2.Config) KernelParam
	NewMount       func(string, *System, util2.Config) Mount
	NewInterface   func(string, *System, util2.Config) Interface
	NewHTTP        func(string, *System, util2.Config) HTTP
	ports          map[string][]GOnetstat.Process
	portsOnce      sync.Once
	procMap        map[string][]ps.Process
	procOnce       sync.Once
}

func (s *System) Ports() map[string][]GOnetstat.Process {
	s.portsOnce.Do(func() {
		s.ports = GetPorts(false)
	})
	return s.ports
}

func (s *System) ProcMap() map[string][]ps.Process {
	s.procOnce.Do(func() {
		s.procMap = GetProcs()
	})
	return s.procMap
}

//New creates the system object which holds all constructors for the system packages
func New() *System {
	sys := &System{
		NewPackage:     NewPackage,
		NewFile:        NewDefFile,
		NewAddr:        NewDefAddr,
		NewPort:        NewDefPort,
		NewUser:        NewDefUser,
		NewGroup:       NewDefGroup,
		NewCommand:     NewDefCommand,
		NewDNS:         NewDefDNS,
		NewProcess:     NewDefProcess,
		NewGossfile:    NewDefGossfile,
		NewKernelParam: NewDefKernelParam,
		NewMount:       NewDefMount,
		NewInterface:   NewDefInterface,
		NewHTTP:        NewDefHTTP,
	}
	sys.detectService()
	return sys
}

// DetectService adds the correct service creation function to a System struct
func (sys *System) detectService() {
	switch DetectService() {
	case "upstart":
		sys.NewService = NewServiceUpstart
	case "systemd":
		sys.NewService = NewServiceSystemd
	case "alpineinit":
		sys.NewService = NewAlpineServiceInit
	default:
		sys.NewService = NewServiceInit
	}
}

// DetectService attempts to detect what kind of service management the system
// is using, "systemd", "upstart", "alpineinit", or "init". It looks for systemctl
// command to detect systemd, and falls back on DetectDistro otherwise. If it can't
// decide, it returns "init".
func DetectService() string {
	if HasCommand("systemctl") {
		return "systemd"
	}
	// Centos Docker container doesn't run systemd, so we detect it or use init.
	switch DetectDistro() {
	case "ubuntu":
		return "upstart"
	case "alpine":
		return "alpineinit"
	case "arch":
		return "systemd"
	}
	return "init"
}

// DetectDistro attempts to detect which Linux distribution this computer is
// using. One of "ubuntu", "redhat" (including Centos), "alpine", "arch", or
// "debian". If it can't decide, it returns an empty string.
func DetectDistro() string {
	if b, e := ioutil.ReadFile("/etc/lsb-release"); e == nil && bytes.Contains(b, []byte("Ubuntu")) {
		return "ubuntu"
	} else if isRedhat() {
		return "redhat"
	} else if _, err := os.Stat("/etc/alpine-release"); err == nil {
		return "alpine"
	} else if _, err := os.Stat("/etc/arch-release"); err == nil {
		return "arch"
	} else if _, err := os.Stat("/etc/debian_version"); err == nil {
		return "debian"
	}
	return ""
}

// HasCommand returns whether or not an executable by this name is on the PATH.
func HasCommand(cmd string) bool {
	if _, err := exec.LookPath(cmd); err == nil {
		return true
	}
	return false
}

func isRedhat() bool {
	if _, err := os.Stat("/etc/redhat-release"); err == nil {
		return true
	} else if _, err := os.Stat("/etc/system-release"); err == nil {
		return true
	}
	return false
}
