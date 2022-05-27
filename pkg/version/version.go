package version

import (
	"fmt"
	"regexp"
	"strconv"
)

var (
	version_hash string = "v0.0.0-1-g123456789"
	commit_hash  string
	build_time   string
)

var ReleasedVersions = []Version{
	{0, 0, 0}, // v0.0.0 dummy version
	{0, 1, 0}, // v0.1.0
}

type Version struct {
	Major int
	Minor int
	Patch int
}

type ExtendedVersion struct {
	Version
	Commit       string
	CommitsAhead int
	BuildTime    string
}

func (v *Version) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v *ExtendedVersion) String() string {
	if v.CommitsAhead > 0 {
		var s string
		if v.CommitsAhead != 1 {
			s = "s"
		}
		return fmt.Sprintf("version %d.%d.%d commit %s (%d commit%s ahead)", v.Major, v.Minor, v.Patch,
			v.Commit, v.CommitsAhead, s)
	} else {
		return fmt.Sprintf("version %s", &v.Version)
	}
}

var re *regexp.Regexp = regexp.MustCompile(`^v([0-9]+)\.([0-9]+)\.([0-9]+)(-([0-9]+)-([0-9a-zA-Z]+)|)$`)

func ParseVersionHash(hash string) (*ExtendedVersion, error) {
	m := re.FindStringSubmatch(hash)
	if len(m) > 0 {
		var commitsAhead int
		var commit string
		var err error
		if len(m[4]) > 0 {
			commitsAhead, err = strconv.Atoi(m[5])
			if err != nil {
				return nil, err
			}
			commit = m[6]
		}
		major, err := strconv.Atoi(m[1])
		if err != nil {
			return nil, err
		}
		minor, err := strconv.Atoi(m[2])
		if err != nil {
			return nil, err
		}
		patch, err := strconv.Atoi(m[3])
		if err != nil {
			return nil, err
		}
		return &ExtendedVersion{
			Version:      Version{major, minor, patch},
			Commit:       commit,
			CommitsAhead: commitsAhead,
			BuildTime:    build_time,
		}, nil
	}
	return nil, fmt.Errorf("could not parse version")
}

func GetVersion() (*ExtendedVersion, error) {
	v, err := ParseVersionHash(version_hash)
	if err != nil {
		return nil, err
	}
	if v.CommitsAhead == 0 {
		v.Commit = commit_hash
	}
	return v, nil
}
