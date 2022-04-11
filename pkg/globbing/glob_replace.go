package globbing

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

type PatternReplace struct {
	from     string
	to       string
	reFrom   string
	reTo     string
	compiled *regexp.Regexp
}

func NewPatternReplace(from, to string) *PatternReplace {
	return &PatternReplace{
		from: from,
		to:   to,
	}
}

func splitGlobOnPlaceholders(glob string) ([]string, []string) {
	re := regexp.MustCompile(`/\[[^/\[\]]+\]/`)
	phIdx := re.FindAllStringIndex(glob, -1)
	curridx := 0
	maxidx := len(glob)
	var vv []string
	var ph []string
	var cur string
	inph := false
	push := func(r byte) { cur += string(r) }
	save := func() {
		if !inph {
			vv = append(vv, cur)
		} else {
			ph = append(ph, cur)
		}
		cur = ""
	}
	for _, pi := range phIdx {
		for _, ppi := range pi {
			for true {
				if curridx == maxidx {
					break
				}
				if curridx == ppi {
					save()
					break
				}
				push(glob[curridx])
				curridx++
			}
			inph = !inph
		}
	}
	for curridx != maxidx {
		push(glob[curridx])
		curridx++
	}
	save()
	return vv, ph
}

func buildGlobRegexpFromPH(vv, ph []string) string {
	idx := 0
	maxidx := len(vv) + len(ph)
	var re string
	re += `^`
	for idx != maxidx {
		if idx == maxidx-1 {
			suff := filepath.Base(vv[idx/2])
			left := strings.TrimSuffix(vv[idx/2], suff)
			re += buildRegexpFromGlob(left)
			re += "(" + buildRegexpFromGlob(suff) + ")"
		} else if idx%2 == 0 {
			re += buildRegexpFromGlob(vv[idx/2])
		} else {
			re += `(.*)`
		}
		idx++
	}
	re += `$`
	return re
}

func buildGlobReplacementRegexpFromPH(vv, ph, oph []string) string {
	idx := 0
	maxidx := len(vv) + len(ph)
	var re string
	for idx != maxidx {
		if idx == maxidx-1 {
			suff := filepath.Base(vv[idx/2])
			left := strings.TrimSuffix(vv[idx/2], suff)
			re += buildRegexpFromGlob(left)
			re += fmt.Sprintf("$%d", len(ph)+1)
		} else if idx%2 == 0 {
			re += vv[idx/2]
		} else {
			pp := ph[idx/2]
			ii := -1
			for i, v := range oph {
				if v == pp {
					ii = i
				}
			}
			if ii == -1 {
				re += pp
			} else {
				re += fmt.Sprintf("$%d", ii+1)
			}
		}
		idx++
	}
	return re
}

func checkRepGlob(from, to string) bool {
	re := regexp.MustCompile(`^[^/\[\]]*/(([^/\[\]]*|\[[^/\[\]]*\])/)*\*\.[a-zA-Z0-9]+$`)
	if !re.MatchString(from) || !re.MatchString(to) {
		return false
	}
	ff := filepath.Base(from)
	ft := filepath.Base(to)
	if ff != ft {
		return false
	}
	return true
}

func globRepToRegexp(from, to string) (reFrom, reTo string, err error) {
	err = nil
	if !checkRepGlob(from, to) {
		err = fmt.Errorf("Couple '%s' -> '%s' is is not correct", from, to)
		return
	}
	fromParts, fromPH := splitGlobOnPlaceholders(from)
	reFrom = buildGlobRegexpFromPH(fromParts, fromPH)
	toParts, toPH := splitGlobOnPlaceholders(to)
	reTo = buildGlobReplacementRegexpFromPH(toParts, toPH, fromPH)
	return
}

func (p *PatternReplace) Compile() error {
	if p.compiled != nil {
		return nil
	}
	reFrom, reTo, err := globRepToRegexp(p.from, p.to)
	if err != nil {
		return err
	}
	p.reFrom = reFrom
	p.reTo = reTo
	p.compiled = regexp.MustCompile(p.reFrom)
	return nil
}

func (p *PatternReplace) Match(s string) bool {
	p.Compile()
	return p.compiled.MatchString(s)
}

func (p *PatternReplace) Replace(s string) string {
	p.Compile()
	return p.compiled.ReplaceAllString(s, p.reTo)
}
