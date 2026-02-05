package cli

import (
	"fmt"
	"io"

	"github.com/autobrr/go-mediainfo/internal/mediainfo"
)

var appVersion = "dev"

func SetVersion(version string) {
	if version != "" {
		appVersion = version
	}
}

func Version(stdout io.Writer) {
	fmt.Fprintf(stdout, "go-mediainfo, %s\n", mediainfo.FormatVersion(appVersion))
}
