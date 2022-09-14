package fileFunc

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

/*
检测文件是否存在
*/
func CheckFileExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

/*
递归创建文件夹
*/
func MakeDir(dirPath string) (err error) {
	if 0 == len(dirPath) {
		err = errors.New("empty dir string")
		return
	}
	sep := string(filepath.Separator)
	for strings.HasSuffix(dirPath, sep) {
		dirPath = dirPath[:len(dirPath)-len(sep)]
	}
	pathArr := strings.Split(dirPath, sep)
	pathLen := len(pathArr)
	for i := 2; i <= pathLen; i++ {
		nowPath := strings.Join(pathArr[:i], sep)
		if CheckFileExist(nowPath) {
			continue
		}
		err = os.Mkdir(nowPath, 0755) // 系统默认文件夹权限，如果需要别的权限创建后可进行修改
		if err != nil {
			return errors.New(nowPath + err.Error())
		}
	}
	return
}

/*
读取文件
*/
func ReadFile(filename string) (val string, err error) {
	if !CheckFileExist(filename) {
		err = errors.New(filename + " is not exist")
		return
	}
	f, err := os.OpenFile(filename, os.O_RDONLY, 0666) // 打开文件
	if nil != err {
		return
	}
	fd, err := ioutil.ReadAll(f)
	f.Close()
	if err != nil {
		return
	}
	val = string(fd)
	return
}

func ReadFileStringLines(filename string, handler func(int, string)) (count int, err error) {
	return ReadFileLines(filename, func(i int, bytes []byte) {
		handler(i, strings.TrimSpace(string(bytes)))
	})
}

/*
读取文件
*/
func ReadFileLines(filename string, handler func(int, []byte)) (count int, err error) {
	if !CheckFileExist(filename) {
		err = errors.New(filename + " is not exist")
		return
	}
	f, err := os.OpenFile(filename, os.O_RDONLY, 0666) // 打开文件
	if nil != err {
		return
	}
	buf := bufio.NewReader(f)
	count = 0
	for {
		line, err := buf.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return count, nil
			}
			return 0, err
		}
		handler(count, line)
		count++
	}
	return
}

/*
读取byte文件
*/
func ReadFileByte(filename string) (val []byte) {
	if !CheckFileExist(filename) {
		return
	}
	f, err := os.Open(filename)
	if nil != err {
		return
	}
	defer f.Close()

	buf := make([]byte, 1024)
	for {
		n, err := f.Read(buf)
		if nil != err && io.EOF != err {
			return
		}
		if 0 == n {
			break
		}
		val = append(val, buf[:n]...)
	}
	return
}

/*
写文件
*/
func WriteFile(filename string, value string) (err error) {
	if !CheckFileExist(filename) {
		// 生成文件
		_, err = os.Create(filename)
	}
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0666) // 打开文件
	if nil != err {
		return
	}
	defer f.Close()
	n, err := io.WriteString(f, value)
	if nil != err {
		return
	}
	if 0 == n {
		err = errors.New("no byte write")
	}
	return
}

func WriteFileCover(filename string, value string) (err error) {
	os.Remove(filename)
	if !CheckFileExist(filename) {
		// 生成文件
		_, err = os.Create(filename)
	}
	f, err := os.OpenFile(filename, os.O_RDWR, 0666) // 打开文件
	if nil != err {
		return
	}
	defer f.Close()
	n, err := io.WriteString(f, value)
	if nil != err {
		return
	}
	if 0 == n {
		err = errors.New("no byte write")
	}
	return
}

/*
写文件,字符类型
*/
func WriteFileByte(filename string, value []byte) (err error) {
	if !CheckFileExist(filename) {
		// 生成文件
		_, err = os.Create(filename)
		if err != nil {
			return
		}
	}
	err = ioutil.WriteFile(filename, value, 0666)
	return
}
func AppendFileByte(filename string, value []byte) (err error) {
	if !CheckFileExist(filename) {
		// 生成文件
		_, err = os.Create(filename)
		if err != nil {
			return
		}
	}
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0666) // 打开文件
	if nil != err {
		return
	}
	defer f.Close()
	n, err := f.Write(value)
	if nil != err {
		return
	}
	if 0 == n {
		err = errors.New("no byte write")
	}
	return
}

func NowPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return dir + string(filepath.Separator)
}

func WriteFileTrunc(actionFilePath, actionValue string) (err error) {
	if !CheckFileExist(actionFilePath) {
		// 生成文件
		_, err = os.Create(actionFilePath)
	}
	actionFile, err := os.OpenFile(actionFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	// 打开文件
	n, err := io.WriteString(actionFile, actionValue)
	if nil != err {
		return
	}
	if 0 == n {
		fmt.Println("no byte write :" + actionFilePath)
	}
	return
}

func ReadJson(value interface{}, file string) (err error) {
	val, err := os.Open(file)
	if err != nil {
		return
	}
	defer val.Close()
	err = json.NewDecoder(val).Decode(&value)
	if err != nil {
		return
	}
	return
}
func WriteJson(value interface{}, file string) (err error) {
	val, err := os.Create(file)
	if err != nil {
		return
	}
	defer val.Close()
	err = json.NewEncoder(val).Encode(&value)
	if err != nil {
		return
	}
	return
}

type ZipDir struct {
	Name     string `json:"name"`
	Children []ZipDir
	Files    []ZipFile
}

type ZipFile struct {
	Name    string `json:"name"`
	Content []byte `json:"content"`
}

func (z ZipDir) CreateZip(path string) (err error) {
	file, err := os.Create(filepath.Join(path))
	if err != nil {
		return
	}
	defer file.Close()
	zipwriter := zip.NewWriter(file)
	defer zipwriter.Close()
	err = z.createTree("", zipwriter)
	return
}

/*
zipwriter.Create("/1/11/1.jpg") 添加一个文件到zip file中，如果包含目录，则目录也会被自动创建。文件路径必须是相对路径，不能以盘符(e.g. C:)开头，而且分隔符要使用斜线(/)而不能使用反斜线(\)，如果只想创建目录而不是文件，可以使用斜线结尾(/name/dir/)。此 Create 方法返回 Writer ，通过这个 Writer 可以将内容写入文件。
文件的内容必须在下一次调用 CreateHeader、Create 或 Close 方法之前全部写入。
*/
func (z ZipDir) createTree(parent string, writer *zip.Writer) (err error) {
	for _, f := range z.Files { // 当前目录下文件写入
		var ioWriter io.Writer
		ioWriter, err = writer.Create("/" + parent + f.Name)
		if err != nil {
			return
		}
		_, err = ioWriter.Write(f.Content)
		if err != nil {
			return
		}
	}
	// 子级写入
	for _, d := range z.Children {
		err = d.createTree(parent+d.Name+"/", writer)
		if err != nil {
			return
		}
	}
	return
}

func CurFileVer() int64 {
	name := os.Args[0]
	f, err := os.Open(name)
	if err != nil {
		fmt.Println("open", name, err.Error())
		return 0
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		fmt.Println("stat", err.Error())
		return 0
	}

	return fi.ModTime().Unix()
}

