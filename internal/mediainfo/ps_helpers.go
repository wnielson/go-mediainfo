package mediainfo

func findPSStreamByKind(streams map[byte]psStream, kind StreamKind) (psStream, bool) {
	for _, st := range streams {
		if st.kind == kind {
			return st, true
		}
	}
	return psStream{}, false
}
