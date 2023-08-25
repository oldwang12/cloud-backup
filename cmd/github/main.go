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
}

func run() {
	g := github.NewGitHub(githubToken, githubOwner, githubRepository, githubBranch)
	backupFunc := func(filePath string) {
		if isDir(filePath) {
			tarFilePath, err := tarFile(filePath)
			if err != nil {
				klog.Error(err)
				return
			}
			filePath = tarFilePath
		}
		backupFileName := generateBackupFileName(filePath)
		if err := g.Upload(filePath, backupFileName); err != nil {
			klog.Error(err)
			return
		}
		klog.Infof("upload %s to %v success.", filePath, backupFileName)
	}

	for {
		for _, filePath := range strings.Split(backupFilePath, ",") {
			backupFunc(filePath)
		}
		time.Sleep(time.Hour * 12)
	}
}

func generateBackupFileName(filePath string) string {
	fileName := path.Base(filePath)
	return fmt.Sprintf("%v_%v.%v",
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
	tarFilePath := fmt.Sprintf("%s.tar.gz", path.Base(filePath))
	cmd := exec.Command("tar", "-zPcf", tarFilePath, filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return path.Join(path.Dir(filePath), tarFilePath), cmd.Run()
}
