package newluabslib

type Profile struct {
	sources []string
	name    string
	_CPP    *CPPProfile
}

func (p *Profile) CPP() *CPPProfile {
	return p._CPP
}
