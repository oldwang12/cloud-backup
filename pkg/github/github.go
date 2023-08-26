package github

import (
	"context"
	"fmt"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
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

type Content struct {
	Date int64
	SHA  *string
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

func (g *GitHub) Delete(backupDir, backupFileName string, reserve int) error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: g.Token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	opts := &github.RepositoryContentGetOptions{}
	_, contents, _, err := client.Repositories.GetContents(ctx, g.Owner, g.Repository, backupDir, opts)
	if err != nil {
		return err
	}

	backupFiles := []*github.RepositoryContent{}

	newFileName := removeFileNameDate(backupFileName)

	for _, content := range contents {
		if removeFileNameDate(*content.Name) == newFileName {
			backupFiles = append(backupFiles, content)
		}
	}

	if len(backupFiles) <= reserve {
		return nil
	}

	m := make(map[string]Content)
	format := "2006_01_02_150405"

	for _, v := range backupFiles {
		timeStr := strings.Replace(*v.Name, fmt.Sprintf("%v_", strings.Split(newFileName, ".")[0]), "", 1)
		timeStr = strings.TrimSpace(timeStr)
		timeStr = strings.Split(timeStr, ".")[0]

		t, err := time.Parse(format, timeStr)
		if err != nil {
			klog.Error(err)
			return err
		}
		m[*v.Name] = Content{
			Date: t.Unix(),
			SHA:  v.SHA,
		}
	}

	var keyValuePairs []struct {
		Key   string
		Value Content
	}

	for k, v := range m {
		keyValuePairs = append(keyValuePairs, struct {
			Key   string
			Value Content
		}{k, v})
	}

	sort.Slice(keyValuePairs, func(i, j int) bool {
		return keyValuePairs[i].Value.Date < keyValuePairs[j].Value.Date
	})

	deleteFunc := func(fileName string, sha *string) error {
		_, _, err = client.Repositories.DeleteFile(ctx,
			g.Owner,
			g.Repository,
			path.Join(backupDir, fileName),
			&github.RepositoryContentFileOptions{
				Message: github.String(fmt.Sprintf("backup time %v", time.Now().Format("2006-01-02 15:04:05"))),
				SHA:     sha,
				Branch:  github.String(g.Branch),
			})
		if err != nil {
			return err
		}
		klog.Infof("delete file %v success.", path.Join(backupDir, fileName))
		return nil
	}

	deleteCount := len(backupFiles) - reserve
	for _, item := range keyValuePairs {
		if deleteCount <= 0 {
			break
		}
		klog.Infof("%s: %d\n", item.Key, item.Value)
		if err := deleteFunc(item.Key, item.Value.SHA); err != nil {
			return err
		}
		deleteCount--
	}

	return nil
}

func removeFileNameDate(filename string) string {
	// 定义用于匹配日期和时间的正则表达式模式
	pattern := `_\d{4}_\d{2}_\d{2}_\d{6}.`
	// 编译正则表达式模式
	re := regexp.MustCompile(pattern)

	if re.MatchString(filename) {
		// 使用正则表达式替换文件名中的日期和时间部分
		return re.ReplaceAllString(filename, ".")
	}
	klog.Warningf("%v not match.", filename)
	return ""
}
