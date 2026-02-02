package mediainfo

import (
	"bytes"
	"strings"
)

func RenderCSV(reports []Report) string {
	var buf bytes.Buffer
	buf.WriteString("ref,track_type,field,value\n")
	for _, report := range reports {
		writeCSVTrack(&buf, report.Ref, report.General)
		for _, stream := range orderTracks(report.Streams) {
			writeCSVTrack(&buf, report.Ref, stream)
		}
	}
	return buf.String()
}

func writeCSVTrack(buf *bytes.Buffer, ref string, stream Stream) {
	for _, field := range orderFieldsForJSON(stream.Kind, stream.Fields) {
		buf.WriteString(csvEscape(ref))
		buf.WriteString(",")
		buf.WriteString(csvEscape(string(stream.Kind)))
		buf.WriteString(",")
		buf.WriteString(csvEscape(field.Name))
		buf.WriteString(",")
		buf.WriteString(csvEscape(field.Value))
		buf.WriteString("\n")
	}
}

func csvEscape(value string) string {
	if value == "" {
		return ""
	}
	if strings.ContainsAny(value, ",\n\"") {
		return "\"" + strings.ReplaceAll(value, "\"", "\"\"") + "\""
	}
	return value
}
