package system

import (
	"os/exec"
)

//-------------------------------------------------------------------------------------------------

type Runtime interface {
	Open(url string)
}

type ExecuteCommand func(cmd string, args ...string) *exec.Cmd

func DefaultRuntime() *SystemRuntime {
	return &SystemRuntime{
		OperatingSystem: Current(),
		ExecuteCommand:  exec.Command,
	}
}

type SystemRuntime struct {
	OperatingSystem OperatingSystem
	ExecuteCommand  ExecuteCommand
}

func (r *SystemRuntime) Open(url string) {
	cmd, args := OpenCommand(r.OperatingSystem, url)
	r.ExecuteCommand(cmd, args...).Start()
}

//-------------------------------------------------------------------------------------------------
