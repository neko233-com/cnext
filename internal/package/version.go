package pkg

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
)

type VersionConstraint struct {
	Min   *semver.Version
	Max   *semver.Version
	Exact *semver.Version
	Range string
}

func ParseConstraint(v string) (*VersionConstraint, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil, fmt.Errorf("empty version constraint")
	}

	c := &VersionConstraint{Range: v}

	if strings.HasPrefix(v, ">=") {
		ver, err := semver.NewVersion(v[2:])
		if err != nil {
			return nil, fmt.Errorf("invalid version %q: %w", v[2:], err)
		}
		c.Min = ver
	} else if strings.HasPrefix(v, ">") {
		ver, err := semver.NewVersion(v[1:])
		if err != nil {
			return nil, fmt.Errorf("invalid version %q: %w", v[1:], err)
		}
		c.Min = ver
	} else if strings.HasPrefix(v, "<=") {
		ver, err := semver.NewVersion(v[2:])
		if err != nil {
			return nil, fmt.Errorf("invalid version %q: %w", v[2:], err)
		}
		c.Max = ver
	} else if strings.HasPrefix(v, "<") {
		ver, err := semver.NewVersion(v[1:])
		if err != nil {
			return nil, fmt.Errorf("invalid version %q: %w", v[1:], err)
		}
		c.Max = ver
	} else if strings.HasPrefix(v, "^") {
		ver, err := semver.NewVersion(v[1:])
		if err != nil {
			return nil, fmt.Errorf("invalid version %q: %w", v[1:], err)
		}
		c.Min = ver
		nextMajor := semver.MustParse(fmt.Sprintf("%d.0.0", ver.Major()+1))
		c.Max = nextMajor
	} else if strings.HasPrefix(v, "~") {
		ver, err := semver.NewVersion(v[1:])
		if err != nil {
			return nil, fmt.Errorf("invalid version %q: %w", v[1:], err)
		}
		c.Min = ver
		nextMinor := semver.MustParse(fmt.Sprintf("%d.%d.0", ver.Major(), ver.Minor()+1))
		c.Max = nextMinor
	} else {
		ver, err := semver.NewVersion(v)
		if err != nil {
			return nil, fmt.Errorf("invalid version %q: %w", v, err)
		}
		c.Exact = ver
	}

	return c, nil
}

func (c *VersionConstraint) Matches(v *semver.Version) bool {
	if c.Exact != nil {
		return v.Equal(c.Exact)
	}

	if c.Min != nil && v.LessThan(c.Min) {
		return false
	}

	if c.Max != nil && !v.LessThan(c.Max) {
		return false
	}

	return true
}

func ResolveVersion(available []*semver.Version, constraint *VersionConstraint) (*semver.Version, error) {
	sort.Slice(available, func(i, j int) bool {
		return available[i].GreaterThan(available[j])
	})

	for _, v := range available {
		if constraint.Matches(v) {
			return v, nil
		}
	}

	return nil, fmt.Errorf("no version matches constraint %s", constraint.Range)
}
