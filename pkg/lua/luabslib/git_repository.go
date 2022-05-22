package luabslib

//go:generate go run ./gen -i ./git_repository.go -c GitRepository -T ./gen/templates -P luabslib -o git_repository_gen.go

import (
	"fmt"

	"github.com/gueckmooh/bs/pkg/git"
)

type GitRepository struct {
	repo *git.GitRepository
}

func NewGitRepository(params map[string]string) (*GitRepository, error) {
	var gitopts []git.GitRepositoryOption
	if v, ok := params["url"]; ok {
		gitopts = append(gitopts, git.WithUpstreamUrl(v))
	} else {
		return nil, fmt.Errorf("Git repository must have an upstream url: please give 'url' param")
	}
	if v, ok := params["path"]; ok {
		gitopts = append(gitopts, git.WithPath(v))
	}
	if v, ok := params["revision"]; ok {
		gitopts = append(gitopts, git.WithRevision(v))
	}
	return &GitRepository{
		repo: git.NewGitRepository(gitopts...),
	}, nil
}

func (g *GitRepository) Clone() error {
	return g.repo.Clone()
}

func (g *GitRepository) Checkout() error {
	return g.repo.Checkout()
}

func (g *GitRepository) CloneAndCheckout() error {
	return g.repo.CloneAndCheckout()
}
