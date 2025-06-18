package system_test

import (
	"os/exec"
	"testing"

	"github.com/vaguevoid/cloud-cli/internal/lib/system"
	"github.com/vaguevoid/cloud-cli/internal/test/assert"
)

const TestURL = "https://play.void.dev"

//-------------------------------------------------------------------------------------------------

func TestDefaultRuntime(t *testing.T) {
	runtime := system.DefaultRuntime()
	assert.NotNil(t, runtime)
	assert.Equal(t, system.Current(), runtime.OperatingSystem)
}

//-------------------------------------------------------------------------------------------------

func TestOpenBrowserWindows(t *testing.T) {
	runtime := system.DefaultRuntime()
	runtime.OperatingSystem = system.OperatingSystemWindows
	runtime.ExecuteCommand = func(cmd string, args ...string) *exec.Cmd {
		assert.Equal(t, "rundll32", cmd)
		assert.Equal(t, []string{"url.dll,FileProtocolHandler", TestURL}, args)
		return &exec.Cmd{}
	}
	runtime.Open(TestURL)
}

func TestOpenBrowserMac(t *testing.T) {
	runtime := system.DefaultRuntime()
	runtime.OperatingSystem = system.OperatingSystemMac
	runtime.ExecuteCommand = func(cmd string, args ...string) *exec.Cmd {
		assert.Equal(t, "open", cmd)
		assert.Equal(t, []string{TestURL}, args)
		return &exec.Cmd{}
	}
	runtime.Open(TestURL)
}

func TestOpenBrowserLinux(t *testing.T) {
	runtime := system.DefaultRuntime()
	runtime.OperatingSystem = system.OperatingSystemLinux
	runtime.ExecuteCommand = func(cmd string, args ...string) *exec.Cmd {
		assert.Equal(t, "xdg-open", cmd)
		assert.Equal(t, []string{TestURL}, args)
		return &exec.Cmd{}
	}
	runtime.Open(TestURL)
}

//-------------------------------------------------------------------------------------------------
