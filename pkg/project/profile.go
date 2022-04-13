package project

type Profile struct {
	Name string

	cppProfile *CPPProfile

	subProfiles   []*Profile
	parentProfile *Profile
}

func NewProfile(name string) *Profile {
	return &Profile{
		Name:          name,
		cppProfile:    nil,
		subProfiles:   []*Profile{},
		parentProfile: nil,
	}
}

func (p *Profile) NewSubProfile(name string) *Profile {
	np := NewProfile(name)
	np.parentProfile = p
	np.subProfiles = append(np.subProfiles, np)
	return np
}

func (p *Profile) SetCPPProfile(cp *CPPProfile) {
	p.cppProfile = cp
}

func (p *Profile) GetCPPProfile() (cp *CPPProfile) {
	cp = p.cppProfile
	return
}
