package version

import "fmt"

// Set via ldflags at build time.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func Full() string {
	return fmt.Sprintf("goforge %s (commit: %s, built: %s)", Version, Commit, Date)
}
