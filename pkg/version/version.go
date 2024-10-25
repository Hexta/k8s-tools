package version

import "runtime/debug"

var defaultVersion = "dev"

func Version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return defaultVersion
	}

	return info.Main.Version
}
