package mediainfo

import (
	"bytes"
	"testing"
)

const fuzzParserMaxBytes = 1 << 20 // 1 MiB

func fuzzLimit(data []byte) []byte {
	if len(data) > fuzzParserMaxBytes {
		return data[:fuzzParserMaxBytes]
	}
	return data
}

func FuzzParseAC3FrameParsers(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{0x0B, 0x77})
	f.Add([]byte{0x0B, 0x77, 0x00, 0x00, 0x10, 0x00, 0x00})
	f.Add([]byte{0x0B, 0x77, 0x77, 0x00, 0x00, 0x40, 0x50, 0x00})

	f.Fuzz(func(t *testing.T, data []byte) {
		data = fuzzLimit(data)
		_, _, _ = parseAC3Frame(data)
		_, _, _ = parseEAC3FrameWithOptions(data, true)
		_, _, _ = parseEAC3FrameWithOptions(data, false)
	})
}

func FuzzParseMPEGTSPacketizers(f *testing.F) {
	f.Add([]byte{}, uint8(128))
	f.Add(bytes.Repeat([]byte{0x47}, 188), uint8(128))
	f.Add(bytes.Repeat([]byte{0x47}, 192), uint8(128))
	f.Add(append([]byte{0x47, 0x40, 0x00, 0x10}, bytes.Repeat([]byte{0xFF}, 184)...), uint8(160))

	f.Fuzz(func(t *testing.T, data []byte, speed uint8) {
		data = fuzzLimit(data)
		parseSpeed := float64(speed) / 255.0

		tsReader := bytes.NewReader(data)
		_, _, _, _ = ParseMPEGTS(tsReader, int64(len(data)), parseSpeed)

		bdavReader := bytes.NewReader(data)
		_, _, _, _ = ParseBDAV(bdavReader, int64(len(data)), parseSpeed)
	})
}

func FuzzParseMatroskaContainers(f *testing.F) {
	f.Add([]byte{}, uint8(128))
	f.Add([]byte{0x1A, 0x45, 0xDF, 0xA3}, uint8(128))
	f.Add([]byte{0x1A, 0x45, 0xDF, 0xA3, 0x9F, 0x42, 0x86, 0x81, 0x01}, uint8(128))

	f.Fuzz(func(t *testing.T, data []byte, speed uint8) {
		data = fuzzLimit(data)
		opts := defaultAnalyzeOptions()
		opts.HasParseSpeed = true
		opts.ParseSpeed = float64(speed) / 255.0
		opts = normalizeAnalyzeOptions(opts)
		_, _ = ParseMatroskaWithOptions(bytes.NewReader(data), int64(len(data)), opts)
	})
}
