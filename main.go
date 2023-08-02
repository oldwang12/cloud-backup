package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	filePath := flag.String("filepath", "", "文件路径")
	waitTime := flag.Duration("wait_time", time.Hour, "")
	alistSrc := flag.String("alist_src", "", "")
	alistDst := flag.String("alist_dir", "", "")
	alistHost := flag.String("alist_host", "", "")
	alistDataDir := flag.String("alist_data_dir", "", "")
	alistToken := flag.String("alist_token", "", "")
	flag.Parse()
	check(*filePath, *alistSrc, *alistDst, *alistHost, *alistDataDir, *alistToken)

	timer := time.NewTimer(*waitTime)
	for {
		fmt.Println("start time: ", time.Now().Format("2006-01-02 15:04:05"))
		run(*filePath, *alistSrc, *alistDst, *alistHost, *alistDataDir, *alistToken)
		<-timer.C
	}
}

func run(filePath, alistSrc, alistDst, alistHost, alistDataDir, alistToken string) {
	cmd := exec.Command("sh", "./tar.sh", filePath, alistSrc, alistDst, alistHost, alistDataDir, alistToken)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(output))
}

func check(filePath, alistSrc, alistDst, alistHost, alistDataDir, alistToken string) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		panic(err)
	}

	if alistSrc == "" {
		panic("empty alist src")
	}

	if alistDst == "" {
		panic("empty alist dst")
	}

	if alistHost == "" {
		panic("empty alist host")
	}

	if alistDataDir == "" {
		panic("empty alist data dir")
	}

	if alistToken == "" {
		panic("empty alist token")
	}
}
