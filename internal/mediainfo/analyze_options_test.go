package mediainfo

import "testing"

func TestDefaultAnalyzeOptions(t *testing.T) {
	opts := defaultAnalyzeOptions()
	if opts.ParseSpeed != 0.5 {
		t.Fatalf("ParseSpeed=%v, want 0.5", opts.ParseSpeed)
	}
	if opts.TestContinuousFileNames {
		t.Fatalf("TestContinuousFileNames=%v, want false", opts.TestContinuousFileNames)
	}
}

func TestNormalizeAnalyzeOptionsDefaults(t *testing.T) {
	opts := normalizeAnalyzeOptions(AnalyzeOptions{})
	if opts.ParseSpeed != 0.5 {
		t.Fatalf("ParseSpeed=%v, want 0.5", opts.ParseSpeed)
	}
	if opts.TestContinuousFileNames {
		t.Fatalf("TestContinuousFileNames=%v, want false", opts.TestContinuousFileNames)
	}
}
