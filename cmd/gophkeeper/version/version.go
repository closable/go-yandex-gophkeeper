// Package version realize versionality server/client
package version

import (
	"fmt"
	"runtime"
)

var (
	// These variables are replaced by ldflags at build time
	buildVersion = "N/A"
	gitCommit    = "N/A"
	buildDate    = "N/A" // build date in ISO8601 format
)

// VersionInfo main version structure
type VersionInfo struct {
	BuildVersion string `json:"buildVersion" yaml:"buildVersion"`
	GitCommit    string `json:"gitCommit" yaml:"gitCommit"`
	BuildDate    string `json:"buildDate" yaml:"buildDate"`
	GoVersion    string `json:"goVersion" yaml:"goVersion"`
	Compiler     string `json:"compiler" yaml:"compiler"`
	Platform     string `json:"platform" yaml:"platform"`
}

// Get returns the overall codebase version. It's for detecting
// what code a binary was built from.
func Get() *VersionInfo {
	// These variables typically come from -ldflags settings and in
	// their absence fallback to the constants above
	return &VersionInfo{
		BuildVersion: buildVersion,
		GitCommit:    gitCommit,
		BuildDate:    buildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
