package version

import "fmt"

type Info struct {
	Version   string
	GitCommit string
	BuildDate string
	GoVersion string
}

// These variables are set via ldflags at build time
var (
	version   = "dev"
	gitCommit = "unknown"
	buildDate = "unknown"
	goVersion = "unknown"
)

// Get returns the populated Info struct
func Get() Info {
	return Info{
		Version:   version,
		GitCommit: gitCommit,
		BuildDate: buildDate,
		GoVersion: goVersion,
	}
}

func (i Info) String() string {
	return fmt.Sprintf("%s (%s) built on %s with %s",
		i.Version, i.GitCommit, i.BuildDate, i.GoVersion)
}
