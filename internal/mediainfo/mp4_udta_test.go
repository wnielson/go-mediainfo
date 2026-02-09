package mediainfo

import (
	"bytes"
	"testing"
)

func TestParseMP4Description(t *testing.T) {
	want := "Packed by Bilibili XCoder v1.0"

	var ilst bytes.Buffer
	writeMP4Box(&ilst, "desc", makeMP4DataBox(want))

	var meta bytes.Buffer
	meta.Write(make([]byte, 4)) // version/flags
	writeMP4Box(&meta, "ilst", ilst.Bytes())

	var udta bytes.Buffer
	writeMP4Box(&udta, "meta", meta.Bytes())

	got := parseMP4Description(udta.Bytes())
	if got != want {
		t.Fatalf("desc=%q, want %q", got, want)
	}
}

func makeMP4DataBox(value string) []byte {
	var b bytes.Buffer
	b.Write(make([]byte, 8)) // type/locale/reserved; parser skips first 8 bytes
	b.WriteString(value)
	var out bytes.Buffer
	writeMP4Box(&out, "data", b.Bytes())
	return out.Bytes()
}
