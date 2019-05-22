package utility

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidSemVerString = errors.New("provided string is not a valid SemVer")
var ErrInvalidSemVer = errors.New("provided sem ver version invalid")

// SemVer represents a semantic version value of the form major.minor.patch
// See https://semver.org/
type SemVer struct {
	major int
	minor int
	patch int
}

var InvalidSemVer = &SemVer{major: 0, minor: 0, patch: 0}
var InvalidSemVerString = ""

// NewSemVer creates a new SemVer using the provided major, minor, and patch values.
//
//
// Parameters
//
// major: major version, indicates breaking changes to API when increased
//
// minor: minor version, indicates feature added when increased
//
// patch: patch version, indicates backwards compatible changes to API
//
//
// Outputs
//
// ver: a pointer to a new SemVer object
//
// err: if nil no error, otherwise ErrInvalidSemVer if any of the values are less than 0
func NewSemVer(major, minor, patch int) (ver *SemVer, err error) {
	if major < 0 || minor < 0 || patch < 0 {
		return nil, ErrInvalidSemVer
	}
	return &SemVer{major: major, minor: minor, patch: patch}, nil
}

// SemVerFromString creates a new SemVer using the provided string.
//
// Parameters
// version: a sem ver version string, {major | major.minor | major.minor.patch}
//
// Outputs
// ver: a pointer to a new SemVer object
// err: if nil no error, otherwise ErrInvalidSemVerString if the string is not a valid SemVer
func SemVerFromString(version string) (ver *SemVer, err error) {
	semVer := &SemVer{}
	dots := strings.Count(version, ".")
	if dots > 2 {
		return InvalidSemVer, ErrInvalidSemVerString
	} else {
		mmp := strings.Split(version, ".")
		for index, e := range mmp {
			verPart, err := strconv.Atoi(e)
			if err != nil {
				return InvalidSemVer, ErrInvalidSemVerString
			}
			if verPart < 0 {
				return InvalidSemVer, ErrInvalidSemVerString
			}
			switch index {
			case 0:
				semVer.major = verPart
				break
			case 1:
				semVer.minor = verPart
				break
			case 2:
				semVer.patch = verPart
			default:
				return InvalidSemVer, ErrInvalidSemVerString
			}
		}
		return semVer, nil
	}
}

func (sv *SemVer) String() string {
	if sv == nil {
		return InvalidSemVerString
	}
	return fmt.Sprintf("%d.%d.%d", sv.major, sv.minor, sv.patch)
}

func (sv *SemVer) StringMajor() string {
	if sv == nil {
		return InvalidSemVerString
	}
	return fmt.Sprintf("%d", sv.major)
}

func (sv *SemVer) StringMajorMinor() string {
	if sv == nil {
		return InvalidSemVerString
	}
	return fmt.Sprintf("%d.%d", sv.major, sv.minor)
}

func (sv *SemVer) GetMajor() int {
	if sv == nil {
		return 0
	}
	return sv.major
}

func (sv *SemVer) GetMinor() int {
	if sv == nil {
		return 0
	}
	return sv.minor
}

func (sv *SemVer) GetPatch() int {
	if sv == nil {
		return 0
	}
	return sv.patch
}

// Compare compares the SemVer with the provided otherSemVer.
//
// Parameters
// otherSemVer: a sem ver to compare the SemVer object against
//
// Outputs
// comp: compare value
//       -2: nil semver
//       -1: SemVer < otherSemVer
//        0: SemVer == otherSemVer
//        1: SemVer > otherSemVer
func (sv *SemVer) Compare(otherSemVer *SemVer) (comp int) {
	if otherSemVer == nil || sv == nil {
		return -2
	}
	if sv.major > otherSemVer.major {
		return 1
	} else if sv.major < otherSemVer.major {
		return -1
	} else {
		if sv.minor > otherSemVer.minor {
			return 1
		} else if sv.minor < otherSemVer.minor {
			return -1
		} else {
			if sv.patch > otherSemVer.patch {
				return 1
			} else if sv.patch < otherSemVer.patch {
				return -1
			} else {
				return 0
			}
		}
	}
}

func (sv *SemVer) Equals(otherSemVer *SemVer) bool {
	if otherSemVer == nil || sv == nil {
		return false
	}
	return sv.Compare(otherSemVer) == 0
}
