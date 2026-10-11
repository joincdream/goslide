package updater

import (
	"testing"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input       string
		wantOK      bool
		wantMajor   int
		wantMinor   int
		wantPatch   int
		wantPre     string
	}{
		{"v1.0.0", true, 1, 0, 0, ""},
		{"v0.2.1-dev", true, 0, 2, 1, "dev"},
		{"1.2.3", true, 1, 2, 3, ""},
		{"0.1", true, 0, 1, 0, ""},
		{"2", true, 2, 0, 0, ""},
		{"invalid", false, 0, 0, 0, ""},
		{"v1.x.3", false, 0, 0, 0, ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			v, ok := ParseVersion(tt.input)
			if ok != tt.wantOK {
				t.Fatalf("ParseVersion(%q) ok = %v, want %v", tt.input, ok, tt.wantOK)
			}
			if ok {
				if v.Major != tt.wantMajor || v.Minor != tt.wantMinor || v.Patch != tt.wantPatch || v.PreRelease != tt.wantPre {
					t.Errorf("ParseVersion(%q) = %+v, want %d.%d.%d-%s", tt.input, v, tt.wantMajor, tt.wantMinor, tt.wantPatch, tt.wantPre)
				}
			}
		})
	}
}

func TestCompare(t *testing.T) {
	tests := []struct {
		v1   string
		v2   string
		want int
	}{
		{"v1.0.0", "v1.0.0", 0},
		{"v1.0.1", "v1.0.0", 1},
		{"v1.0.0", "v1.0.1", -1},
		{"v1.1.0", "v1.0.9", 1},
		{"v2.0.0", "v1.9.9", 1},
		{"v1.0.0", "v1.0.0-dev", 1},
		{"v1.0.0-dev", "v1.0.0", -1},
		{"v0.1.0", "v0.1.0-dev", 1},
		{"v0.2.0", "v0.1.0-dev", 1},
		{"v0.1.0-alpha", "v0.1.0-beta", -1},
	}

	for _, tt := range tests {
		t.Run(tt.v1+"_vs_"+tt.v2, func(t *testing.T) {
			got := Compare(tt.v1, tt.v2)
			if got != tt.want {
				t.Errorf("Compare(%q, %q) = %d, want %d", tt.v1, tt.v2, got, tt.want)
			}
		})
	}
}

func TestIsNewer(t *testing.T) {
	if !IsNewer("v0.2.0", "v0.1.0") {
		t.Errorf("expected v0.2.0 to be newer than v0.1.0")
	}
	if !IsNewer("v0.1.0", "v0.1.0-dev") {
		t.Errorf("expected v0.1.0 to be newer than v0.1.0-dev")
	}
	if IsNewer("v0.1.0", "v0.2.0") {
		t.Errorf("expected v0.1.0 NOT to be newer than v0.2.0")
	}
	if IsNewer("v0.1.0", "v0.1.0") {
		t.Errorf("expected v0.1.0 NOT to be newer than v0.1.0")
	}
}
