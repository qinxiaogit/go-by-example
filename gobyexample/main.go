package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	path, err := getCurrentDir()
	if err != nil {
		fmt.Println("error")
	}
	readDir(path)
	fmt.Println("success")
}
func getCurrentDir() (string, error) {
	file, err := exec.LookPath(os.Args[0])

	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return "", errors.New(`error: Can't find "/" or "\".`)
	}
	return string(path[0 : i+1]), nil
}

func readDir(path string) {
	Dirs, err := ioutil.ReadDir(path)
	if nil != err {

	}
	for _, file := range Dirs {
		if file.IsDir() && !checkHideDir(file.Name()) {
			readDir(path + file.Name())
		}
		if !file.IsDir() {
			delFile(path, file.Name())
		}
	}

}

/*
 *检查隐藏目录
 */
func checkHideDir(dirName string) bool {
	return dirName[0:1] == "."
}

/*
 * 删除文件
 */
func delFile(path string, fileName string) bool {
	index := strings.LastIndex(fileName, ".")
	if index == -1 || fileName[index:] == ".exe" {
		re := os.Remove(path + "/" + fileName)
		fmt.Println(path+"/"+fileName, re)
		return true
	}
	return false
}
