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
	Update       int
	Alpha        int
	Beta         int
}

func (v *Version) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v *ExtendedVersion) String() string {
	var u string
	if v.Update != 0 {
		u = fmt.Sprintf(" (update %d)", v.Update)
	} else if v.Alpha != 0 {
		u = fmt.Sprintf(" (alpha %d)", v.Alpha)
	} else if v.Beta != 0 {
		u = fmt.Sprintf(" (beta %d)", v.Beta)
	}
	if v.CommitsAhead > 0 {
		var s string
		if v.CommitsAhead != 1 {
			s = "s"
		}
		return fmt.Sprintf("version %d.%d.%d%s commit %s (%d commit%s ahead)", v.Major, v.Minor, v.Patch, u,
			v.Commit, v.CommitsAhead, s)
	} else {
		return fmt.Sprintf("version %s%s", &v.Version, u)
	}
}

var re *regexp.Regexp = regexp.MustCompile(`^v([0-9]+)\.([0-9]+)\.([0-9]+)(-(update|alpha|beta)(.([0-9]+)|)|)(-([0-9]+)-g([0-9a-zA-Z]+)|)$`)

func ParseVersionHash(hash string) (*ExtendedVersion, error) {
	m := re.FindStringSubmatch(hash)
	if len(m) > 0 {
		var commitsAhead int
		var commit string
		var update int
		var beta int
		var alpha int
		var err error
		if len(m[8]) > 0 {
			commitsAhead, err = strconv.Atoi(m[9])
			if err != nil {
				return nil, err
			}
			commit = m[10]
		}
		if len(m[4]) > 0 {
			var v int = 1
			if len(m[6]) > 0 {
				v, err = strconv.Atoi(m[7])
				if err != nil {
					return nil, err
				}
			}
			switch m[5] {
			case "update":
				update = v
			case "alpha":
				alpha = v
			case "beta":
				beta = v
			}
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
			Update:       update,
			Alpha:        alpha,
			Beta:         beta,
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
