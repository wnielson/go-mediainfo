package mediainfo

import "testing"

type bitWriter struct {
	b   []byte
	bit int
}

func (w *bitWriter) writeBits(v uint32, n int) {
	for i := n - 1; i >= 0; i-- {
		bit := (v >> i) & 1
		byteIdx := w.bit >> 3
		bitIdx := 7 - (w.bit & 7)
		if bit != 0 {
			w.b[byteIdx] |= 1 << bitIdx
		}
		w.bit++
	}
}

func makeEAC3Frame(t *testing.T, frameSize int, dialnormCode uint32) []byte {
	t.Helper()
	if frameSize < 8 {
		t.Fatalf("frameSize too small: %d", frameSize)
	}
	frmsiz := uint32(frameSize/2 - 1)
	if frmsiz > 0x7FF {
		t.Fatalf("frameSize too large: %d", frameSize)
	}

	b := make([]byte, frameSize)
	w := bitWriter{b: b}
	w.writeBits(0x0B77, 16) // syncword
	w.writeBits(0, 2)       // strmtyp: independent
	w.writeBits(0, 3)       // substreamid
	w.writeBits(frmsiz, 11)
	w.writeBits(0, 2)  // fscod: 48kHz
	w.writeBits(3, 2)  // numblkscod: 6 blocks
	w.writeBits(2, 3)  // acmod: 2/0 (stereo)
	w.writeBits(0, 1)  // lfeon
	w.writeBits(16, 5) // bsid (>=10 for E-AC-3 sanity)
	w.writeBits(dialnormCode&0x1F, 5)
	w.writeBits(0, 1) // compre
	return b
}

func TestProbeMatroskaAudio_EAC3MultiFrameNonLacedPacket(t *testing.T) {
	const trackID = uint64(1)
	const frameSize = 16

	f1 := makeEAC3Frame(t, frameSize, 1)
	f2 := makeEAC3Frame(t, frameSize, 2)
	payload := append(append([]byte{}, f1...), f2...)

	// Non-laced packets may contain multiple E-AC-3 syncframes back-to-back; ensure the
	// probe path aggregates them when packetAligned=false.
	probes := map[uint64]*matroskaAudioProbe{
		trackID: {format: "E-AC-3", collect: true},
	}
	probeMatroskaAudio(probes, trackID, payload, 1, int64(len(payload)), false)
	p := probes[trackID]
	if !p.ok {
		t.Fatalf("expected probe ok")
	}
	if got := p.info.dialnormCount; got != 2 {
		t.Fatalf("expected dialnormCount=2, got %d", got)
	}
	if p.info.dialnormMin != -2 || p.info.dialnormMax != -1 {
		t.Fatalf("expected dialnormMin=-2 dialnormMax=-1, got min=%d max=%d", p.info.dialnormMin, p.info.dialnormMax)
	}

	// With packetAligned=true, probe expects exactly one frame per packet and rejects mismatched sizes.
	probes2 := map[uint64]*matroskaAudioProbe{
		trackID: {format: "E-AC-3", collect: true},
	}
	probeMatroskaAudio(probes2, trackID, payload, 1, int64(len(payload)), true)
	p2 := probes2[trackID]
	if p2.ok {
		t.Fatalf("expected probe not ok with packetAligned=true")
	}
	if got := p2.info.dialnormCount; got != 0 {
		t.Fatalf("expected dialnormCount=0, got %d", got)
	}
}

func TestApplyMatroskaStats_AudioDurationAlsoSetsJSON(t *testing.T) {
	info := MatroskaInfo{
		Tracks: []Stream{
			{
				Kind: StreamAudio,
				Fields: []Field{
					{Name: "ID", Value: "1"},
					{Name: "Format", Value: "AAC"},
				},
			},
		},
	}

	stats := map[uint64]*matroskaTrackStats{
		1: {
			hasTime:   true,
			minTimeNs: 0,
			maxTimeNs: int64(4.321 * 1e9),
		},
	}

	applyMatroskaStats(&info, stats, 0)

	if got := findField(info.Tracks[0].Fields, "Duration"); got == "" {
		t.Fatalf("expected Duration field set")
	}
	if info.Tracks[0].JSON == nil || info.Tracks[0].JSON["Duration"] == "" {
		t.Fatalf("expected JSON Duration set")
	}
}
