package mediainfo

import "testing"

type ac3BitWriter struct {
	buf    []byte
	bitPos int
}

func (w *ac3BitWriter) writeBits(v uint32, n int) {
	for i := n - 1; i >= 0; i-- {
		bit := (v >> uint(i)) & 1
		bytePos := w.bitPos >> 3
		shift := 7 - (w.bitPos & 7)
		if bit == 1 {
			w.buf[bytePos] |= 1 << uint(shift)
		}
		w.bitPos++
	}
}

func TestParseAC3Frame_Mono_NoCmixlev(t *testing.T) {
	// Frame layout (subset) matching parseAC3Frame():
	// syncword, crc1, fscod, frmsizecod, bsid, bsmod, acmod=1 (mono),
	// lfeon, dialnorm, compre+compr, then enough padding for parseAC3Dynrng().
	const (
		fscod      = 0 // 48 kHz
		frmsizecod = 0 // 32 kbps @ 48 kHz -> 128 bytes
		bsid       = 6 // common in the wild; also exercises xbsi1/xbsi2 parsing branch
		bsmod      = 0
		acmod      = 1 // mono (C) -> must NOT carry cmixlev
		lfeon      = 0
		dialnorm   = 24 // -24 dB
		compre     = 1
		compr      = 0x20
	)

	frameSize := ac3FrameSizeBytes(fscod, frmsizecod)
	if frameSize == 0 {
		t.Fatalf("unexpected frameSize=0 for fscod=%d frmsizecod=%d", fscod, frmsizecod)
	}
	buf := make([]byte, frameSize)
	bw := ac3BitWriter{buf: buf}

	bw.writeBits(0x0B77, 16)    // syncword
	bw.writeBits(0x0000, 16)    // crc1
	bw.writeBits(fscod, 2)      // fscod
	bw.writeBits(frmsizecod, 6) // frmsizecod
	bw.writeBits(bsid, 5)       // bsid
	bw.writeBits(bsmod, 3)      // bsmod
	bw.writeBits(acmod, 3)      // acmod
	bw.writeBits(lfeon, 1)      // lfeon
	bw.writeBits(dialnorm, 5)   // dialnorm
	bw.writeBits(compre, 1)     // compre
	bw.writeBits(compr, 8)      // compr
	bw.writeBits(0, 1)          // langcode
	bw.writeBits(0, 1)          // audprodie
	bw.writeBits(0, 1)          // copyrightb
	bw.writeBits(0, 1)          // origbs
	bw.writeBits(0, 1)          // xbsi1e (bsid==6)
	bw.writeBits(0, 1)          // xbsi2e (bsid==6)
	bw.writeBits(0, 1)          // addbsie
	bw.writeBits(0, 1)          // (ac3 dynrng) skip bit 1/2 for nfchans=1
	bw.writeBits(0, 1)          // (ac3 dynrng) skip bit 2/2 for nfchans=1
	bw.writeBits(0, 1)          // dynrnge=0

	info, gotSize, ok := parseAC3Frame(buf)
	if !ok {
		t.Fatalf("parseAC3Frame returned ok=false")
	}
	if gotSize != frameSize {
		t.Fatalf("frameSize mismatch: got=%d want=%d", gotSize, frameSize)
	}
	if info.acmod != acmod {
		t.Fatalf("acmod mismatch: got=%d want=%d", info.acmod, acmod)
	}
	if info.lfeon != lfeon {
		t.Fatalf("lfeon mismatch: got=%d want=%d", info.lfeon, lfeon)
	}
	if info.dialnorm != -24 {
		t.Fatalf("dialnorm mismatch: got=%d want=-24", info.dialnorm)
	}
	if info.hasCmixlev {
		t.Fatalf("mono stream must not set cmixlev (bitstream alignment regression)")
	}
	if !info.compre || !info.hasComprField || !info.hasCompr || info.comprCode != compr {
		t.Fatalf("compr mismatch: compre=%v hasComprField=%v hasCompr=%v comprCode=0x%02x want=0x%02x", info.compre, info.hasComprField, info.hasCompr, info.comprCode, compr)
	}
}
