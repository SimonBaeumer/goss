package util

import (
	"bytes"
	//"fmt"
	"os/exec"
	"syscall"
)

// Command represents a command which can be executed
type Command struct {
	name           string
	Cmd            *exec.Cmd
	Stdout, Stderr bytes.Buffer
	Err            error
	Status         int
}

// NewCommand creates a command
func NewCommand(name string, arg ...string) *Command {
	//fmt.Println(arg)
	command := new(Command)
	command.name = name
	command.Cmd = exec.Command(name, arg...)
	return command
}

// Run executes the command and writes the results to its properties
func (c *Command) Run() error {
	c.Cmd.Stdout = &c.Stdout
	c.Cmd.Stderr = &c.Stderr

	if _, err := exec.LookPath(c.name); err != nil {
		c.Err = err
		return c.Err
	}

	if err := c.Cmd.Start(); err != nil {
		c.Err = err
		//log.Fatalf("Cmd.Start: %v")
	}

	if err := c.Cmd.Wait(); err != nil {
		c.Err = err
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				c.Status = status.ExitStatus()
				//log.Printf("Exit Status: %d", status.ExitStatus())
			}
		}
	} else {
		c.Status = 0
	}
	return c.Err
}
