package mediainfo

import (
	"fmt"
	"math"
)

func formatPixels(value uint64) string {
	if value == 0 {
		return ""
	}
	return formatThousands(int64(value)) + " pixels"
}

func formatChannels(value uint64) string {
	if value == 0 {
		return ""
	}
	if value == 1 {
		return "1 channel"
	}
	return fmt.Sprintf("%d channels", value)
}

func formatSampleRate(rate float64) string {
	if rate <= 0 {
		return ""
	}
	if rate >= 1000 {
		return fmt.Sprintf("%.1f kHz", rate/1000)
	}
	return fmt.Sprintf("%.0f Hz", rate)
}

func formatBitDepth(bits uint8) string {
	if bits == 0 {
		return ""
	}
	return fmt.Sprintf("%d bits", bits)
}

func formatAspectRatio(width, height uint64) string {
	if width == 0 || height == 0 {
		return ""
	}
	g := gcd(width, height)
	reducedW := width / g
	reducedH := height / g
	if reducedW <= 50 && reducedH <= 50 {
		return fmt.Sprintf("%d:%d", reducedW, reducedH)
	}
	ratio := float64(width) / float64(height)
	common := []float64{1.33, 1.37, 1.66, 1.78, 1.85, 2.00, 2.20, 2.40}
	for _, target := range common {
		if math.Abs(ratio-target) < 0.02 {
			return fmt.Sprintf("%.2f:1", target)
		}
	}
	return fmt.Sprintf("%.2f:1", ratio)
}

func gcd(a, b uint64) uint64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func formatBitsPerPixelFrame(bitrate float64, width, height uint64, fps float64) string {
	if bitrate <= 0 || width == 0 || height == 0 || fps <= 0 {
		return ""
	}
	value := bitrate / (float64(width) * float64(height) * fps)
	return fmt.Sprintf("%.3f", value)
}

func formatStreamSize(bytes int64, total int64) string {
	if bytes <= 0 || total <= 0 {
		return ""
	}
	percent := int(math.Round(float64(bytes) * 100 / float64(total)))
	return fmt.Sprintf("%s (%d%%)", formatBytes(bytes), percent)
}
