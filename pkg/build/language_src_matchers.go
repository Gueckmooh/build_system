package build

import (
	"github.com/gueckmooh/bs/pkg/globbing"
	"github.com/gueckmooh/bs/pkg/project"
)

func getLanguageSrcMatcher(langID project.LanguageID) *globbing.Pattern {
	switch langID {
	case project.LangCPP:
		return globbing.NewRawPattern(`.*\.(cpp|C|cc|cxx)`)
	}
	return nil
}

func getLanguageHeaderMatcher(langID project.LanguageID) *globbing.Pattern {
	switch langID {
	case project.LangCPP:
		return globbing.NewRawPattern(`.*\.(hpp|h|hh|hxx)`)
	}
	return nil
}
