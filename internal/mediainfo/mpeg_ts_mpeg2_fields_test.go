package mediainfo

import "testing"

func TestMPEG2TSGOPValueVariableWins(t *testing.T) {
	info := mpeg2VideoInfo{
		ScanType:    "Interlaced",
		GOPVariable: true,
		GOPM:        3,
		GOPN:        13,
		GOPLength:   13,
	}
	if got := mpeg2TSGOPValue(info); got != "Variable" {
		t.Fatalf("mpeg2TSGOPValue()=%q, want %q", got, "Variable")
	}
}

func TestMPEG2TSGOPValueInterlacedMN(t *testing.T) {
	info := mpeg2VideoInfo{
		ScanType: "Interlaced",
		GOPM:     3,
		GOPN:     12,
	}
	if got := mpeg2TSGOPValue(info); got != "M=3, N=12" {
		t.Fatalf("mpeg2TSGOPValue()=%q, want %q", got, "M=3, N=12")
	}
}
