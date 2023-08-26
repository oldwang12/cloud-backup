package github

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/oldwang12/cloud-backup/pkg/github"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

type Github struct {
}

var (
	githubRepository string
	githubToken      string
	githubBranch     string
	githubOwner      string
	backupFilePath   string
	source           string
	reserve          int
	backupTime       int
)

var Command = &cobra.Command{
	Use:   "github",
	Short: "Backup to github.",
	Long:  `Backup to gtihub,support one or more (file or directory).`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	Command.Flags().StringVar(&githubRepository, "repo", "", "Github Repository,which repository will you backup.")
	Command.Flags().StringVar(&githubToken, "token", "", "Github Token, create new token see 'https://github.com/settings/tokens/new'.")
	Command.Flags().StringVar(&githubBranch, "branch", "", "Github branch.")
	Command.Flags().StringVar(&githubOwner, "owner", "", "Github owner.")
	Command.Flags().StringVarP(&backupFilePath, "local_filepath", "f", "", "Local file path, support one or more. Example: /root/test1.sql,/root/test2.sql")
	Command.Flags().StringVarP(&source, "source", "s", "backup", "File prefix.")
	Command.Flags().IntVar(&reserve, "reserve", 7, "Reserve")
	Command.Flags().IntVarP(&backupTime, "backuptime", "t", 24, "backup time, 默认单位: Hour,默认值: 24")
}

func run() {
	if err := check(); err != nil {
		klog.Fatal(err)
	}
	g := github.NewGitHub(githubToken, githubOwner, githubRepository, githubBranch)
	backupFunc := func(filePath string) {
		if _, err := os.Stat(filePath); err != nil {
			klog.Errorf("%v not exist, %v", filePath, err)
			return
		}

		backupDir := fmt.Sprintf("%v_%v", source, path.Base(filePath))
		klog.Info(backupDir)

		if isDir(filePath) {
			tarFilePath, err := tarFile(filePath)
			if err != nil {
				klog.Error(err)
				return
			}
			filePath = tarFilePath
			backupDir = fmt.Sprintf("%v_%v", source, strings.Split(filePath, ".")[0])
		}

		backupFileName := generateBackupFileName(filePath, source)
		backupFilePath := path.Join(backupDir, backupFileName)

		if err := g.Upload(filePath, backupFilePath); err != nil {
			klog.Error("upload %s to %v failed, %v", filePath, backupFilePath, err)
			return
		}
		klog.Infof("upload %s to %v success.", filePath, backupFilePath)

		time.Sleep(30 * time.Second)

		if err := g.Delete(backupDir, backupFileName, reserve); err != nil {
			klog.Error(err)
			return
		}
	}

	sleepTime := time.Hour * time.Duration(backupTime)
	klog.Info("repo: ", g.Repository)
	klog.Info("owner: ", g.Owner)
	klog.Info("branch: ", g.Branch)
	klog.Info("source: ", source)
	klog.Info("reserve: ", reserve)
	klog.Info("backupTime: ", sleepTime)

	for {
		for _, filePath := range strings.Split(backupFilePath, ",") {
			backupFunc(filePath)
		}
		time.Sleep(sleepTime)
	}
}

func generateBackupFileName(filePath, source string) string {
	fileName := path.Base(filePath)
	return fmt.Sprintf("%v_%v_%v.%v",
		source,
		strings.Split(fileName, ".")[0],
		time.Now().Format("2006_01_02_150405"),
		strings.Join(strings.Split(fileName, ".")[1:], "."),
	)
}

func isDir(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		klog.Fatal(err)
	}
	return info.IsDir()
}

func tarFile(filePath string) (string, error) {
	tarFileName := fmt.Sprintf("%s.tar.gz", path.Base(filePath))
	cmd := exec.Command("tar", "-zcf", tarFileName, filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// return path.Join(path.Dir(filePath), tarFileName), cmd.Run()
	return tarFileName, cmd.Run()
}

func check() error {
	if githubToken == "" {
		return fmt.Errorf("github token is empty")
	}
	if githubRepository == "" {
		return fmt.Errorf("github repository is empty")
	}
	if githubBranch == "" {
		return fmt.Errorf("github branch is empty")
	}
	if githubOwner == "" {
		return fmt.Errorf("github owner is empty")
	}
	if backupFilePath == "" {
		return fmt.Errorf("backup file path is empty")
	}
	return nil
}
