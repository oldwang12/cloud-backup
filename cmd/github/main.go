package github

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
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
	source           string
	reserve          int
	backupTime       int
	dir              = "/root"
	sizeStr          string
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
			klog.Errorf("%v 不存在, %v", filePath, err)
			return
		}

		backupDir := fmt.Sprintf("%v_%v", source, path.Base(filePath))
		klog.Infof("开始备份文件 %v", backupDir)

		if isDir(filePath) {
			tarFileName, err := tarFile(filePath)
			if err != nil {
				klog.Error(err)
				return
			}

			tarFilePath := filepath.Join(dir, tarFileName)

			tarFilePathInfo, err := os.Stat(tarFilePath)
			if err != nil {
				klog.Errorf("无法获取文件信息: %v", err)
				return
			}
			fileSizeInMegaBytes := float64(tarFilePathInfo.Size()) / 1024 / 1024

			if fileSizeInMegaBytes > 100 {
				klog.Warningf("取消上传文件: %s, 超出 Github 单个文件最大限制 100M, 当前文件大小: %.2fMB\n", filePath, fileSizeInMegaBytes)
				return
			}
			filePath = tarFileName
			backupDir = fmt.Sprintf("%v_%v", source, strings.Split(filePath, ".")[0])
		}

		klog.Infof("文件 %v 大小为 %v", filePath, sizeStr)

		if !fileSizeLessThan(filePath) {
			klog.Warningf("取消上传文件: %s, 文件大小大于 100M\n", filePath)
			return
		}

		backupFileName := generateBackupFileName(filePath, source)
		backupFilePath := path.Join(backupDir, backupFileName)

		if err := g.Upload(filePath, backupFilePath); err != nil {
			klog.Error("上传 %s to %v 失败, %v", filePath, backupFilePath, err)
			return
		}
		klog.Infof("上传 %s to %v 成功.", filePath, backupFilePath)

		// 这里需要等待一段时间，否则删除时 github 可能还没有同步到最新的上传
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

	// 获取 $0
	exePath, err := os.Executable()
	if err != nil {
		klog.Fatalf("获取可执行文件路径时发生错误：%v", err)
	}

	// 列出 /root 下所有文件
	rootfiles, err := listRootDirFiles()
	if err != nil {
		klog.Fatal(err)
	}

	klog.Info("即将备份下列文件...")
	for _, v := range rootfiles {
		klog.Info(v)
	}

	for {
		for _, filePath := range rootfiles {
			// 跳过备份 $0
			if filePath == exePath {
				continue
			}
			backupFunc(filePath)
		}
		time.Sleep(sleepTime)
	}
}

func generateBackupFileName(filePath, source string) string {
	fileName := path.Base(filePath)
	join := ""
	if strings.Contains(fileName, ".") {
		join = "."
	}
	return fmt.Sprintf("%v_%v_%v.%v",
		source,
		strings.Split(fileName, ".")[0],
		time.Now().Format("2006_01_02_150405"),
		strings.Join(strings.Split(fileName, ".")[1:], join),
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
	return nil
}

func listRootDirFiles() ([]string, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		log.Fatal(err)
	}
	return files, nil
}

// 文件或文件夹是否小于100M
func fileSizeLessThan(file string) bool {
	var sizeB float64

	fileInfo, err := os.Stat(file)
	if err != nil {
		klog.Errorf("无法获取文件信息: %v", err)
		return false
	}

	if fileInfo.IsDir() {
		sizeStr, sizeB, err = getDirSize(file)
	} else {
		sizeStr, sizeB, err = getFileSize(file)
	}
	if err != nil {
		klog.Errorf("无法获取文件信息: %v", err)
		return false
	}
	return sizeB/1024/1024 < 100
}

func getDirSize(path string) (string, float64, error) {
	var size float64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += float64(info.Size())
		}
		return err
	})

	sizeB := size
	sizeK := size / 1024
	sizeM := size / 1024 / 1024
	sizeG := size / 1024 / 1024 / 1024

	var returnSize string
	if size > 1024*1024*1024 {
		returnSize = fmt.Sprintf("%vG", fmt.Sprintf("%.2f", sizeG))
	} else if size > 1024*1024 {
		returnSize = fmt.Sprintf("%vM", fmt.Sprintf("%.2f", sizeM))
	} else if size > 1024 {
		returnSize = fmt.Sprintf("%vK", fmt.Sprintf("%.2f", sizeK))
	} else {
		returnSize = fmt.Sprintf("%vB", fmt.Sprintf("%.2f", sizeB))
	}
	return returnSize, sizeB, err
}

func getFileSize(path string) (string, float64, error) {
	var returnSize string

	if !exists(path) {
		return returnSize, 0, fmt.Errorf("file not exists: %v", path)
	}
	fileInfo, err := os.Stat(path)
	if err != nil {
		return returnSize, 0, err
	}

	size := float64(fileInfo.Size())
	sizeB := size
	sizeK := size / 1024
	sizeM := size / 1024 / 1024
	sizeG := size / 1024 / 1024 / 1024

	if size > 1024*1024*1024 {
		returnSize = fmt.Sprintf("%vG", fmt.Sprintf("%.2f", sizeG))
	} else if size > 1024*1024 {
		returnSize = fmt.Sprintf("%vM", fmt.Sprintf("%.2f", sizeM))
	} else if size > 1024 {
		returnSize = fmt.Sprintf("%vK", fmt.Sprintf("%.2f", sizeK))
	} else {
		returnSize = fmt.Sprintf("%vB", fmt.Sprintf("%.2f", sizeB))
	}
	return returnSize, sizeB, nil
}

// exists Whether the path exists
func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
