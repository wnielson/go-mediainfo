package mediainfo

import (
	"math"
	"testing"
)

func writeBitsVC1(dst []byte, bitPos *int, value uint32, n int) {
	for i := n - 1; i >= 0; i-- {
		b := (value >> uint(i)) & 1
		bytePos := *bitPos / 8
		shift := 7 - (*bitPos % 8)
		if b != 0 {
			dst[bytePos] |= 1 << uint(shift)
		}
		*bitPos++
	}
}

func TestParseVC1AnnexBMeta_AdvancedSequenceHeader(t *testing.T) {
	payload := make([]byte, 32)
	pos := 0

	writeBitsVC1(payload, &pos, 3, 2) // profile=Advanced
	writeBitsVC1(payload, &pos, 3, 3) // level=3
	writeBitsVC1(payload, &pos, 1, 2) // colordiff_format=1 (4:2:0)
	writeBitsVC1(payload, &pos, 0, 3) // frmrtq_postproc
	writeBitsVC1(payload, &pos, 0, 5) // bitrtq_postproc
	writeBitsVC1(payload, &pos, 0, 1) // postprocflag
	writeBitsVC1(payload, &pos, 959, 12)
	writeBitsVC1(payload, &pos, 539, 12)
	writeBitsVC1(payload, &pos, 0, 1) // pulldown
	writeBitsVC1(payload, &pos, 0, 1) // interlace
	writeBitsVC1(payload, &pos, 0, 1) // tfcntrflag
	writeBitsVC1(payload, &pos, 0, 1) // finterpflag
	writeBitsVC1(payload, &pos, 0, 1) // reserved
	writeBitsVC1(payload, &pos, 0, 1) // psf

	writeBitsVC1(payload, &pos, 1, 1)     // display_ext
	writeBitsVC1(payload, &pos, 1919, 14) // display_horiz_size (1920)
	writeBitsVC1(payload, &pos, 1079, 14) // display_vert_size (1080)

	writeBitsVC1(payload, &pos, 1, 1) // aspectratio_flag
	writeBitsVC1(payload, &pos, 1, 4) // aspect_ratio=1 (1:1)

	writeBitsVC1(payload, &pos, 1, 1) // framerate_flag
	writeBitsVC1(payload, &pos, 0, 1) // framerateind=0 (code form)
	writeBitsVC1(payload, &pos, 1, 8) // frameratenr=24000
	writeBitsVC1(payload, &pos, 2, 4) // frameratedr=1001

	writeBitsVC1(payload, &pos, 0, 1) // color_format_flag

	writeBitsVC1(payload, &pos, 1, 1) // hrd_param_flag
	writeBitsVC1(payload, &pos, 1, 5) // hrd_num_leaky_buckets=1
	writeBitsVC1(payload, &pos, 0, 4) // bitrate_exponent
	writeBitsVC1(payload, &pos, 5, 4) // buffer_size_exponent
	writeBitsVC1(payload, &pos, 0, 16)
	writeBitsVC1(payload, &pos, 58592, 16) // hrd_buffer => 3749952 bytes

	n := (pos + 7) / 8
	seq := append([]byte{0x00, 0x00, 0x01, 0x0F}, payload[:n]...)

	meta, ok := parseVC1AnnexBMeta(seq)
	if !ok {
		t.Fatalf("expected ok")
	}
	if meta.Profile != "Advanced" {
		t.Fatalf("Profile=%q", meta.Profile)
	}
	if meta.Level != 3 {
		t.Fatalf("Level=%d", meta.Level)
	}
	if meta.Width != 1920 || meta.Height != 1080 {
		t.Fatalf("size=%dx%d", meta.Width, meta.Height)
	}
	if meta.ChromaSubsampling != "4:2:0" {
		t.Fatalf("ChromaSubsampling=%q", meta.ChromaSubsampling)
	}
	if math.Abs(meta.PixelAspectRatio-1.0) > 1e-6 {
		t.Fatalf("PixelAspectRatio=%f", meta.PixelAspectRatio)
	}
	if meta.ScanType != "Progressive" {
		t.Fatalf("ScanType=%q", meta.ScanType)
	}
	if meta.BufferSize != 3749952 {
		t.Fatalf("BufferSize=%d", meta.BufferSize)
	}
	if meta.FrameRateNum != 24000 || meta.FrameRateDen != 1001 {
		t.Fatalf("FrameRate ratio=%d/%d", meta.FrameRateNum, meta.FrameRateDen)
	}
	if math.Abs(meta.FrameRate-23.976) > 0.002 {
		t.Fatalf("FrameRate=%f", meta.FrameRate)
	}
}
