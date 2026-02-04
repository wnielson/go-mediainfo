package mediainfo

import "testing"

func TestPtsTrackerSegments(t *testing.T) {
	var tracker ptsTracker
	tracker.add(1000)
	tracker.add(2000)
	tracker.add(1500) // reset
	tracker.add(2500)

	if !tracker.hasResets() {
		t.Fatalf("expected reset tracking")
	}

	if got, want := tracker.duration(), float64(1500)/90000.0; got != want {
		t.Fatalf("duration mismatch: got %v want %v", got, want)
	}
	if got, want := tracker.durationLastSegment(), float64(1000)/90000.0; got != want {
		t.Fatalf("last segment mismatch: got %v want %v", got, want)
	}
	if got, want := tracker.durationTotal(), float64(2000)/90000.0; got != want {
		t.Fatalf("total mismatch: got %v want %v", got, want)
	}
}
