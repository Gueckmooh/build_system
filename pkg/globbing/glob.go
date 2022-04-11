package globbing

import (
	"regexp"
	"strings"

	"github.com/gueckmooh/bs/pkg/functional"
)

type Pattern struct {
	value      string
	regexp     string
	compiled   *regexp.Regexp
	rawPattern bool
}

type Patterns []*Pattern

func NewPattern(value string) *Pattern {
	return &Pattern{
		value:      value,
		regexp:     "",
		compiled:   nil,
		rawPattern: false,
	}
}

func NewRawPattern(reg string) *Pattern {
	return &Pattern{
		value:      "",
		regexp:     reg,
		compiled:   nil,
		rawPattern: true,
	}
}

func buildRegexpFromGlob(glob string) string {
	var re string
	lastIsStar := false
	for _, rune := range glob {
		// Special case for * and **
		if rune == '*' {
			if lastIsStar {
				re += `.*`
				lastIsStar = false
			} else {
				lastIsStar = true
			}
			continue
		} else if lastIsStar {
			re += `[^/]*`
			lastIsStar = false
		}
		switch rune {
		case '?':
			re += `.`
		case '.':
			re += `\.`
		case '*':
			re += `\*`
		case '+':
			re += `\+`
		case '(':
			re += `\(`
		case ')':
			re += `\)`
		case '[':
			re += `\[`
		case ']':
			re += `\]`
		case '|':
			re += `\|`
		case '^':
			re += `\^`
		case '$':
			re += `\$`
		default:
			re += string(rune)
		}
	}
	if lastIsStar {
		re += `[^/]*`
		lastIsStar = false
	}
	return re
}

func globToRegexp(glob string) string {
	var re string
	re += `^`
	re += buildRegexpFromGlob(glob)
	re += `$`
	return re
}

func simplifiedGlobToRegexp(glob string) string {
	re := globToRegexp(glob)
	if re[len(re)-2] == '/' {
		re = strings.TrimSuffix(re, "$") + ".*$"
	}
	return re
}

func (p *Pattern) Compile() *Pattern {
	if p.compiled != nil {
		return p
	}
	if p.rawPattern {
		p.compiled = regexp.MustCompile(p.regexp)
	} else {
		p.regexp = simplifiedGlobToRegexp(p.value)
		p.compiled = regexp.MustCompile(p.regexp)
	}
	return p
}

func (p *Pattern) Match(s string) bool {
	p.Compile()
	return p.compiled.MatchString(s)
}

func (p Patterns) Match(s string) bool {
	return functional.ListAnyOf(p, func(p *Pattern) bool { return p.Match(s) })
}
