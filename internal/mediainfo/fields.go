package mediainfo

func appendFieldUnique(fields []Field, field Field) []Field {
	for _, existing := range fields {
		if existing.Name == field.Name {
			return fields
		}
	}
	return append(fields, field)
}
