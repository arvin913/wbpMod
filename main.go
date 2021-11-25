package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"wbpMod/utils"
)

const version string = "1.03"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("请输入要修改的版本号（格式为wbpMod 版本号）")
	} else {
		if os.Args[1] != "" {
			if os.Args[1] == "--version" {
				fmt.Printf("版本号：%s\n", version)
			} else {
				//获取当前文件夹路径
				currentDir, _ := os.Getwd()
				wbpFileList := utils.GetWbpFile(currentDir)
				for _, wbpFile := range wbpFileList {
					//str := strings.Split(wbpFile, "\\")
					//获取当前的文件名
					//fileName := str[len(str)-1]
					fileName := filepath.Base(wbpFile)
					//文件夹名
					dirName := strings.Split(fileName, ".")[0]
					//开始解压缩
					isUnzipFinished := make(chan bool)
					go utils.Unzip(wbpFile, dirName, isUnzipFinished)
					if <-isUnzipFinished {
						fmt.Println("===================================")
						fmt.Printf("%s解压完毕\n", fileName)
					}
					//解压缩之后所在的文件夹
					deCompressDir := utils.GetFileDir(wbpFile) + "\\" + dirName

					isModVersionFinished:= make(chan bool)
					go utils.ModWbpVersion(deCompressDir,os.Args[1],isModVersionFinished)
					if <-isModVersionFinished {
						fmt.Printf("%s修改完毕\n", fileName)
					}

					//压缩文件夹
					isZipFinished := make(chan bool)
					go utils.Zip(deCompressDir, currentDir, dirName, fileName, isZipFinished)
					if <-isZipFinished {
						//time.Sleep(10*1000)
						os.RemoveAll(deCompressDir)
						fmt.Printf("%s压缩完毕\n", fileName)
						fmt.Println("===================================")
					}
				}
			}
		}
	}
}
