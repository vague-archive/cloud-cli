package system

import (
	"runtime"
)

//-------------------------------------------------------------------------------------------------

type OperatingSystem int

const (
	OperatingSystemUnknown OperatingSystem = iota
	OperatingSystemWindows
	OperatingSystemMac
	OperatingSystemLinux
)

func Current() OperatingSystem {
	return Identify(runtime.GOOS)
}

func Identify(name string) OperatingSystem {
	switch name {
	case "darwin":
		return OperatingSystemMac
	case "windows":
		return OperatingSystemWindows
	case "linux":
		return OperatingSystemLinux
	default:
		return OperatingSystemUnknown
	}
}

//-------------------------------------------------------------------------------------------------

func OpenCommand(OperatingSystem OperatingSystem, url string) (string, []string) {
	var cmd string
	var args []string
	switch OperatingSystem {
	case OperatingSystemMac:
		cmd = "open"
	case OperatingSystemWindows:
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler"}
	case OperatingSystemLinux:
		cmd = "xdg-open"
	}
	args = append(args, url)
	return cmd, args
}

//-------------------------------------------------------------------------------------------------
