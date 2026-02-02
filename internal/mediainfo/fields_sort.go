package mediainfo

import "sort"

func sortFields(kind StreamKind, fields []Field) {
	order := streamFieldOrder
	if kind == StreamGeneral {
		order = generalFieldOrder
	}

	sort.SliceStable(fields, func(i, j int) bool {
		ai, aok := order[fields[i].Name]
		aj, bok := order[fields[j].Name]
		switch {
		case aok && bok:
			return ai < aj
		case aok:
			return true
		case bok:
			return false
		default:
			return fields[i].Name < fields[j].Name
		}
	})
}
