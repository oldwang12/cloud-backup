package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	fileName := flag.String("filename", "", "")
	srcDir := flag.String("src_dir", "", "")
	dstDir := flag.String("dst_dir", "", "")
	alistHost := flag.String("alist_host", "", "")
	alistToken := flag.String("alist_token", "", "")
	flag.Parse()
	defer time.Sleep(time.Hour)
	if *fileName == "" || *srcDir == "" || *dstDir == "" || *alistHost == "" || *alistToken == "" {
		fmt.Println("缺少参数")
	}

	for {
		time.Sleep(time.Minute)
		fmt.Println("=========================================")
		// 获取当前时间
		currentTime := time.Now().Format("2006-01-02_15h04m")
		fmt.Println("当前时间:", currentTime)

		// 拼接备份文件名
		fileNameTarGz := filepath.Join("/root/data", fmt.Sprintf("%s_%s.tar", currentTime, *fileName))
		fmt.Println("备份文件名:", fileNameTarGz)

		// 创建压缩文件
		err := tarGzFile(*fileName, fileNameTarGz)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("压缩成功")

		// 发送请求
		err = sendRequest(fileNameTarGz, *alistHost, *alistToken, *srcDir, *dstDir)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(fileNameTarGz + "上传百度云成功")
		time.Sleep(6 * time.Hour)
	}
}

func tarGzFile(sourceDir, tarFilePath string) error {
	// 创建压缩文件
	tarFile, err := os.Create(tarFilePath)
	if err != nil {
		return err
	}
	defer tarFile.Close()

	// 创建 tar.Writer
	tarWriter := tar.NewWriter(tarFile)
	defer tarWriter.Close()

	// 压缩文件夹
	err = filepath.Walk(sourceDir, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 获取相对于源文件夹的相对路径
		relPath, err := filepath.Rel(sourceDir, filePath)
		if err != nil {
			return err
		}

		// 创建 tar.Entry
		header, err := tar.FileInfoHeader(fileInfo, "")
		if err != nil {
			return err
		}
		header.Name = relPath

		// 写入 header
		err = tarWriter.WriteHeader(header)
		if err != nil {
			return err
		}

		if !fileInfo.Mode().IsRegular() {
			return nil
		}

		// 写入文件内容
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(tarWriter, file)
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

func gzipWriter(inputFile, outputFile *os.File) error {
	fileStat, err := inputFile.Stat()
	if err != nil {
		return err
	}

	// 创建gzip Writer
	gzipWriter := gzip.NewWriter(outputFile)
	defer gzipWriter.Close()

	gzipWriter.Name = fileStat.Name()
	gzipWriter.ModTime = fileStat.ModTime()

	// 写入文件内容
	_, err = io.Copy(gzipWriter, inputFile)
	if err != nil {
		return err
	}

	return nil
}

func sendRequest(fileNameTarGz, host, token, src, dst string) error {
	url := host + "/api/fs/copy"
	method := "POST"

	payload := fmt.Sprintf(`{
		"src_dir": "%s",
		"dst_dir": "%s",
		"names": ["%s"]
	}`, src, dst, fileNameTarGz)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(payload))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}
