package github

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/go-github/v38/github"
	"golang.org/x/oauth2"
	"k8s.io/klog/v2"
)

type GitHub struct {
	Token      string
	Owner      string
	Repository string
	Branch     string
}

func NewGitHub(token, owner, repository, branch string) *GitHub {
	return &GitHub{
		Token:      token,
		Owner:      owner,
		Repository: repository,
		Branch:     branch,
	}
}

func (g *GitHub) Upload(localFilePath, remoteFilePath string) error {
	// Create a GitHub client using OAuth2 authentication
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: g.Token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Read the file content
	fileContent, err := os.ReadFile(localFilePath)
	if err != nil {
		return err
	}

	if string(fileContent) == "" {
		klog.Fatal("not support upload empty file")
	}

	// Create a new file in the repository
	opts := &github.RepositoryContentFileOptions{
		Message: github.String(fmt.Sprintf("backup time %v", time.Now().Format("2006-01-02 15:04:05"))),
		Content: fileContent,
		Branch:  github.String(g.Branch),
	}

	_, _, err = client.Repositories.CreateFile(ctx, g.Owner, g.Repository, remoteFilePath, opts)
	if err != nil {
		return err
	}
	return nil
}
