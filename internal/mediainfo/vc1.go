package mediainfo

import "bytes"

type vc1Meta struct {
	Profile           string
	Level             int
	Width             uint64
	Height            uint64
	ChromaSubsampling string
	PixelAspectRatio  float64
	ScanType          string
	BufferSize        int64
	FrameRate         float64
	FrameRateNum      int
	FrameRateDen      int
}

func vc1FrameRateENR(code uint8) int {
	switch code {
	case 0x01:
		return 24000
	case 0x02:
		return 25000
	case 0x03:
		return 30000
	case 0x04:
		return 50000
	case 0x05:
		return 60000
	case 0x06:
		return 48000
	case 0x07:
		return 72000
	default:
		return 0
	}
}

func vc1FrameRateDR(code uint8) int {
	switch code {
	case 0x01:
		return 1000
	case 0x02:
		return 1001
	default:
		return 0
	}
}

var vc1PixelAspectRatio = []float64{
	1,         // Reserved
	1,         // 1:1
	12.0 / 11, // 12:11
	10.0 / 11, // 10:11
	16.0 / 11, // 16:11
	40.0 / 33, // 40:33
	24.0 / 11, // 24:11
	20.0 / 11, // 20:11
	32.0 / 11, // 32:11
	80.0 / 33, // 80:33
	18.0 / 11, // 18:11
	15.0 / 11, // 15:11
	64.0 / 33, // 64:33
	160.0 / 99,
	1, // Reserved
	1, // Custom
}

// parseVC1AnnexBMeta extracts a small subset of VC-1 sequence header fields as emitted by
// official MediaInfo for Blu-ray/TS.
func parseVC1AnnexBMeta(data []byte) (vc1Meta, bool) {
	// SequenceHeader start code.
	start := bytes.Index(data, []byte{0x00, 0x00, 0x01, 0x0F})
	if start < 0 || start+4 >= len(data) {
		return vc1Meta{}, false
	}
	br := newBitReader(data[start+4:])

	profile := int(br.readBitsValue(2))
	if profile != 3 { // Advanced
		return vc1Meta{}, false
	}

	level := int(br.readBitsValue(3))
	colordiffFormat := int(br.readBitsValue(2))
	_ = br.readBitsValue(3) // frmrtq_postproc
	_ = br.readBitsValue(5) // bitrtq_postproc
	_ = br.readBitsValue(1) // postprocflag
	codedWidth := int(br.readBitsValue(12))
	codedHeight := int(br.readBitsValue(12))
	_ = br.readBitsValue(1) // pulldown
	interlace := int(br.readBitsValue(1))
	_ = br.readBitsValue(1) // tfcntrflag
	_ = br.readBitsValue(1) // finterpflag
	_ = br.readBitsValue(1) // reserved
	_ = br.readBitsValue(1) // psf

	meta := vc1Meta{
		Profile: "Advanced",
		Level:   level,
		Width:   uint64((codedWidth + 1) * 2),
		Height:  uint64((codedHeight + 1) * 2),
	}

	if interlace == 1 {
		meta.ScanType = "Interlaced"
	} else {
		meta.ScanType = "Progressive"
	}

	if colordiffFormat == 1 {
		meta.ChromaSubsampling = "4:2:0"
	}

	displayExt := br.readBitsValue(1) // display_ext
	if displayExt == 1 {
		_ = br.readBitsValue(14) // display_horiz_size
		_ = br.readBitsValue(14) // display_vert_size

		aspectRatioFlag := br.readBitsValue(1)
		if aspectRatioFlag == 1 {
			arCode := int(br.readBitsValue(4))
			if arCode >= 0 && arCode < len(vc1PixelAspectRatio) {
				meta.PixelAspectRatio = vc1PixelAspectRatio[arCode]
			}
			if arCode == 0x0F {
				arX := float64(br.readBitsValue(8))
				arY := float64(br.readBitsValue(8))
				if arX > 0 && arY > 0 {
					meta.PixelAspectRatio = arX / arY
				}
			}
		}

		frPresent := br.readBitsValue(1)
		if frPresent == 1 {
			frForm := br.readBitsValue(1) // framerateind
			if frForm == 1 {
				frExp := br.readBitsValue(16)
				meta.FrameRate = float64(frExp+1) / 64.0
			} else {
				enr := uint8(br.readBitsValue(8))
				dr := uint8(br.readBitsValue(4))
				num := vc1FrameRateENR(enr)
				den := vc1FrameRateDR(dr)
				if num > 0 && den > 0 {
					meta.FrameRateNum = num
					meta.FrameRateDen = den
					meta.FrameRate = float64(num) / float64(den)
				}
			}
		}

		colorFormatFlag := br.readBitsValue(1)
		if colorFormatFlag == 1 {
			_ = br.readBitsValue(8) // color_prim
			_ = br.readBitsValue(8) // transfer_char
			_ = br.readBitsValue(8) // matrix_coef
		}
	}

	hrdParamFlag := br.readBitsValue(1)
	if hrdParamFlag == 1 {
		hrdBuckets := int(br.readBitsValue(5))
		_ = br.readBitsValue(4) // bitrate_exponent
		bufExp := int(br.readBitsValue(4))
		for i := 0; i < hrdBuckets; i++ {
			_ = br.readBitsValue(16) // hrd_rate
			hrdBuf := int64(br.readBitsValue(16))
			val := (hrdBuf + 1) * (1 << (1 + bufExp))
			if val > meta.BufferSize {
				meta.BufferSize = val
			}
		}
	}

	return meta, true
}
