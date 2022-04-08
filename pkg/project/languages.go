package project

type LanguageID int8

const (
	LangCPP LanguageID = iota
	LangUnknown
)

func LanguageIDFromString(IDStr string) LanguageID {
	switch IDStr {
	case "CPP":
		return LangCPP
	}
	return LangUnknown
}
