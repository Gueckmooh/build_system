package version

import (
	"fmt"
	"regexp"
	"strconv"
)

var (
	version_hash string
	commit_hash  string
	build_time   string
)

type Version struct {
	Major        string
	Minor        string
	Patch        string
	Commit       string
	CommitsAhead int
	BuildTime    string
}

func (v *Version) String() string {
	return fmt.Sprintf("v%s.%s.%s", v.Major, v.Minor, v.Patch)
}

var re *regexp.Regexp = regexp.MustCompile(`^v([0-9]+)\.([0-9]+)\.([0-9]+)(-([0-9]+)-([0-9a-zA-Z]+)|)$`)

func ParseVersionHash(hash string) (*Version, error) {
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
		major := m[1]
		minor := m[2]
		patch := m[3]
		return &Version{
			Major:        major,
			Minor:        minor,
			Patch:        patch,
			Commit:       commit,
			CommitsAhead: commitsAhead,
			BuildTime:    build_time,
		}, nil
	}
	return nil, fmt.Errorf("could not parse version")
}

func GetVersion() (*Version, error) {
	v, err := ParseVersionHash(version_hash)
	if err != nil {
		return nil, err
	}
	if v.CommitsAhead == 0 {
		v.Commit = commit_hash
	}
	return v, nil
}
