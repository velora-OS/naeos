package version

import (
	"fmt"
	"strconv"
	"strings"
)

type VersionInfo struct {
	NEIRVersion    string
	SchemaVersion  string
	ProjectVersion string
}

type SemVer struct {
	Major int
	Minor int
	Patch int
}

func Default() VersionInfo {
	return VersionInfo{NEIRVersion: "0.1.0", SchemaVersion: "1.0", ProjectVersion: "0.1.0"}
}

func ParseSemVer(s string) (SemVer, error) {
	s = strings.TrimPrefix(s, "v")
	parts := strings.Split(s, ".")
	if len(parts) != 3 {
		return SemVer{}, fmt.Errorf("invalid semver format: %s", s)
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return SemVer{}, fmt.Errorf("invalid major version: %s", parts[0])
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return SemVer{}, fmt.Errorf("invalid minor version: %s", parts[1])
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return SemVer{}, fmt.Errorf("invalid patch version: %s", parts[2])
	}
	if major < 0 || minor < 0 || patch < 0 {
		return SemVer{}, fmt.Errorf("version components must be non-negative")
	}
	return SemVer{Major: major, Minor: minor, Patch: patch}, nil
}

func (v SemVer) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func Compare(a, b SemVer) int {
	if a.Major != b.Major {
		if a.Major < b.Major {
			return -1
		}
		return 1
	}
	if a.Minor != b.Minor {
		if a.Minor < b.Minor {
			return -1
		}
		return 1
	}
	if a.Patch != b.Patch {
		if a.Patch < b.Patch {
			return -1
		}
		return 1
	}
	return 0
}

func IsCompatible(required, actual SemVer) bool {
	if required.Major == 0 {
		return required.Major == actual.Major && required.Minor == actual.Minor
	}
	return required.Major == actual.Major && actual.Minor >= required.Minor
}

func (vi VersionInfo) Validate() error {
	if vi.NEIRVersion == "" {
		return fmt.Errorf("NEIRVersion must not be empty")
	}
	if _, err := ParseSemVer(vi.NEIRVersion); err != nil {
		return fmt.Errorf("invalid NEIRVersion: %w", err)
	}
	if vi.SchemaVersion == "" {
		return fmt.Errorf("SchemaVersion must not be empty")
	}
	if _, err := ParseSemVer(vi.SchemaVersion); err != nil {
		return fmt.Errorf("invalid SchemaVersion: %w", err)
	}
	return nil
}

func (vi VersionInfo) IsCompatibleWith(other VersionInfo) bool {
	neirA, errA := ParseSemVer(vi.NEIRVersion)
	neirB, errB := ParseSemVer(other.NEIRVersion)
	if errA != nil || errB != nil {
		return false
	}
	return IsCompatible(neirA, neirB)
}
