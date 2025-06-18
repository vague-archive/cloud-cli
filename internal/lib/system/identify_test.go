package system_test

import (
	"runtime"
	"testing"

	"github.com/vaguevoid/cloud-cli/internal/lib/system"
	"github.com/vaguevoid/cloud-cli/internal/test/assert"
)

//-------------------------------------------------------------------------------------------------

func TestIdentify(t *testing.T) {
	assert.Equal(t, system.OperatingSystemMac, system.Identify("darwin"))
	assert.Equal(t, system.OperatingSystemWindows, system.Identify("windows"))
	assert.Equal(t, system.OperatingSystemLinux, system.Identify("linux"))
	assert.Equal(t, system.OperatingSystemUnknown, system.Identify("c64"))
}

//-------------------------------------------------------------------------------------------------

func TestCurrent(t *testing.T) {
	goos := runtime.GOOS
	switch goos {
	case "darwin":
		assert.Equal(t, system.OperatingSystemMac, system.Current())
	case "windows":
		assert.Equal(t, system.OperatingSystemWindows, system.Current())
	case "linux":
		assert.Equal(t, system.OperatingSystemLinux, system.Current())
	default:
		assert.Fail(t, "unexpected os", goos)
	}
}

//-------------------------------------------------------------------------------------------------

func TestOpenCommand(t *testing.T) {
	url := "https://play.void.dev"

	cmd, args := system.OpenCommand(system.OperatingSystemLinux, url)
	assert.Equal(t, "xdg-open", cmd)
	assert.Equal(t, []string{url}, args)

	cmd, args = system.OpenCommand(system.OperatingSystemMac, url)
	assert.Equal(t, "open", cmd)
	assert.Equal(t, []string{url}, args)

	cmd, args = system.OpenCommand(system.OperatingSystemWindows, url)
	assert.Equal(t, "rundll32", cmd)
	assert.Equal(t, []string{"url.dll,FileProtocolHandler", url}, args)
}

//-------------------------------------------------------------------------------------------------
