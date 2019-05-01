package version

import (
	"fmt"
)

// Version variables (cf. https://github.com/ahmetb/govvv)
var (
	GitCommit  string
	GitBranch  string
	GitState   string
	GitSummary string
	BuildDate  string
	Version    string
)

// Long computes long version
// eg: v0.0.1 master#b1e8b63b (2017-10-27T08:41:03Z)
func Long() string {
	result := Short()

	if BuildDate != "" {
		result += fmt.Sprintf(" (%s)", BuildDate)
	}

	return result
}

// Short computes short version
// eg: v0.0.1 master#b1e8b63b
func Short() string {
	result := Version

	if GitBranch != "" {
		if result != "" {
			result += " "
		}

		result += GitBranch
		if GitCommit != "" {
			result += fmt.Sprintf("#%s", GitCommit)
		}
	}

	if result == "" {
		return "unknown"
	}

	return result
}
