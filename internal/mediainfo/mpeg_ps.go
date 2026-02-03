package mediainfo

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"sort"
)

func ParseMPEGPS(file io.ReadSeeker, size int64) (ContainerInfo, []Stream, bool) {
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return ContainerInfo{}, nil, false
	}

	reader := bufio.NewReaderSize(file, 1<<20)
	data, err := io.ReadAll(reader)
	if err != nil || len(data) < 4 {
		return ContainerInfo{}, nil, false
	}

	streams := map[byte]*psStream{}
	streamOrder := []byte{}
	videoParsers := map[byte]*mpeg2VideoParser{}
	var videoPTS ptsTracker
	var anyPTS ptsTracker

	for i := 0; i+4 <= len(data); {
		if data[i] != 0x00 || data[i+1] != 0x00 || data[i+2] != 0x01 {
			i++
			continue
		}
		streamID := data[i+3]
		kind, format := mapPSStream(streamID)
		if kind != "" {
			entry, exists := streams[streamID]
			if !exists {
				entry = &psStream{id: streamID, kind: kind, format: format}
				streams[streamID] = entry
				streamOrder = append(streamOrder, streamID)
			}
			entry.kind = kind
			entry.format = format
		} else {
			i += 4
			continue
		}

		if i+9 >= len(data) {
			i += 4
			continue
		}
		pesLen := int(binary.BigEndian.Uint16(data[i+4 : i+6]))
		flags := data[i+7]
		headerLen := int(data[i+8])
		payloadStart := i + 9 + headerLen
		if payloadStart > len(data) {
			i += 4
			continue
		}
		if (flags&0x80) != 0 && i+9+headerLen <= len(data) {
			if pts, ok := parsePTS(data[i+9:]); ok {
				anyPTS.add(pts)
				if entry := streams[streamID]; entry != nil {
					entry.pts.add(pts)
					if entry.kind == StreamVideo {
						videoPTS.add(pts)
					}
				}
			}
		}

		payloadLen := 0
		if pesLen > 0 {
			payloadLen = pesLen - 3 - headerLen
			if payloadLen < 0 {
				payloadLen = 0
			}
			if payloadStart+payloadLen > len(data) {
				payloadLen = len(data) - payloadStart
			}
		} else {
			next := nextStartCode(data, payloadStart)
			payloadLen = next - payloadStart
		}

		if payloadLen > 0 {
			if entry := streams[streamID]; entry != nil {
				entry.bytes += uint64(payloadLen)
				if entry.kind == StreamVideo {
					parser := videoParsers[streamID]
					if parser == nil {
						parser = &mpeg2VideoParser{}
						videoParsers[streamID] = parser
					}
					parser.consume(data[payloadStart : payloadStart+payloadLen])
				}
			}
		}

		i = payloadStart + payloadLen
		if i <= payloadStart {
			i = payloadStart + 1
		}
	}

	var streamsOut []Stream
	sort.Slice(streamOrder, func(i, j int) bool { return streamOrder[i] < streamOrder[j] })
	var videoFrameRate float64
	for _, id := range streamOrder {
		if st := streams[id]; st != nil && st.kind == StreamVideo {
			if parser := videoParsers[id]; parser != nil {
				info := parser.finalize()
				videoFrameRate = info.FrameRate
			}
			break
		}
	}
	for _, id := range streamOrder {
		st := streams[id]
		if st == nil {
			continue
		}
		fields := []Field{{Name: "ID", Value: formatID(uint64(st.id))}}
		if st.format != "" {
			fields = append(fields, Field{Name: "Format", Value: st.format})
		}
		if st.kind == StreamVideo {
			info := mpeg2VideoInfo{}
			if parser := videoParsers[id]; parser != nil {
				info = parser.finalize()
			}
			if info.Version != "" {
				fields = append(fields, Field{Name: "Format version", Value: info.Version})
			}
			if info.Profile != "" {
				fields = append(fields, Field{Name: "Format profile", Value: info.Profile})
			}
			if info.BVOP != nil {
				fields = append(fields, Field{Name: "Format settings, BVOP", Value: formatYesNo(*info.BVOP)})
			}
			if info.Matrix != "" {
				fields = append(fields, Field{Name: "Format settings, Matrix", Value: info.Matrix})
			}
			if info.GOPLength > 0 {
				fields = append(fields, Field{Name: "Format settings, GOP", Value: fmt.Sprintf("N=%d", info.GOPLength)})
			}
			duration := videoPTS.duration()
			if duration > 0 && info.FrameRate > 0 {
				duration += 2.0 / info.FrameRate
			}
			if duration > 0 {
				fields = addStreamDuration(fields, duration)
				fields = append(fields, Field{Name: "Bit rate mode", Value: "Variable"})
				bitrate := 0.0
				kbps := int64(0)
				if st.bytes > 0 {
					bitrate = (float64(st.bytes) * 8) / duration
					kbps = int64(bitrate / 1000.0)
					if kbps < 0 {
						kbps = 0
					}
					if value := formatBitrateKbps(kbps); value != "" {
						fields = appendFieldUnique(fields, Field{Name: "Bit rate", Value: value})
					}
				}
				if info.Width > 0 {
					fields = append(fields, Field{Name: "Width", Value: formatPixels(info.Width)})
				}
				if info.Height > 0 {
					fields = append(fields, Field{Name: "Height", Value: formatPixels(info.Height)})
				}
				if info.AspectRatio != "" {
					fields = append(fields, Field{Name: "Display aspect ratio", Value: info.AspectRatio})
				}
				if info.FrameRateNumer > 0 && info.FrameRateDenom > 0 {
					fields = append(fields, Field{Name: "Frame rate", Value: formatFrameRateRatio(info.FrameRateNumer, info.FrameRateDenom)})
				} else if info.FrameRate > 0 {
					fields = append(fields, Field{Name: "Frame rate", Value: formatFrameRate(info.FrameRate)})
				}
				if info.ColorSpace != "" {
					fields = append(fields, Field{Name: "Color space", Value: info.ColorSpace})
				}
				if info.ChromaSubsampling != "" {
					fields = append(fields, Field{Name: "Chroma subsampling", Value: info.ChromaSubsampling})
				}
				if info.BitDepth != "" {
					fields = append(fields, Field{Name: "Bit depth", Value: info.BitDepth})
				}
				if info.ScanType != "" {
					fields = append(fields, Field{Name: "Scan type", Value: info.ScanType})
				}
				fields = append(fields, Field{Name: "Compression mode", Value: "Lossy"})
				if bitrate > 0 && info.Width > 0 && info.Height > 0 {
					if bits := formatBitsPerPixelFrame(bitrate, info.Width, info.Height, info.FrameRate); bits != "" {
						fields = append(fields, Field{Name: "Bits/(Pixel*Frame)", Value: bits})
					}
				}
				if info.TimeCode != "" {
					fields = append(fields, Field{Name: "Time code of first frame", Value: info.TimeCode})
				}
				if info.TimeCodeSource != "" {
					fields = append(fields, Field{Name: "Time code source", Value: info.TimeCodeSource})
				}
				if info.GOPOpenClosed != "" {
					fields = append(fields, Field{Name: "GOP, Open/Closed", Value: info.GOPOpenClosed})
				}
				if info.GOPFirstClosed != "" {
					fields = append(fields, Field{Name: "GOP, Open/Closed of first frame", Value: info.GOPFirstClosed})
				}
				if kbps > 0 && duration > 0 {
					streamSizeBytes := int64(float64(kbps*1000)*duration/8.0 + 0.5)
					if streamSize := formatStreamSize(streamSizeBytes, size); streamSize != "" {
						fields = append(fields, Field{Name: "Stream size", Value: streamSize})
					}
				} else if st.bytes > 0 {
					if streamSize := formatStreamSize(int64(st.bytes), size); streamSize != "" {
						fields = append(fields, Field{Name: "Stream size", Value: streamSize})
					}
				}
			}
		} else if st.kind != StreamVideo {
			if duration := st.pts.duration(); duration > 0 {
				fields = addStreamDuration(fields, duration)
			}
		}
		streamsOut = append(streamsOut, Stream{Kind: st.kind, Fields: fields})
	}

	info := ContainerInfo{}
	if duration := videoPTS.duration(); duration > 0 {
		if videoFrameRate > 0 {
			duration += 2.0 / videoFrameRate
		}
		info.DurationSeconds = duration
	} else if duration := anyPTS.duration(); duration > 0 {
		info.DurationSeconds = duration
	}

	return info, streamsOut, true
}

func nextStartCode(data []byte, start int) int {
	for i := start; i+3 < len(data); i++ {
		if data[i] == 0x00 && data[i+1] == 0x00 && data[i+2] == 0x01 {
			return i
		}
	}
	return len(data)
}

func mapPSStream(streamID byte) (StreamKind, string) {
	switch {
	case streamID >= 0xE0 && streamID <= 0xEF:
		return StreamVideo, "MPEG Video"
	case streamID >= 0xC0 && streamID <= 0xDF:
		return StreamAudio, "MPEG Audio"
	case streamID == 0xBD:
		return StreamText, "Private"
	default:
		return "", ""
	}
}
