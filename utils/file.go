package utils

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

func GetWbpInfoFile(path string) string {
	wbpInfoFilePath := ""
	fs, _ := ioutil.ReadDir(path)
	for _, file := range fs{
		if !file.IsDir(){
			if strings.Contains(file.Name(), "wbpinfo") {
				wbpInfoFilePath =path+"\\"+ file.Name()
			}
		}
	}
	return wbpInfoFilePath
}


func GetWbpFile(path string) []string{
	wbpFilePath := make([]string, 0)
	fs, _ := ioutil.ReadDir(path)
	for _, file := range fs {
		if !file.IsDir() {
			//fmt.Println(path+"\\"+file.Name())
			if strings.Contains(file.Name(), ".wbp") {
				wbpFilePath = append(wbpFilePath, path+"\\"+file.Name())
			}
		}
	}
	return wbpFilePath
}

func GetFileDir(path string) string {
	return filepath.Dir(path)
}