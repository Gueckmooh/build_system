package globbing

import (
	"fmt"
	"regexp"
)

type Pattern struct {
	value    string
	regexp   string
	compiled *regexp.Regexp
}

func NewPattern(value string) *Pattern {
	return &Pattern{
		value:    value,
		regexp:   "",
		compiled: nil,
	}
}

func globToRegexp(glob string) string {
	lastIsStar := false
	var re string
	re += `^`
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
	} else if glob[len(glob)-1] == '/' {
		re += `.*`
	}
	re += `$`
	return re
}

func (p *Pattern) Compile() *Pattern {
	if p.compiled != nil {
		return p
	}
	fmt.Printf("globToRegexp(p.value): %v\n", globToRegexp(p.value))
	p.regexp = globToRegexp(p.value)
	p.compiled = regexp.MustCompile(p.regexp)
	return p
}

func (p *Pattern) Match(s string) bool {
	p.Compile()
	return p.compiled.MatchString(s)
}
