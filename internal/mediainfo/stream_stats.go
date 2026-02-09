package mediainfo

import (
	"math"
	"strconv"
)

func frameCountFromFields(fields []Field) (string, bool) {
	duration, durOk := parseDurationSeconds(findField(fields, "Duration"))
	fps, fpsOk := parseFPS(findField(fields, "Frame rate"))
	if !durOk || !fpsOk {
		return "", false
	}
	return strconv.Itoa(int(math.Round(duration * fps))), true
}

func sumStreamSizes(streams []Stream, includeFieldFallback bool) int64 {
	var sum int64
	for _, stream := range streams {
		if stream.JSON != nil {
			if value, ok := stream.JSON["StreamSize"]; ok {
				if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
					sum += parsed
				}
			}
			continue
		}
		if includeFieldFallback {
			if sizeValue := findField(stream.Fields, "Stream size"); sizeValue != "" {
				if parsed, ok := parseSizeBytes(sizeValue); ok {
					sum += parsed
				}
			}
		}
	}
	return sum
}

func setRemainingStreamSize(json map[string]string, total int64, streamSizeSum int64) {
	if streamSizeSum <= 0 {
		return
	}
	remaining := total - streamSizeSum
	if remaining >= 0 {
		json["StreamSize"] = strconv.FormatInt(remaining, 10)
	}
}

func setOverallBitRate(json map[string]string, size int64, duration float64) {
	if duration <= 0 {
		return
	}
	// Match MediaInfo: bitrate computations are effectively based on integer milliseconds.
	durationMs := int64(math.Round(duration * 1000))
	if durationMs <= 0 {
		return
	}
	overall := (size*8000 + durationMs/2) / durationMs
	if overall > 0 {
		json["OverallBitRate"] = strconv.FormatInt(overall, 10)
	}
}
