package git_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/gueckmooh/bs/pkg/git"
)

func TestGitClone(t *testing.T) {
	tmpdir := t.TempDir()
	err := os.Chdir(tmpdir)
	if err != nil {
		t.Fatal(err)
	}
	gr := git.NewGitRepository(
		git.WithUpstreamUrl("https://github.com/gueckmooh/build_system"),
		git.WithPath("here"))
	err = gr.Clone()
	if err != nil {
		t.Fatal(err)
	}
	stat, err := os.Stat("here")
	if err != nil {
		t.Fatal(err)
	}
	if !stat.IsDir() {
		t.Fatal("here must be a directory")
	}
}

func TestGetGitRepoName(t *testing.T) {
	const gitRepoName = "https://github.com/gueckmooh/build_system"
	re := regexp.MustCompile(`.*/([^/]*$)`)
	m := re.FindStringSubmatch(gitRepoName)
	if len(m) < 2 {
		t.Fatal("Match not found")
	}
	if m[1] != "build_system" {
		t.Fatal(`"build_system" not found in the name`)
	}
}

func TestGitCloneRevision(t *testing.T) {
	tmpdir := t.TempDir()
	err := os.Chdir(tmpdir)
	if err != nil {
		t.Fatal(err)
	}
	gr := git.NewGitRepository(
		git.WithUpstreamUrl("https://github.com/gueckmooh/build_system"),
		git.WithPath("here"),
		git.WithRevision("e4c8f6474dcebf9188fb0395f7f870e4807dfa91"),
	)
	err = gr.Clone()
	if err != nil {
		t.Fatal(err)
	}
	err = gr.Checkout()
	if err != nil {
		t.Fatal(err)
	}
	stat, err := os.Stat("here")
	if err != nil {
		t.Fatal(err)
	}
	if !stat.IsDir() {
		t.Fatal("here must be a directory")
	}
}

func TestGitCloneAndCheckout(t *testing.T) {
	tmpdir := t.TempDir()
	err := os.Chdir(tmpdir)
	if err != nil {
		t.Fatal(err)
	}
	gr := git.NewGitRepository(
		git.WithUpstreamUrl("https://github.com/gueckmooh/build_system"),
		git.WithPath("here"),
		git.WithRevision("e4c8f6474dcebf9188fb0395f7f870e4807dfa91"),
	)
	err = gr.CloneAndCheckout()
	if err != nil {
		t.Fatal(err)
	}
	stat, err := os.Stat("here")
	if err != nil {
		t.Fatal(err)
	}
	if !stat.IsDir() {
		t.Fatal("here must be a directory")
	}
}

func TestGitCloneAndCheckoutWithoutName(t *testing.T) {
	tmpdir := t.TempDir()
	err := os.Chdir(tmpdir)
	if err != nil {
		t.Fatal(err)
	}
	gr := git.NewGitRepository(
		git.WithUpstreamUrl("https://github.com/gueckmooh/build_system"),
		git.WithRevision("e4c8f6474dcebf9188fb0395f7f870e4807dfa91"),
	)
	err = gr.CloneAndCheckout()
	if err != nil {
		t.Fatal(err)
	}
	stat, err := os.Stat("build_system")
	if err != nil {
		t.Fatal(err)
	}
	if !stat.IsDir() {
		t.Fatal("here must be a directory")
	}
}
