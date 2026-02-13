package mediainfo

import (
	"encoding/binary"
	"strconv"
)

func consumeDVBSubtitle(entry *tsStream, payload []byte) {
	if entry == nil || len(payload) < 2 {
		return
	}
	// data_identifier (0x20) + subtitle_stream_id.
	if payload[0] != 0x20 {
		return
	}

	if len(entry.dvbSubPageIDs) == 0 && len(entry.dvbSubRegionIDs) == 0 {
		entry.dvbSubStreamID = payload[1]
	}

	const maxRegions = 16
	const maxAddrs = 8

	pos := 2
	for pos+6 <= len(payload) {
		// subtitle_sync_byte
		if payload[pos] != 0x0F {
			pos++
			continue
		}
		segType := payload[pos+1]
		pageID := binary.BigEndian.Uint16(payload[pos+2 : pos+4])
		segLen := int(binary.BigEndian.Uint16(payload[pos+4 : pos+6]))
		pos += 6
		if segLen < 0 || pos+segLen > len(payload) {
			break
		}
		seg := payload[pos : pos+segLen]
		pos += segLen

		switch segType {
		case 0x10: // page composition segment
			if len(entry.dvbSubPageIDs) == 0 {
				entry.dvbSubPageIDs = append(entry.dvbSubPageIDs, pageID)
			}
			if len(seg) < 2 {
				continue
			}
			// region loop: region_id (8) + reserved (8) + region_horizontal_address (16) + region_vertical_address (16)
			for i := 2; i+6 <= len(seg) && len(entry.dvbSubRegionX) < maxAddrs && len(entry.dvbSubRegionY) < maxAddrs; i += 6 {
				x := binary.BigEndian.Uint16(seg[i+2:i+4]) & 0x0FFF
				y := binary.BigEndian.Uint16(seg[i+4:i+6]) & 0x0FFF
				entry.dvbSubRegionX = append(entry.dvbSubRegionX, x)
				entry.dvbSubRegionY = append(entry.dvbSubRegionY, y)
			}
		case 0x11: // region composition segment
			if len(seg) < 8 || len(entry.dvbSubRegionIDs) >= maxRegions {
				continue
			}
			regionID := seg[0]
			w := binary.BigEndian.Uint16(seg[2:4])
			h := binary.BigEndian.Uint16(seg[4:6])
			depthCode := (seg[6] >> 2) & 0x07
			depth := byte(0)
			switch depthCode {
			case 1:
				depth = 2
			case 2:
				depth = 4
			case 3:
				depth = 8
			}
			entry.dvbSubRegionIDs = append(entry.dvbSubRegionIDs, regionID)
			entry.dvbSubRegionW = append(entry.dvbSubRegionW, w)
			entry.dvbSubRegionH = append(entry.dvbSubRegionH, h)
			entry.dvbSubRegionDepth = append(entry.dvbSubRegionDepth, depth)
		default:
			// Ignore other segment types (CLUT/object/display definition).
		}

		if len(entry.dvbSubRegionIDs) >= 4 && len(entry.dvbSubRegionX) >= 2 && len(entry.dvbSubRegionW) >= 4 {
			// Enough metadata to match MediaInfo's typical first-scan output.
			return
		}
	}
}

func formatByteList(vals []byte) string {
	if len(vals) == 0 {
		return ""
	}
	out := make([]byte, 0, len(vals)*4)
	for i, v := range vals {
		if i > 0 {
			out = append(out, ' ', '/', ' ')
		}
		out = append(out, []byte(strconv.FormatUint(uint64(v), 10))...)
	}
	return string(out)
}

func formatUint16List(vals []uint16) string {
	if len(vals) == 0 {
		return ""
	}
	out := make([]byte, 0, len(vals)*6)
	for i, v := range vals {
		if i > 0 {
			out = append(out, ' ', '/', ' ')
		}
		out = append(out, []byte(strconv.FormatUint(uint64(v), 10))...)
	}
	return string(out)
}

func formatRepeatByte(v byte, n int) string {
	if n <= 0 {
		return ""
	}
	out := make([]byte, 0, n*4)
	for i := 0; i < n; i++ {
		if i > 0 {
			out = append(out, ' ', '/', ' ')
		}
		out = append(out, []byte(strconv.FormatUint(uint64(v), 10))...)
	}
	return string(out)
}

func formatRepeatUint16(v uint16, n int) string {
	if n <= 0 {
		return ""
	}
	out := make([]byte, 0, n*6)
	for i := 0; i < n; i++ {
		if i > 0 {
			out = append(out, ' ', '/', ' ')
		}
		out = append(out, []byte(strconv.FormatUint(uint64(v), 10))...)
	}
	return string(out)
}
