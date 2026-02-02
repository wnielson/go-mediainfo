package mediainfo

func addStreamDuration(fields []Field, duration float64) []Field {
	if duration <= 0 {
		return fields
	}
	return appendFieldUnique(fields, Field{Name: "Duration", Value: formatDuration(duration)})
}

func addStreamBitrate(fields []Field, bits float64) []Field {
	if bits <= 0 {
		return fields
	}
	return appendFieldUnique(fields, Field{Name: "Bit rate", Value: formatBitrate(bits)})
}
