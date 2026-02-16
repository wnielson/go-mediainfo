package mediainfo

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestM2TSSmoke(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.m2ts")

	// Minimal BDAV packet structure: 192-byte packets with TS sync byte at offset 4.
	packet := make([]byte, 192)
	packet[4] = 0x47
	data := make([]byte, 0, 192*4)
	for i := 0; i < 4; i++ {
		data = append(data, packet...)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write m2ts sample: %v", err)
	}

	report, err := AnalyzeFile(path)
	if err != nil {
		t.Fatalf("analyze m2ts sample: %v", err)
	}

	if out := RenderText([]Report{report}); out == "" {
		t.Fatalf("empty text output")
	}
	jsonOut := RenderJSON([]Report{report})
	var root any
	if err := json.Unmarshal([]byte(jsonOut), &root); err != nil {
		t.Fatalf("parse json output: %v", err)
	}
}
