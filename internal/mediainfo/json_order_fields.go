package mediainfo

func orderFieldsForJSON(kind StreamKind, fields []Field) []Field {
	ordered := append([]Field(nil), fields...)
	sortFields(kind, ordered)
	return ordered
}
