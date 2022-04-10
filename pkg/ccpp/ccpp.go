package ccpp

import (
	"regexp"
	"strings"

	"github.com/gueckmooh/bs/pkg/functional"
)

var (
	CPPSourceExts                  = []string{"cpp", "cc", "cxx", "C"}
	CPPSourceExtsRe *regexp.Regexp = nil
)

func buildCPPSourceExtsRe() {
	if CPPSourceExtsRe == nil {
		re := `^.*\.(`
		re += strings.Join(CPPSourceExts, "|")
		re += ")$"
		CPPSourceExtsRe = regexp.MustCompile(re)
	}
}

func FilterCPPSourceFiles(files []string) []string {
	buildCPPSourceExtsRe()
	return functional.ListFilter(files,
		func(s string) bool {
			return CPPSourceExtsRe.MatchString(s)
		})
}

func IsCPPSourceFile(file string) bool {
	buildCPPSourceExtsRe()
	return CPPSourceExtsRe.MatchString(file)
}
