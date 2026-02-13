package mediainfo

// Minimal LATM/LOAS (AAC over MPEG-TS, stream_type 0x11) parsing for MediaInfo parity.
// We only need StreamMuxConfig -> AudioSpecificConfig to expose LC, sample rate, and channel count.

// consumeLATM parses LOAS frames (syncword 0x2B7) from the payload and extracts AAC metadata.
func consumeLATM(entry *tsStream, payload []byte) {
	if len(payload) == 0 {
		return
	}
	entry.audioBuffer = append(entry.audioBuffer, payload...)

	i := 0
	for i+3 <= len(entry.audioBuffer) {
		// LOAS syncword (11 bits) at byte boundary: 0x56, then top 3 bits of next byte == 0b111.
		if entry.audioBuffer[i] != 0x56 || (entry.audioBuffer[i+1]&0xE0) != 0xE0 {
			i++
			continue
		}
		length := (int(entry.audioBuffer[i+1]&0x1F) << 8) | int(entry.audioBuffer[i+2])
		if length <= 0 {
			i++
			continue
		}
		total := 3 + length
		if i+total > len(entry.audioBuffer) {
			break
		}
		frame := entry.audioBuffer[i+3 : i+total]
		entry.audioFrames++
		if !entry.hasAudioInfo {
			if objType, sr, ch, ok := parseLATMAudioSpecificConfig(frame); ok && objType > 0 && sr > 0 && ch > 0 {
				entry.audioProfile = mapAACProfile(objType)
				entry.audioObject = objType
				entry.audioRate = sr
				entry.audioChannels = uint64(ch)
				// AAC-LC defaults to 1024 SPF (MediaInfo output).
				if entry.audioSpf == 0 {
					entry.audioSpf = 1024
				}
				entry.hasAudioInfo = true
			}
		}
		i += total
	}
	if i > 0 {
		entry.audioBuffer = append(entry.audioBuffer[:0], entry.audioBuffer[i:]...)
	}
}

func parseLATMAudioSpecificConfig(audioMuxElement []byte) (objType int, sampleRate float64, channels int, ok bool) {
	if len(audioMuxElement) == 0 {
		return 0, 0, 0, false
	}
	br := newBitReader(audioMuxElement)
	useSame := br.readBitsValue(1)
	if useSame == ^uint64(0) {
		return 0, 0, 0, false
	}
	// When useSameStreamMux==1, the StreamMuxConfig is omitted; we can't learn metadata from this element.
	if useSame == 1 {
		return 0, 0, 0, false
	}

	audioMuxVersion := br.readBitsValue(1)
	if audioMuxVersion != 0 {
		return 0, 0, 0, false
	}
	allStreamsSameTimeFraming := br.readBitsValue(1)
	if allStreamsSameTimeFraming != 1 {
		return 0, 0, 0, false
	}
	if br.readBitsValue(6) == ^uint64(0) { // numSubFrames
		return 0, 0, 0, false
	}
	numProgram := br.readBitsValue(4)
	if numProgram != 0 {
		return 0, 0, 0, false
	}
	numLayer := br.readBitsValue(3)
	if numLayer != 0 {
		return 0, 0, 0, false
	}

	objType, sampleRate, channels, ok = parseAACAudioSpecificConfigBits(br)
	if !ok {
		return 0, 0, 0, false
	}

	frameLengthType := br.readBitsValue(3)
	if frameLengthType == ^uint64(0) {
		return 0, 0, 0, false
	}
	// Skip common frameLengthType values (0 is most common).
	switch frameLengthType {
	case 0:
		_ = br.readBitsValue(8) // latmBufferFullness
	case 1:
		_ = br.readBitsValue(9)
	case 3, 4, 5:
		_ = br.readBitsValue(6)
	case 6, 7:
		_ = br.readBitsValue(1)
	default:
	}

	otherDataPresent := br.readBitsValue(1)
	if otherDataPresent == ^uint64(0) {
		return 0, 0, 0, false
	}
	if otherDataPresent == 1 {
		// audioMuxVersion==0: otherDataLenBits is encoded as a sum of 8-bit chunks.
		otherLen := 0
		for {
			v := br.readBitsValue(8)
			if v == ^uint64(0) {
				return 0, 0, 0, false
			}
			otherLen += int(v)
			if v != 255 {
				break
			}
		}
		// Skip other_data bits.
		if otherLen > 0 {
			bits := otherLen * 8
			for bits > 0 {
				n := bits
				if n > 255 {
					n = 255
				}
				if br.readBitsValue(uint8(n)) == ^uint64(0) {
					return 0, 0, 0, false
				}
				bits -= n
			}
		}
	}
	crcPresent := br.readBitsValue(1)
	if crcPresent == ^uint64(0) {
		return 0, 0, 0, false
	}
	if crcPresent == 1 {
		_ = br.readBitsValue(8)
	}

	return objType, sampleRate, channels, true
}

func parseAACAudioSpecificConfigBits(br *bitReader) (objType int, sampleRate float64, channels int, ok bool) {
	objType, ok = readAACAudioObjectType(br)
	if !ok || objType <= 0 {
		return 0, 0, 0, false
	}

	sfIndex := br.readBitsValue(4)
	if sfIndex == ^uint64(0) {
		return 0, 0, 0, false
	}
	if sfIndex == 0xF {
		sf := br.readBitsValue(24)
		if sf == ^uint64(0) {
			return 0, 0, 0, false
		}
		sampleRate = float64(sf)
	} else {
		sampleRate = adtsSampleRate(int(sfIndex))
	}

	chCfg := br.readBitsValue(4)
	if chCfg == ^uint64(0) {
		return 0, 0, 0, false
	}
	channels = int(chCfg)

	// Handle explicit SBR/PS signaling: update objType to the extension object type.
	// (We keep the base sample rate and channels for parity with most TS LATM streams.)
	if objType == 5 || objType == 29 {
		extIndex := br.readBitsValue(4)
		if extIndex == ^uint64(0) {
			return objType, sampleRate, channels, true
		}
		if extIndex == 0xF {
			_ = br.readBitsValue(24)
		}
		next, ok := readAACAudioObjectType(br)
		if ok && next > 0 {
			objType = next
		}
	}

	return objType, sampleRate, channels, true
}
