package cli

import (
	"fmt"
	"io"
)

var appVersion = "dev"

func SetVersion(version string) {
	if version != "" {
		appVersion = version
	}
}

func Version(stdout io.Writer) {
	fmt.Fprintf(stdout, "MediaInfo Command line, %s\n", appVersion)
}
