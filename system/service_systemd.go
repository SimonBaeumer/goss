package system

import (
	"fmt"
	"strings"

	"github.com/SimonBaeumer/goss/util"
)

type ServiceSystemd struct {
	service string
}

func NewServiceSystemd(service string, system *System, config util.Config) Service {
	return &ServiceSystemd{
		service: service,
	}
}

func (s *ServiceSystemd) Service() string {
	return s.service
}

func (s *ServiceSystemd) Exists() (bool, error) {
	if invalidService(s.service) {
		return false, nil
	}

	cmd := util.NewCommand("systemctl", "-q", "list-unit-files", "--type=service")
	if err := cmd.Run(); err != nil {
		return false, err
	}

	if strings.Contains(cmd.Stdout.String(), fmt.Sprintf("%s.service", s.service)) {
		return true, cmd.Err
	}
	// Fallback on sysv
	sysv := &ServiceInit{service: s.service}
	if e, err := sysv.Exists(); e && err == nil {
		return true, nil
	}
	return false, nil
}

func (s *ServiceSystemd) Enabled() (bool, error) {
	if invalidService(s.service) {
		return false, nil
	}
	cmd := util.NewCommand("systemctl", "-q", "is-enabled", s.service)
	if err := cmd.Run(); err != nil {
		return false, err
	}

	if cmd.Status == 0 {
		return true, cmd.Err
	}
	// Fallback on sysv
	sysv := &ServiceInit{service: s.service}
	if en, err := sysv.Enabled(); en && err == nil {
		return true, nil
	}
	return false, nil
}

func (s *ServiceSystemd) Running() (bool, error) {
	if invalidService(s.service) {
		return false, nil
	}

	cmd := util.NewCommand("systemctl", "-q", "is-active", s.service)
	if err := cmd.Run(); err != nil {
		return false, err
	}

	if cmd.Status == 0 {
		return true, cmd.Err
	}
	// Fallback on sysv
	sysv := &ServiceInit{service: s.service}
	if r, err := sysv.Running(); r && err == nil {
		return true, nil
	}
	return false, nil
}
