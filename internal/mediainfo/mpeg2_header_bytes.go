package mediainfo

func consumeMPEG2HeaderBytes(entry *psStream, payload []byte, hasPTS bool) {
	if entry == nil || len(payload) == 0 {
		return
	}
	entry.videoHeaderCarry = append(entry.videoHeaderCarry, payload...)
	buf := entry.videoHeaderCarry
	scanMPEG2StartCodes(buf, 0, func(_ int, code byte) bool {
		switch code {
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
			entry.videoGOPBytes += 8
		case 0x00:
			entry.videoHeaderBytes += 6
		default:
			if code >= 0x01 && code <= 0xAF {
				entry.videoHeaderBytes += 6
			}
		}
		return true
	})
	if len(buf) >= 3 {
		entry.videoHeaderCarry = append(entry.videoHeaderCarry[:0], buf[len(buf)-3:]...)
	} else {
		entry.videoHeaderCarry = append(entry.videoHeaderCarry[:0], buf...)
	}
}

func consumeMPEG2FrameBytes(entry *psStream, payload []byte) {
	if entry == nil || len(payload) == 0 || entry.videoFrameBytesCount >= 2 {
		if entry != nil && entry.videoFrameBytesCount < 2 {
			entry.videoFramePos += int64(len(payload))
		}
		return
	}
	entry.videoFrameCarry = append(entry.videoFrameCarry, payload...)
	buf := entry.videoFrameCarry
	basePos := entry.videoFramePos - int64(len(entry.videoFrameCarry))
	scanMPEG2StartCodes(buf, 0, func(i int, code byte) bool {
		if code != 0x00 {
			return true
		}
		pos := basePos + int64(i)
		if !entry.videoFrameStartSet {
			entry.videoFrameStartSet = true
			entry.videoFrameStart = pos
			return true
		}
		frameBytes := pos - entry.videoFrameStart
		if frameBytes > 0 {
			entry.videoFrameBytes += uint64(frameBytes)
			entry.videoFrameBytesCount++
		}
		entry.videoFrameStart = pos
		return entry.videoFrameBytesCount < 2
	})
	if len(buf) >= 3 {
		entry.videoFrameCarry = append(entry.videoFrameCarry[:0], buf[len(buf)-3:]...)
	} else {
		entry.videoFrameCarry = append(entry.videoFrameCarry[:0], buf...)
	}
	entry.videoFramePos += int64(len(payload))
}
