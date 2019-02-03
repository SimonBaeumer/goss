package system

import (
	"github.com/SimonBaeumer/goss/util"
	"github.com/opencontainers/runc/libcontainer/user"
)

type Group interface {
	Groupname() string
	Exists() (bool, error)
	GID() (int, error)
}

type DefGroup struct {
	groupname string
}

func NewDefGroup(groupname string, system *System, config util.Config) Group {
	return &DefGroup{groupname: groupname}
}

// Groupname returns the groupname
func (u *DefGroup) Groupname() string {
	return u.groupname
}

// Exists checks if the group exists on the host system
func (u *DefGroup) Exists() (bool, error) {
	_, err := user.LookupGroup(u.groupname)
	if err != nil {
		return false, nil
	}
	return true, nil
}

// GID returns the GID
func (u *DefGroup) GID() (int, error) {
	group, err := user.LookupGroup(u.groupname)
	if err != nil {
		return 0, err
	}

	return group.Gid, nil
}
