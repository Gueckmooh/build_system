package argparse

import (
	"fmt"
	"strings"
)

const (
	prologue = `# -*-shell-script-*-

_%s_completions()
{
`
	inits_prologue = `    local cur prev opts fileopts diropts keywords subcommands sub argopts_match
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    sub="${COMP_WORDS[1]}"
`

	subcommand_list = `    subcommands="%s"
`
	subcommand_test = `    if [[ ${COMP_CWORD} == 1 ]]; then
        COMPREPLY=( $(compgen -W "${subcommands} ${opts}" -- ${cur}) )
        return 0
    fi
`

	gopts_list = `    opts="%s"
`

	gopts_test = `    if [[ ${cur} == * ]] ; then
        COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
        return 0
    fi
`

	gargopts_match_list = `    argopts_match="%s"
`

	gargopts_test = `    if [[ ${prev} =~ ${argopts_match} ]] ; then
        _filedir
        return 0
    fi
`

	subopts_list = `    %s_subopts="%s"
`

	subopts_test = `    if [[ ${sub} == "%s" ]] && [[ ${cur} == * ]] ; then
        COMPREPLY=( $(compgen -W "${opts} ${%s_subopts}" -- ${cur}) )
        return 0
    fi
`

	subargopts_match_list = `    %s_subargopts_match="%s"
`

	subargopts_test = `    if [[ ${sub} == "%s" ]] && ([[ ${prev} =~ ${%s_subargopts_match} ]] || [[ ${prev} =~ ${argopts_match} ]]); then
        _filedir
        return 0
    fi
`

	epilogue = `}
complete -F _%s_completions %s`
)

func contains(l []string, s string) bool {
	for _, v := range l {
		if v == s {
			return true
		}
	}
	return false
}

func (p *Parser) GenBashAutoComplete() {
	script := fmt.Sprintf(prologue, p.Command.name)
	inits := inits_prologue
	var body string

	var gsopts []string
	var glopts []string
	var gsaropts []string
	var glaropts []string
	if len(p.Command.args) > 0 {
		for _, arg := range p.Command.args {
			if arg.sname != "" {
				if arg.size > 1 {
					gsaropts = append(gsaropts, "-"+arg.sname)
				}
				gsopts = append(gsopts, "-"+arg.sname)
			}
			if arg.size > 1 {
				glaropts = append(glaropts, "--"+arg.lname)
			}
			glopts = append(glopts, "--"+arg.lname)
		}
	}

	if len(p.Command.commands) > 0 {
		var cmds []string
		for _, cmd := range p.Command.commands {
			cmds = append(cmds, cmd.name)
		}
		inits += fmt.Sprintf(subcommand_list, strings.Join(cmds, " "))
		body += subcommand_test
		for _, cmd := range p.Command.commands {
			if len(cmd.args) > 0 {
				var sopts []string
				var lopts []string
				var saropts []string
				var laropts []string
				for _, arg := range cmd.args {
					if arg.sname != "" && !contains(gsopts, "-"+arg.sname) {
						if arg.size > 1 {
							saropts = append(saropts, "-"+arg.sname)
						}
						sopts = append(sopts, "-"+arg.sname)
					}
					if !contains(glopts, "--"+arg.lname) {
						if arg.size > 1 {
							laropts = append(laropts, "--"+arg.lname)
						}
						lopts = append(lopts, "--"+arg.lname)
					}
				}
				opts := strings.Trim(strings.Join(sopts, " ")+" "+strings.Join(lopts, " "), " ")
				opts_arg := strings.Trim(strings.Join(saropts, "|")+"|"+strings.Join(laropts, "|"), "|")
				if opts_arg != "" {
					inits += fmt.Sprintf(subargopts_match_list, cmd.name, opts_arg)
					body += fmt.Sprintf(subargopts_test, cmd.name, cmd.name)
				}
				if opts != "" {
					inits += fmt.Sprintf(subopts_list, cmd.name, opts)
					body += fmt.Sprintf(subopts_test, cmd.name, cmd.name)
				}
			}
		}
	}

	if len(p.Command.args) > 0 {
		opts := strings.Trim(strings.Join(gsopts, " ")+" "+strings.Join(glopts, " "), " ")
		opts_arg := strings.Trim(strings.Join(gsaropts, "|")+"|"+strings.Join(glaropts, "|"), "|")
		if opts_arg != "" {
			inits += fmt.Sprintf(gargopts_match_list, opts_arg)
			body += gargopts_test
		}
		if opts != "" {
			inits += fmt.Sprintf(gopts_list, opts)
			body += gopts_test
		}
	}

	script += inits + "\n" + body

	script += fmt.Sprintf(epilogue, p.Command.name, p.Command.name)
	fmt.Println(script)
}
