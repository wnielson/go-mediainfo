package mediainfo

import "testing"

func TestParseH264SPS_ExtendedSAR(t *testing.T) {
	// SPS extracted from a real-world MP4 sample; official mediainfo reports PixelAspectRatio=1.001 (1001/1000).
	sps := []byte{
		0x67, 0x64, 0x00, 0x1e, 0xac, 0xd9, 0x40, 0xa0, 0x2f, 0xf9,
		0x7f, 0xf0, 0x50, 0x10, 0x50, 0x01, 0x00, 0x00, 0x03, 0x00,
		0x01, 0x00, 0x00, 0x03, 0x00, 0x28, 0x0f, 0x16, 0x2d, 0x96,
	}
	info := parseH264SPS(sps)
	if !info.HasSAR {
		t.Fatalf("expected HasSAR")
	}
	if info.SARWidth != 1281 || info.SARHeight != 1280 {
		t.Fatalf("sar=%d:%d, want 1281:1280", info.SARWidth, info.SARHeight)
	}
}
