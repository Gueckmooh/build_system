package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/alessio/shellescape"
	log "github.com/gueckmooh/bs/pkg/logging"
)

const (
	GitBin = "git"
)

type GitRepository struct {
	upstreamURL string
	revision    string
	path        string
}

func extractGitRepoName(url string) (string, error) {
	re := regexp.MustCompile(`.*/([^/]*$)`)
	m := re.FindStringSubmatch(url)
	if len(m) < 2 {
		return "", fmt.Errorf("Match not found")
	}
	return m[1], nil
}

func runCommand(cmd []string) (string, string, error) {
	exe := exec.Command(cmd[0], cmd[1:]...)
	var outb, errb bytes.Buffer
	exe.Stdout = &outb
	exe.Stderr = &errb
	log.Printf("%s\n", shellescape.QuoteCommand(cmd))
	err := exe.Run()
	return outb.String(), errb.String(), err
}

type GitRepositoryOption func(*GitRepository)

func WithUpstreamUrl(url string) GitRepositoryOption {
	return func(g *GitRepository) {
		g.upstreamURL = url
	}
}

func WithRevision(revision string) GitRepositoryOption {
	return func(g *GitRepository) {
		g.revision = revision
	}
}

func WithPath(path string) GitRepositoryOption {
	return func(g *GitRepository) {
		g.path = path
	}
}

func NewGitRepository(opts ...GitRepositoryOption) *GitRepository {
	gr := &GitRepository{}
	for _, opt := range opts {
		opt(gr)
	}
	return gr
}

func (g *GitRepository) Clone() error {
	if len(g.path) == 0 {
		path, err := extractGitRepoName(g.upstreamURL)
		if err != nil {
			return err
		}
		g.path = path
	}
	if len(g.upstreamURL) == 0 {
		return fmt.Errorf("Cannont clone repository without upstream url")
	}
	_, _, err := runCommand([]string{
		GitBin,
		"clone",
		"-n",
		g.upstreamURL,
		g.path,
	})
	if err != nil {
		return err
	}
	return nil
}

func (g *GitRepository) Checkout() error {
	if len(g.path) == 0 {
		return fmt.Errorf("Could not checkout a non existant repository")
	}
	stat, err := os.Stat(g.path)
	if os.IsNotExist(err) || !stat.IsDir() {
		return fmt.Errorf("Could not checkout a non existant repository")
	}
	if len(g.revision) == 0 { // Do nothing
		return nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	err = os.Chdir(g.path)
	defer func() { os.Chdir(cwd) }()
	if err != nil {
		return err
	}
	_, _, err = runCommand([]string{
		GitBin,
		"checkout",
		g.revision,
	})
	if err != nil {
		return err
	}
	return nil
}

func (g *GitRepository) CloneAndCheckout() error {
	err := g.Clone()
	if err != nil {
		return err
	}
	err = g.Checkout()
	if err != nil {
		return err
	}
	return nil
}
