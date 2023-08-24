package github

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/google/go-github/v38/github"
	"golang.org/x/oauth2"
	"k8s.io/klog/v2"
)

func uploadToGitHub(localFilePath, remoteFilePath string) error {
	// Create a GitHub client using OAuth2 authentication
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
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
		Branch:  github.String(githubBranch),
	}

	_, _, err = client.Repositories.CreateFile(ctx, githubOwner, githubRepository, remoteFilePath, opts)
	if err != nil {
		return err
	}
	return nil
}

func getBackupFiles() []string {
	var backupFilePaths []string
	info, err := os.Stat(path.Join("/root", backupFilePath))
	if err != nil {
		klog.Fatal(err)
	}
	if info.IsDir() {
		klog.Fatal("not support dir")
		// err := filepath.Walk(backupFilePath, func(path string, info os.FileInfo, err error) error {
		// 	if err != nil {
		// 		klog.Fatal(err)
		// 	}

		// 	if !info.IsDir() {
		// 		backupFilePaths = append(backupFilePaths, path)
		// 	}

		// 	return nil
		// })
		// if err != nil {
		// 	klog.Fatal(err)
		// }
	} else {
		backupFilePaths = append(backupFilePaths, backupFilePath)
	}
	return backupFilePaths
}
