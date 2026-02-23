package mediainfo_test

import (
	"testing"

	"github.com/autobrr/go-mediainfo/pkg/mediainfo"
)

func TestProxyAPI(t *testing.T) {
	// Smoke test to ensure the proxy can be imported and types are consistent
	var _ mediainfo.Report
	var _ mediainfo.StreamKind = mediainfo.StreamGeneral
}
