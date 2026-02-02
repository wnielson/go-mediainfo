package mediainfo

import (
	"fmt"
	"math"
)

func formatFrameRate(rate float64) string {
	if rate <= 0 {
		return ""
	}
	if math.Abs(rate-math.Round(rate)) < 0.0005 {
		return fmt.Sprintf("%.0f FPS", rate)
	}
	return fmt.Sprintf("%.3f FPS", rate)
}
