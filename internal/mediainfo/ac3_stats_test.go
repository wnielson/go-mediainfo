package mediainfo

import "testing"

func TestAC3MergeFrame_DynrngStatsSkipUnset(t *testing.T) {
	var acc ac3Info

	// dynrnge=1 but dynrngCode=0xFF is treated as "unset" in MediaInfoLib stats.
	acc.mergeFrame(ac3Info{dynrngParsed: true, dynrnge: true, dynrngCode: 0xFF})
	if got := acc.dynrngs[0xFF]; got != 0 {
		t.Fatalf("dynrng[0xFF]=%d, want 0", got)
	}

	// dynrnge=0 still contributes to stats as code 0 (MediaInfoLib counts 0 when dynrnge is absent).
	acc.mergeFrame(ac3Info{dynrngParsed: true, dynrnge: false, dynrngCode: 0})
	if got := acc.dynrngs[0x00]; got != 1 {
		t.Fatalf("dynrng[0x00]=%d, want 1", got)
	}

	// A real dynrng code should be counted.
	acc.mergeFrame(ac3Info{dynrngParsed: true, dynrnge: true, dynrngCode: 0x10})
	if got := acc.dynrngs[0x10]; got != 1 {
		t.Fatalf("dynrng[0x10]=%d, want 1", got)
	}
}
