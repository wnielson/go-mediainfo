package mediainfo

import "math"

func estimateTSFrameRate(sampleCount uint64, duration float64) string {
	if duration <= 0 || sampleCount == 0 {
		return ""
	}
	rate := float64(sampleCount) / duration
	if math.IsInf(rate, 0) || math.IsNaN(rate) {
		return ""
	}
	return formatFrameRate(rate)
}
