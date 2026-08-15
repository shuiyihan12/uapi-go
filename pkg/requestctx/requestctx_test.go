package requestctx

import "testing"

// TestNormalizeRegion pins the region normalization contract: the three
// Production regions are matched case-insensitively and
// whitespace-tolerantly; anything else — including typos — normalizes to the
// empty string, which makes pkg/client fall back to the default endpoint
// (with a warning log since the invalid-region fallback fix).
func TestNormalizeRegion(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"", ""},
		{"   ", ""},
		{"apac", "apac"},
		{"AMERICAS", "americas"},
		{"  Emea ", "emea"},
		{"APA", ""},    // typo of apac
		{"europe", ""}, // not a Production region
		{"1G", ""},     // provider code, not a region
	}
	for _, tt := range tests {
		if got := NormalizeRegion(tt.in); got != tt.want {
			t.Errorf("NormalizeRegion(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
