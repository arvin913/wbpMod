package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var DeCompressDir string
//压缩文件
//files 文件数组，可以是不同dir下的文件或者文件夹
//dest 压缩文件存放地址
func Compress(dirname string, files []*os.File, dest string) error {
	d, _ := os.Create(dest)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	for _, file := range files {
		err := compress(dirname, file, "", w)
		if err != nil {
			return err
		}
	}
	return nil
}

func compress(dirname string, file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	//fmt.Println(info.Name())
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		//fmt.Println("prefix:" + prefix)
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(dirname, f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)

		//去除压缩包内文件夹
		header.Name = strings.Replace(prefix+"/"+header.Name, "/"+dirname, "", 1)
		if strings.Count(header.Name, "/") > 0 {
			header.Name = header.Name[1:len(header.Name)]
		}
		header.Method = zip.Deflate
		//header.Modified = info.ModTime().Add(-10*time.Hour)//不改变文件修改时间
		//fmt.Println("header.Name:" + header.Name)
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

//解压
func DeCompress(zipFile, dest string) (err error) {
	//目标文件夹不存在则创建
	if _, err = os.Stat(dest); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dest, 0755)
		}
	}

	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}

	defer reader.Close()

	for _, file := range reader.File {
		//    log.Println(file.Name)

		if file.FileInfo().IsDir() {

			err := os.MkdirAll(dest+"/"+file.Name, 0755)
			if err != nil {
				log.Println(err)
			}
			continue
		} else {
			err = os.MkdirAll(getDir(dest+"/"+file.Name), 0755)
			if err != nil {
				return err
			}
		}

		rc, err := file.Open()
		if err != nil {
			return err
		}
		//defer rc.Close()

		filename := dest + "/" + file.Name
		//err = os.MkdirAll(getDir(filename), 0755)
		//if err != nil {
		//    return err
		//}

		w, err := os.Create(filename)

		//err=os.Chtimes(filename,file.Modified,file.Modified)
		if err != nil {
			return err
		}
		//defer w.Close()

		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
		w.Close()
		rc.Close()
		//文件原始修改时间
		err=os.Chtimes(GetFileDir(zipFile) + "/"+filename,file.Modified,file.Modified)
	}
	return
}

func getDir(path string) string {
	//return subString(path, 0, strings.LastIndex(path, "/"))
	return GetFileDir(path)
}

func subString(str string, start, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < start || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}

func CompressZip(src string, dest string) (err error) {
	f, err := ioutil.ReadDir(src)
	if err != nil {
		log.Println(err)
	}

	fzip, _ := os.Create(dest)
	w := zip.NewWriter(fzip)

	defer fzip.Close()
	defer w.Close()
	for _, file := range f {
		fw, _ := w.Create(file.Name())
		fileContent, err := ioutil.ReadFile(src + file.Name())
		if err != nil {
			log.Println(err)
		}
		_, err = fw.Write(fileContent)

		if err != nil {
			log.Println(err)
		}
	}
	return
}

func Zip(deCompressDir, currentDir, dirName, fileName string, isZipFinished chan<- bool) {
	f, err := os.Open(deCompressDir)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	files := []*os.File{f}
	outDir := currentDir + "/out/"
	err = os.MkdirAll(outDir, 0755)
	if err != nil {
		fmt.Println(err)
	}
	compressErr := Compress(dirName, files, outDir+fileName)
	if compressErr != nil {
		fmt.Println(compressErr)
	}
	isZipFinished <- true
}

func Unzip(wbpFile, dirName string, isUnzipFinished chan<- bool) {
	deCompressErr := DeCompress(wbpFile, dirName)
	if deCompressErr != nil {
		fmt.Println(deCompressErr)
	}
	isUnzipFinished <- true
}