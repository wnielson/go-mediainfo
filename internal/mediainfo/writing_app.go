package mediainfo

import "strings"

func normalizeWritingApplication(raw string) string {
	name, _, versionRaw := splitWritingApplication(raw)
	if name == "mkvmerge" && strings.HasPrefix(versionRaw, "v") {
		return name + " " + strings.TrimPrefix(versionRaw, "v")
	}
	return raw
}

func splitWritingApplication(raw string) (string, string, string) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", "", ""
	}
	parts := strings.Fields(raw)
	if len(parts) == 0 {
		return "", "", ""
	}
	name := parts[0]
	versionRaw := strings.TrimSpace(strings.TrimPrefix(raw, name))
	version := strings.TrimPrefix(versionRaw, "v")
	return name, version, versionRaw
}
