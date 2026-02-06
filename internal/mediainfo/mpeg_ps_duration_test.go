package mediainfo

import "testing"

func TestFinalizeMPEGPSFallsBackToDerivedVideoDuration(t *testing.T) {
	streams := map[uint16]*psStream{
		psStreamKey(0xE0, psSubstreamNone): {
			id:              0xE0,
			subID:           psSubstreamNone,
			kind:            StreamVideo,
			format:          "MPEG Video",
			derivedDuration: 0.033,
		},
	}
	streamOrder := []uint16{psStreamKey(0xE0, psSubstreamNone)}

	info, _, ok := finalizeMPEGPS(streams, streamOrder, nil, ptsTracker{}, ptsTracker{}, 8<<10, mpegPSOptions{dvdParsing: true, parseSpeed: 0.5})
	if !ok {
		t.Fatalf("expected ok")
	}
	if info.DurationSeconds == 0 {
		t.Fatalf("expected DurationSeconds > 0")
	}
	if info.DurationSeconds < 0.032 || info.DurationSeconds > 0.034 {
		t.Fatalf("DurationSeconds = %f, want ~0.033", info.DurationSeconds)
	}
}
