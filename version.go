package main

import (
	"runtime/debug"
)

const defaultVersion = "development"

var version = ""

func versionStr() string {
	// 1. if version is set (by goreleaser for example), show it
	if version != "" {
		return version
	}

	// 2. if runtime version is set (by go get for example), show it
	info, ok := debug.ReadBuildInfo()
	if ok {
		return info.Main.Version
	}

	// 3. show default version
	return defaultVersion
}
