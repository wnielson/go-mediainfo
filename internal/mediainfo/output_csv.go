package mediainfo

import (
	"bytes"
	"fmt"
)

func RenderCSV(reports []Report) string {
	var buf bytes.Buffer
	for _, report := range reports {
		writeCSVTrack(&buf, string(report.General.Kind), report.General)

		kindCounts := map[StreamKind]int{}
		for _, stream := range report.Streams {
			kindCounts[stream.Kind]++
		}
		kindIndex := map[StreamKind]int{}
		for _, stream := range report.Streams {
			kindIndex[stream.Kind]++
			title := csvStreamTitle(stream.Kind, kindIndex[stream.Kind], kindCounts[stream.Kind])
			writeCSVTrack(&buf, title, stream)
		}
	}
	return buf.String()
}

func csvStreamTitle(kind StreamKind, index int, total int) string {
	if total > 1 {
		return fmt.Sprintf("%s,%d", kind, index)
	}
	return string(kind)
}

func writeCSVTrack(buf *bytes.Buffer, trackType string, stream Stream) {
	buf.WriteString(trackType)
	buf.WriteString("\n")
	for _, field := range stream.Fields {
		buf.WriteString(field.Name)
		buf.WriteString(",")
		buf.WriteString(field.Value)
		buf.WriteString("\n")
	}
	buf.WriteString("\n")
}
