package updater

import (
	"strconv"
	"strings"

	"golang.org/x/mod/semver"
)

// Version represents a parsed semantic version.
type Version struct {
	Major      int
	Minor      int
	Patch      int
	PreRelease string
	Raw        string
}

// Canonical ensures a version string starts with "v" and conforms to standard semver.
func Canonical(v string) string {
	s := strings.TrimSpace(v)
	if s == "" {
		return ""
	}
	if !strings.HasPrefix(s, "v") {
		s = "v" + s
	}

	// Normalize shorthand versions like v1 or v1.2
	parts := strings.Split(s, "-")
	base := parts[0]
	pre := ""
	if len(parts) > 1 {
		pre = "-" + strings.Join(parts[1:], "-")
	}

	dots := strings.Count(base, ".")
	if dots == 0 {
		base += ".0.0"
	} else if dots == 1 {
		base += ".0"
	}

	return semver.Canonical(base + pre)
}

// ParseVersion parses a version string.
func ParseVersion(v string) (Version, bool) {
	c := Canonical(v)
	if !semver.IsValid(c) {
		return Version{Raw: v}, false
	}

	majStr := strings.TrimPrefix(semver.Major(c), "v")
	maj, _ := strconv.Atoi(majStr)

	majMin := semver.MajorMinor(c)
	minStr := strings.TrimPrefix(majMin, semver.Major(c)+".")
	min, _ := strconv.Atoi(minStr)

	// Extract patch
	withoutPre := strings.Split(strings.TrimPrefix(c, "v"), "-")[0]
	patchParts := strings.Split(withoutPre, ".")
	patch := 0
	if len(patchParts) >= 3 {
		patch, _ = strconv.Atoi(patchParts[2])
	}

	pre := strings.TrimPrefix(semver.Prerelease(c), "-")

	return Version{
		Major:      maj,
		Minor:      min,
		Patch:      patch,
		PreRelease: pre,
		Raw:        v,
	}, true
}

// Compare compares two version strings using official Go module semver semantics.
// Returns 1 if v1 > v2, -1 if v1 < v2, and 0 if v1 == v2.
func Compare(v1Str, v2Str string) int {
	c1 := Canonical(v1Str)
	c2 := Canonical(v2Str)

	if !semver.IsValid(c1) && !semver.IsValid(c2) {
		return strings.Compare(v1Str, v2Str)
	}
	if !semver.IsValid(c1) {
		return -1
	}
	if !semver.IsValid(c2) {
		return 1
	}

	return semver.Compare(c1, c2)
}

// IsNewer returns true if targetVersion is strictly newer than currentVersion.
func IsNewer(targetVersion, currentVersion string) bool {
	return Compare(targetVersion, currentVersion) > 0
}
