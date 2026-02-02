package mediainfo

import "fmt"

func formatBytes(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div := float64(size)
	exp := 0
	units := []string{"KiB", "MiB", "GiB", "TiB", "PiB"}
	for div >= unit && exp < len(units)-1 {
		div /= unit
		exp++
	}
	return fmt.Sprintf("%.2f %s", div, units[exp])
}
