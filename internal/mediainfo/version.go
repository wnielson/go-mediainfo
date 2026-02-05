package mediainfo

import "strings"

const (
	AppName = "go-mediainfo"
	AppURL  = "https://github.com/autobrr/go-mediainfo"
)

var AppVersion = "dev"

func SetAppVersion(version string) {
	if version != "" {
		AppVersion = version
	}
}

func FormatVersion(version string) string {
	if version == "" || version == "dev" {
		return "dev"
	}
	if strings.HasPrefix(version, "v") {
		return version
	}
	return "v" + version
}
