package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func ReadWbpInfo(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(io.Reader(fi))
	if err != nil {
		fmt.Println(err)
	}
	return string(fd)
}

func WriteWbpInfo(path,content string){
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		//_,err=f.Write([]byte(decoder.ConvertString(content)))
		_,err=f.Write([]byte(content))
	}
}

func ModVersion(content,newVersion string) string {
	oldVersion:=""
	wbpVersion := regexp.MustCompile(`<Version>(.*?)</Version>`)
	if wbpVersion==nil{
		fmt.Println("regexp.MustCompile error")
		return ""
	}
	temp := wbpVersion.FindAllStringSubmatch(content, 1)
	for _,data:=range temp{
		oldVersion=data[1]
	}
	content=strings.Replace(content,oldVersion,newVersion,-1)
	return content
}


func ModWbpVersion(deCompressDir,newVersion string,isModVersionFinished chan<- bool){
	//修改版本号
	content := ModVersion(ReadWbpInfo(GetWbpInfoFile(deCompressDir)), newVersion)
	//写入wbpinfo
	WriteWbpInfo(GetWbpInfoFile(deCompressDir), content)
	isModVersionFinished<- true
}
