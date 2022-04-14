package project

type Profile struct {
	Name string

	cppProfile *CPPProfile

	subProfiles   []*Profile
	parentProfile *Profile
	Sources       []FilesPattern
}

func NewProfile(name string) *Profile {
	return &Profile{
		Name:          name,
		cppProfile:    nil,
		subProfiles:   []*Profile{},
		parentProfile: nil,
	}
}

func DummyProfile(name string) *Profile {
	return &Profile{
		Name:          name,
		cppProfile:    NewCPPProfile(),
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

func (p *Profile) Clone() *Profile {
	np := &Profile{
		Name:          p.Name,
		cppProfile:    p.cppProfile.Clone(),
		parentProfile: p.parentProfile,
		subProfiles:   p.subProfiles,
		Sources:       p.Sources,
	}
	return np
}

func (p *Profile) Merge(op *Profile) *Profile {
	np := p.Clone()
	np.Name = op.Name
	np.cppProfile = np.cppProfile.Merge(op.cppProfile)
	np.Sources = append(np.Sources, op.Sources...)
	return np
}

func (p *Profile) AddSubProfile(sp *Profile) {
	p.subProfiles = append(p.subProfiles, sp)
	sp.parentProfile = p
}

func (p *Profile) SetCPPProfile(cp *CPPProfile) {
	p.cppProfile = cp
}

func (p *Profile) GetCPPProfile() (cp *CPPProfile) {
	cp = p.cppProfile
	return
}

func (p *Profile) GetSubProfiles() []*Profile {
	return p.subProfiles
}
