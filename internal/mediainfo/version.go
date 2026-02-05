package mediainfo

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
