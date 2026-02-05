package mediainfo

func consumeMPEG2HeaderBytes(entry *psStream, payload []byte, hasPTS bool) {
	if entry == nil || len(payload) == 0 {
		return
	}
	buf := append(entry.videoHeaderCarry, payload...)
	for i := 0; i+4 <= len(buf); i++ {
		if buf[i] != 0x00 || buf[i+1] != 0x00 || buf[i+2] != 0x01 {
			continue
		}
		switch buf[i+3] {
		case 0xB3:
			entry.videoHeaderBytes += 12
			if hasPTS {
				entry.videoSeqExtBytes += 12
			}
		case 0xB5:
			entry.videoHeaderBytes += 4
			if hasPTS {
				entry.videoSeqExtBytes += 4
			}
		case 0xB8:
			entry.videoHeaderBytes += 8
		case 0x00:
			entry.videoHeaderBytes += 6
		default:
			if buf[i+3] >= 0x01 && buf[i+3] <= 0xAF {
				entry.videoHeaderBytes += 6
			}
		}
	}
	if len(buf) >= 3 {
		entry.videoHeaderCarry = append(entry.videoHeaderCarry[:0], buf[len(buf)-3:]...)
	} else {
		entry.videoHeaderCarry = append(entry.videoHeaderCarry[:0], buf...)
	}
}

func mpeg2HeaderSize(code byte) uint64 {
	switch {
	case code == 0xB3:
		return 12
	case code == 0xB5:
		return 4
	case code == 0xB8:
		return 8
	case code == 0x00:
		return 6
	case code >= 0x01 && code <= 0xAF:
		return 6
	default:
		return 0
	}
}
