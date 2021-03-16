package info

import (
	"context"
	"errors"
	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
	"os"
	"strings"
)

var GitHub GitHubType

type GitHubType struct {
	client *github.Client
}

// init github client if GITHUB_TOKEN is set
func init() {
	if len(os.Getenv("GITHUB_TOKEN")) > 0 {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
		)
		tc := oauth2.NewClient(ctx, ts)
		GitHub.client = github.NewClient(tc)
	}
}

func (g *GitHubType) GetToken() string {
	return os.Getenv("GITHUB_TOKEN")
}

func (g *GitHubType) GetRepositoryNameWithOwner() (string, error) {
	ownerAndRepo := strings.Split(os.Getenv("GITHUB_REPOSITORY"), "/")
	if len(ownerAndRepo) != 2 {
		return "", errors.New("wrong github repository info")
	}
	return strings.Join(strings.Split(os.Getenv("GITHUB_REPOSITORY"), "/"), ":"), nil
}

func (g *GitHubType) GetRepository() (*github.Repository, error) {
	ownerAndRepo := strings.Split(os.Getenv("GITHUB_REPOSITORY"), "/")
	if len(ownerAndRepo) != 2 {
		return nil, errors.New("wrong github repository info")
	}
	repo, _, err := g.client.Repositories.Get(context.Background(), ownerAndRepo[0], ownerAndRepo[1])
	if err != nil {
		return nil, err
	}
	return repo, nil
}
