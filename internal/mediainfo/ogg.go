package mediainfo

import (
	"bytes"
	"encoding/binary"
	"io"
)

func ParseOgg(file io.ReadSeeker, size int64) (ContainerInfo, []Stream, bool) {
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return ContainerInfo{}, nil, false
	}

	var (
		sampleRate  uint32
		channels    uint8
		lastGranule uint64
		dataBytes   uint64
		format      string
	)
	var segTable [255]byte

	for {
		var header [27]byte
		if _, err := io.ReadFull(file, header[:]); err != nil {
			break
		}
		if header[0] != 'O' || header[1] != 'g' || header[2] != 'g' || header[3] != 'S' {
			return ContainerInfo{}, nil, false
		}
		granule := binary.LittleEndian.Uint64(header[6:14])
		segCount := int(header[26])
		if segCount > len(segTable) {
			return ContainerInfo{}, nil, false
		}
		if _, err := io.ReadFull(file, segTable[:segCount]); err != nil {
			return ContainerInfo{}, nil, false
		}
		dataLen := 0
		for _, seg := range segTable[:segCount] {
			dataLen += int(seg)
		}
		if dataLen > 0 {
			dataBytes += uint64(dataLen)
			if sampleRate == 0 {
				peek := min(dataLen, 64)
				data := make([]byte, peek)
				if _, err := io.ReadFull(file, data); err != nil {
					return ContainerInfo{}, nil, false
				}
				if sr, ch, fmt := parseOggIdentification(data); sr > 0 {
					sampleRate = sr
					channels = ch
					format = fmt
				}
				if dataLen > peek {
					if _, err := file.Seek(int64(dataLen-peek), io.SeekCurrent); err != nil {
						return ContainerInfo{}, nil, false
					}
				}
			} else {
				if _, err := file.Seek(int64(dataLen), io.SeekCurrent); err != nil {
					return ContainerInfo{}, nil, false
				}
			}
		}
		if granule != ^uint64(0) && granule > lastGranule {
			lastGranule = granule
		}
	}

	if sampleRate == 0 {
		return ContainerInfo{}, nil, false
	}

	duration := 0.0
	if lastGranule > 0 {
		duration = float64(lastGranule) / float64(sampleRate)
	}

	bitrate := 0.0
	if duration > 0 {
		bitrate = (float64(dataBytes) * 8) / duration
	}

	info := ContainerInfo{
		DurationSeconds: duration,
		BitrateMode:     "Variable",
	}

	if format == "" {
		format = "Unknown"
	}

	fields := []Field{
		{Name: "Format", Value: format},
	}
	fields = appendChannelFields(fields, uint64(channels))
	fields = appendSampleRateField(fields, float64(sampleRate))
	fields = append(fields, Field{Name: "Bit rate mode", Value: "Variable"})
	fields = addStreamCommon(fields, duration, bitrate)

	_ = size
	return info, []Stream{{Kind: StreamAudio, Fields: fields}}, true
}

func parseOggIdentification(data []byte) (uint32, uint8, string) {
	if len(data) < 16 {
		return 0, 0, ""
	}
	if data[0] == 0x01 && bytes.Equal(data[1:7], []byte("vorbis")) {
		channels := data[11]
		sampleRate := binary.LittleEndian.Uint32(data[12:16])
		return sampleRate, channels, "Vorbis"
	}
	if bytes.HasPrefix(data, []byte("OpusHead")) && len(data) >= 19 {
		channels := data[9]
		return 48000, channels, "Opus"
	}
	return 0, 0, ""
}
