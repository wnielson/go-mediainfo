package mediainfo

import (
	"fmt"
	"strings"
)

var languageNames = map[string]string{
	"en": "English",
	"fr": "French",
	"es": "Spanish",
	"de": "German",
	"it": "Italian",
	"pt": "Portuguese",
}

var languageMap3To2 = map[string]string{
	"eng": "en",
	"fra": "fr",
	"fre": "fr",
	"spa": "es",
	"deu": "de",
	"ger": "de",
	"ita": "it",
	"por": "pt",
}

func normalizeLanguageCode(code string) string {
	code = strings.TrimSpace(code)
	if code == "" {
		return ""
	}
	code = strings.ReplaceAll(code, "_", "-")
	parts := strings.Split(code, "-")
	if len(parts) == 0 {
		return code
	}
	lang := strings.ToLower(parts[0])
	if lang == "und" {
		return ""
	}
	if mapped, ok := languageMap3To2[lang]; ok {
		lang = mapped
	}
	if len(parts) > 1 {
		region := strings.ToUpper(parts[1])
		return fmt.Sprintf("%s-%s", lang, region)
	}
	return lang
}

func formatLanguage(code string) string {
	normalized := normalizeLanguageCode(code)
	if normalized == "" {
		return ""
	}
	parts := strings.Split(normalized, "-")
	name := languageNames[parts[0]]
	if name == "" {
		return code
	}
	if len(parts) > 1 {
		return fmt.Sprintf("%s (%s)", name, strings.ToUpper(parts[1]))
	}
	return name
}
