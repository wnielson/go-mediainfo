package mediainfo

import (
	"os"
	"path/filepath"
	"testing"
)

func writeSizedFile(t *testing.T, path string, size int) {
	t.Helper()
	data := make([]byte, size)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write file %s: %v", path, err)
	}
}

func TestDetectContinuousFileSetIncludesSparseNumericFiles(t *testing.T) {
	dir := t.TempDir()
	start := filepath.Join(dir, "00000.m2ts")
	writeSizedFile(t, start, 3)
	writeSizedFile(t, filepath.Join(dir, "00001.m2ts"), 5)
	writeSizedFile(t, filepath.Join(dir, "00003.m2ts"), 7)

	set, ok := detectContinuousFileSet(start)
	if !ok {
		t.Fatalf("detectContinuousFileSet() ok=false, want true")
	}
	if set.Count != 3 {
		t.Fatalf("Count=%d, want 3", set.Count)
	}
	wantLast := filepath.Join(dir, "00003.m2ts")
	if set.LastPath != wantLast {
		t.Fatalf("LastPath=%q, want %q", set.LastPath, wantLast)
	}
	if set.LastSize != 7 {
		t.Fatalf("LastSize=%d, want 7", set.LastSize)
	}
	if set.TotalSize != 15 {
		t.Fatalf("TotalSize=%d, want 15", set.TotalSize)
	}
}

func TestDetectContinuousFileSetRequiresFollowingFile(t *testing.T) {
	dir := t.TempDir()
	start := filepath.Join(dir, "00000.m2ts")
	writeSizedFile(t, start, 3)

	if _, ok := detectContinuousFileSet(start); ok {
		t.Fatalf("detectContinuousFileSet() ok=true, want false")
	}
}
