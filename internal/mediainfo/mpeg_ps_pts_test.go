package mediainfo

import "testing"

func TestPtsDurationPS_Resets(t *testing.T) {
	var tracker ptsTracker
	tracker.add(0)
	tracker.add(90_000 * 100) // 100 seconds
	tracker.add(0)            // discontinuity
	tracker.add(90_000 * 10)  // 10 seconds

	t.Run("parseSpeed<1 uses total", func(t *testing.T) {
		got := ptsDurationPS(tracker, mpegPSOptions{parseSpeed: 0.5})
		if got != 110 {
			t.Fatalf("duration=%v, want 110", got)
		}
	})

	t.Run("parseSpeed>=1 uses minmax span (non-dvd)", func(t *testing.T) {
		got := ptsDurationPS(tracker, mpegPSOptions{parseSpeed: 1})
		if got != 100 {
			t.Fatalf("duration=%v, want 100", got)
		}
	})

	t.Run("dvdParsing uses last segment", func(t *testing.T) {
		got := ptsDurationPS(tracker, mpegPSOptions{parseSpeed: 1, dvdParsing: true})
		if got != 10 {
			t.Fatalf("duration=%v, want 10", got)
		}
	})
}
