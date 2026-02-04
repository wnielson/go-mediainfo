package mediainfo

import "testing"

func TestParseDVDUserDataOddField(t *testing.T) {
	data := []byte{
		'C', 'C', 0x01, 0xF8, 0x82,
		0xFF, 0x14, 0x2F,
		0xFE, 0x80, 0x80,
	}
	hasCC, ccType, hasCommand, hasDisplay := parseDVDUserData(data)
	if !hasCC {
		t.Fatalf("expected CC data")
	}
	if ccType != 1 {
		t.Fatalf("expected CC3 (odd field), got %d", ccType)
	}
	if !hasCommand {
		t.Fatalf("expected command detection")
	}
	if !hasDisplay {
		t.Fatalf("expected display detection")
	}
}

func TestParseDVDUserDataEvenField(t *testing.T) {
	data := []byte{
		'C', 'C', 0x01, 0xF8, 0x82,
		0xFF, 0x80, 0x80,
		0xFE, 0x14, 0x2F,
	}
	hasCC, ccType, hasCommand, hasDisplay := parseDVDUserData(data)
	if !hasCC {
		t.Fatalf("expected CC data")
	}
	if ccType != 0 {
		t.Fatalf("expected CC1 (even field), got %d", ccType)
	}
	if !hasCommand {
		t.Fatalf("expected command detection")
	}
	if !hasDisplay {
		t.Fatalf("expected display detection")
	}
}
