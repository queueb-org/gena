package app

import (
	"runtime/debug"
	"strings"
)

var _version string

// buildVersion returns application version taken from
// internally built info.
func buildVersion() (version string) {
	version = "devel"

	if info, ok := debug.ReadBuildInfo(); ok {
		version = moduleVersion(info.Main.Version)
	}

	return
}

// Module version parses build version and translates it to
// Human-Readable Version
func moduleVersion(version string) string {
	if version == "(devel)" || version == "" {
		return "devel"
	}

	return strings.TrimPrefix(version, "v")
}

// Version returns application version taken by default
func Version() string {
	if _version == "" {
		_version = moduleVersion(buildVersion())
	}

	return _version
}
