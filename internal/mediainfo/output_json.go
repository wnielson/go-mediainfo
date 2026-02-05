package mediainfo

import (
	"bytes"
	"encoding/json"
)

type jsonMediaOut struct {
	Ref    string
	Tracks []jsonTrackOut
}

type jsonTrackOut struct {
	Fields []jsonKV
}

func RenderJSON(reports []Report) string {
	if len(reports) == 1 {
		return renderJSONPayload(buildJSONPayload(reports[0])) + "\n"
	}
	payloads := make([]jsonPayloadOut, 0, len(reports))
	for _, report := range reports {
		payloads = append(payloads, buildJSONPayload(report))
	}
	return renderJSONPayloads(payloads) + "\n"
}

type jsonPayloadOut struct {
	CreatingLibrary []jsonKV
	Media           jsonMediaOut
}

func buildJSONPayload(report Report) jsonPayloadOut {
	return jsonPayloadOut{
		CreatingLibrary: jsonCreatingLibraryFields(),
		Media:           buildJSONMedia(report),
	}
}

func jsonCreatingLibraryFields() []jsonKV {
	return []jsonKV{
		{Key: "name", Val: AppName},
		{Key: "version", Val: FormatVersion(AppVersion)},
		{Key: "url", Val: AppURL},
	}
}

func renderJSONPayload(payload jsonPayloadOut) string {
	var buf bytes.Buffer
	buf.WriteString("{\n")
	writeJSONField(&buf, "creatingLibrary", renderJSONObject(payload.CreatingLibrary, false), true)
	buf.WriteString(",\n")
	writeJSONField(&buf, "media", renderJSONMedia(payload.Media), true)
	buf.WriteString("\n}")
	return buf.String()
}

func renderJSONPayloads(payloads []jsonPayloadOut) string {
	var buf bytes.Buffer
	buf.WriteString("[\n")
	for i, payload := range payloads {
		if i > 0 {
			buf.WriteString(",\n")
		}
		buf.WriteString(renderJSONPayload(payload))
	}
	buf.WriteString("\n]")
	return buf.String()
}

func renderJSONMedia(media jsonMediaOut) string {
	tracks := make([]string, 0, len(media.Tracks))
	for _, track := range media.Tracks {
		tracks = append(tracks, renderJSONTrack(track.Fields))
	}
	return renderJSONMediaObject(media.Ref, tracks)
}

func renderJSONMediaObject(ref string, tracks []string) string {
	var buf bytes.Buffer
	buf.WriteString("{")
	writeJSONField(&buf, "@ref", ref, false)
	buf.WriteString(",")
	writeJSONField(&buf, "track", renderJSONArray(tracks, false), true)
	buf.WriteString("}")
	return buf.String()
}

func renderJSONArray(items []string, multiline bool) string {
	var buf bytes.Buffer
	buf.WriteString("[")
	for i, item := range items {
		if i > 0 {
			if multiline {
				buf.WriteString(",\n")
			} else {
				buf.WriteString(",")
			}
		}
		buf.WriteString(item)
	}
	buf.WriteString("]")
	return buf.String()
}

func renderJSONTrack(fields []jsonKV) string {
	var buf bytes.Buffer
	buf.WriteString("{")
	inlineCount := 2
	if len(fields) > 2 && fields[1].Key == "@typeorder" && fields[2].Key == "StreamOrder" {
		inlineCount = 3
	}
	if len(fields) > 2 && fields[1].Key == "@typeorder" && fields[2].Key == "ID" {
		inlineCount = 3
	}
	for i, field := range fields {
		if i > 0 {
			if i < inlineCount {
				buf.WriteString(",")
			} else {
				buf.WriteString(",\n")
			}
		}
		writeJSONField(&buf, field.Key, field.Val, field.Raw)
	}
	buf.WriteString("}")
	return buf.String()
}

func renderJSONObject(fields []jsonKV, multiline bool) string {
	var buf bytes.Buffer
	buf.WriteString("{")
	for i, field := range fields {
		if i > 0 {
			if multiline {
				buf.WriteString(",\n")
			} else {
				buf.WriteString(",")
			}
		}
		writeJSONField(&buf, field.Key, field.Val, field.Raw)
	}
	buf.WriteString("}")
	return buf.String()
}

func writeJSONField(buf *bytes.Buffer, key, value string, raw bool) {
	buf.WriteString("\"")
	buf.WriteString(key)
	buf.WriteString("\":")
	if raw {
		buf.WriteString(value)
		return
	}
	buf.WriteString(renderJSONString(value))
}

func renderJSONString(value string) string {
	data, _ := json.Marshal(value)
	return string(data)
}
