package project

import "fmt"

const (
	DialectCPP98 int8 = iota
	DialectCPP03
	DialectCPP11
	DialectCPP0x
	DialectCPP14
	DialectCPP1y
	DialectCPP17
	DialectCPP1z
	DialectCPP20
	DialectCPP2a
	DialectCPP23
	DialectCPP2b
	DialectCPPGNU98
	DialectCPPGNU03
	DialectCPPGNU11
	DialectCPPGNU0x
	DialectCPPGNU14
	DialectCPPGNU1y
	DialectCPPGNU17
	DialectCPPGNU1z
	DialectCPPGNU20
	DialectCPPGNU2a
	DialectCPPGNU23
	DialectCPPGNU2b
	DialectCPPUnknown
)

type CPPProfile struct {
	Dialect      int8
	BuildOptions []string
}

func NewCPPProfile() *CPPProfile {
	return &CPPProfile{
		Dialect: DialectCPPUnknown,
	}
}

func (p *CPPProfile) Clone() *CPPProfile {
	np := &CPPProfile{
		Dialect:      p.Dialect,
		BuildOptions: p.BuildOptions,
	}
	return np
}

func (p *CPPProfile) Merge(op *CPPProfile) (np *CPPProfile) {
	np = p.Clone()
	np.BuildOptions = append(np.BuildOptions, op.BuildOptions...)
	return np
}

func (p *CPPProfile) SetDialectFromString(s string) error {
	p.Dialect = cppDialectFromString(s)
	if p.Dialect == DialectCPPUnknown {
		return fmt.Errorf("Unknown CPP dialect '%s'", s)
	}
	return nil
}

func cppDialectFromString(s string) int8 {
	switch s {
	case "CPP98":
		return DialectCPP98
	case "CPP03":
		return DialectCPP03
	case "CPP11":
		return DialectCPP11
	case "CPP0x":
		return DialectCPP0x
	case "CPP14":
		return DialectCPP14
	case "CPP1y":
		return DialectCPP1y
	case "CPP17":
		return DialectCPP17
	case "CPP1z":
		return DialectCPP1z
	case "CPP20":
		return DialectCPP20
	case "CPP2a":
		return DialectCPP2a
	case "CPP23":
		return DialectCPP23
	case "CPP2b":
		return DialectCPP2b
	case "GNY98":
		return DialectCPPGNU98
	case "GNU03":
		return DialectCPPGNU03
	case "GNU11":
		return DialectCPPGNU11
	case "GNU0x":
		return DialectCPPGNU0x
	case "GNU14":
		return DialectCPPGNU14
	case "GNU1y":
		return DialectCPPGNU1y
	case "GNU17":
		return DialectCPPGNU17
	case "GNU1z":
		return DialectCPPGNU1z
	case "GNU20":
		return DialectCPPGNU20
	case "GNU2a":
		return DialectCPPGNU2a
	case "GNU23":
		return DialectCPPGNU23
	case "GNU2b":
		return DialectCPPGNU2b
	}
	return DialectCPPUnknown
}
