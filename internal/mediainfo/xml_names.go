package mediainfo

import "strings"

func xmlFieldName(name string) string {
	var b strings.Builder
	b.Grow(len(name))
	for _, r := range name {
		switch {
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '_':
			b.WriteRune(r)
		case r == ' ' || r == '-' || r == '/' || r == '.':
			b.WriteRune('_')
		case r == '(' || r == ')':
			continue
		default:
			b.WriteRune('_')
		}
	}
	if b.Len() == 0 {
		return "Field"
	}
	return b.String()
}
