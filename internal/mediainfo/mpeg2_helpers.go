package mediainfo

import (
	"fmt"
	"strings"
)

func scanMPEG2StartCodes(buf []byte, start int, fn func(i int, code byte) bool) {
	for i := start; i+4 <= len(buf); i++ {
		if buf[i] != 0x00 || buf[i+1] != 0x00 || buf[i+2] != 0x01 {
			continue
		}
		if !fn(i, buf[i+3]) {
			return
		}
	}
}

func nextStartCode(data []byte, start int) int {
	for i := start; i+3 < len(data); i++ {
		if data[i] == 0x00 && data[i+1] == 0x00 && data[i+2] == 0x01 {
			return i
		}
	}
	return -1
}

func maybeCaptureMPEG2Matrix(br *bitReader) string {
	matrix := make([]byte, 0, 64)
	for range 64 {
		value := br.readBitsValue(8)
		if len(matrix) < 64 {
			matrix = append(matrix, byte(value))
		}
	}
	if len(matrix) != 64 {
		return ""
	}
	var builder strings.Builder
	builder.Grow(128)
	for _, b := range matrix {
		fmt.Fprintf(&builder, "%02X", b)
	}
	return builder.String()
}
