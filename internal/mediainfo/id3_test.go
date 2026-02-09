package mediainfo

import "testing"

func utf16LEWithBOM(s string) []byte {
	out := []byte{0xFF, 0xFE}
	for i := 0; i < len(s); i++ {
		out = append(out, s[i], 0x00)
	}
	return out
}

func TestID3COMMNormalizesNewlinesLikeMediaInfo(t *testing.T) {
	// Encoding: UTF-8, language: "eng", desc: "", comment: lines with blank line.
	data := []byte{0x03, 'e', 'n', 'g', 0x00}
	data = append(data, []byte("a\r\n\r\nb")...)
	got, ok := parseID3COMM(data)
	if !ok {
		t.Fatalf("parseID3COMM ok=false")
	}
	want := "a /  / b"
	if got != want {
		t.Fatalf("parseID3COMM=%q want %q", got, want)
	}
}

func TestID3WXXXAllowsEmptyDescription(t *testing.T) {
	// Encoding: UTF-8, desc: "", url: "http://x".
	data := []byte{0x03, 0x00}
	data = append(data, []byte("http://x")...)
	desc, url, ok := parseID3WXXX(data)
	if !ok {
		t.Fatalf("parseID3WXXX ok=false")
	}
	if desc != "URL" {
		t.Fatalf("desc=%q want %q", desc, "URL")
	}
	if url != "http://x" {
		t.Fatalf("url=%q want %q", url, "http://x")
	}
}

func TestID3TXXXUTF16TerminatorAlignmentDoesNotDropLastChar(t *testing.T) {
	// Regression: naive bytes.Index(0x00,0x00) can match across a UTF-16LE code-unit boundary and truncate.
	descBytes := utf16LEWithBOM("major_brand")
	valBytes := utf16LEWithBOM("isom")
	data := []byte{0x01}              // UTF-16 with BOM
	data = append(data, descBytes...) // description
	data = append(data, 0x00, 0x00)   // terminator
	data = append(data, valBytes...)  // value
	desc, value, ok := parseID3TXXX(data)
	if !ok {
		t.Fatalf("parseID3TXXX ok=false")
	}
	if desc != "major_brand" {
		t.Fatalf("desc=%q want %q", desc, "major_brand")
	}
	if value != "isom" {
		t.Fatalf("value=%q want %q", value, "isom")
	}
}
