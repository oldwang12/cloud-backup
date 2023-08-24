package github

import (
	"fmt"
	"path"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

type Github struct {
}

var (
	githubRepository string
	// githubRepositoryPath string
	githubToken    string
	githubBranch   string
	githubOwner    string
	backupFilePath string
)

var Command = &cobra.Command{
	Use:   "github",
	Short: "backup to github",
	Long:  `backup to gtihub`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	Command.Flags().StringVar(&githubRepository, "repo", "", "Github Repository,which repository will you backup.")
	// Command.Flags().StringVar(&githubRepositoryPath, "github_path", "", "Github Repository Dir, default: /bakcup.")
	Command.Flags().StringVar(&githubToken, "token", "", "Github Token, create new token see 'https://github.com/settings/tokens/new'.")
	Command.Flags().StringVar(&githubBranch, "branch", "", "Github branch.")
	Command.Flags().StringVar(&githubOwner, "owner", "", "Github owner.")
	Command.Flags().StringVarP(&backupFilePath, "local_filepath", "f", "", "Local file path, example: /root/test.sql.")
}

func run() {
	now := time.Now().Format("2006_01_02_150405")
	backupFilePaths := getBackupFiles()
	for _, localFileName := range backupFilePaths {
		remoteFileName := path.Join(fmt.Sprintf("%v_%v", localFileName, now))

		if err := uploadToGitHub(localFileName, remoteFileName); err != nil {
			klog.Fatal(err)
		}
		klog.Infof("upload %s to %v success.", localFileName, remoteFileName)
	}
}
