package mediainfo

import "testing"

func TestExtractLeadingNumber(t *testing.T) {
	cases := []struct {
		value string
		want  string
	}{
		{value: "1 920 pixels", want: "1920"},
		{value: "640", want: "640"},
		{value: "  29.970 FPS", want: "29.970"},
		{value: "", want: ""},
	}
	for _, tc := range cases {
		if got := extractLeadingNumber(tc.value); got != tc.want {
			t.Fatalf("extractLeadingNumber(%q)=%q want %q", tc.value, got, tc.want)
		}
	}
}
