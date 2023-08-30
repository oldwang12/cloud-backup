package main

import (
	"github.com/oldwang12/cloud-backup/cmd"
	"k8s.io/klog/v2"
)

func main() {
	version := "2.1.1"
	klog.Infof("当前版本: %v", version)

	cmd.Execute()
}
