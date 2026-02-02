package mediainfo

import (
	"bytes"
	"strings"
)

func RenderText(reports []Report) string {
	var buf bytes.Buffer
	for i, report := range reports {
		if i > 0 {
			buf.WriteString("\n")
		}
		writeStream(&buf, report.General)
		for _, stream := range report.Streams {
			buf.WriteString("\n")
			writeStream(&buf, stream)
		}
	}
	return strings.TrimRight(buf.String(), "\n")
}

func writeStream(buf *bytes.Buffer, stream Stream) {
	buf.WriteString(string(stream.Kind))
	buf.WriteString("\n")
	for _, field := range stream.Fields {
		buf.WriteString(padRight(field.Name, 36))
		buf.WriteString(": ")
		buf.WriteString(field.Value)
		buf.WriteString("\n")
	}
}

func padRight(value string, width int) string {
	if len(value) >= width {
		return value
	}
	return value + strings.Repeat(" ", width-len(value))
}
