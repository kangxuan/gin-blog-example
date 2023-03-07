package file

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
)

// GetSize 获取文件大小
func GetSize(f multipart.File) (int, error) {
	content, err := io.ReadAll(f)
	return len(content), err
}

// GetExt 获取文件的后缀名
func GetExt(fileName string) string {
	return path.Ext(fileName)
}

// CheckExisted 检查文件是否存在
func CheckExisted(src string) bool {
	_, err := os.Stat(src)

	return os.IsNotExist(err)
}

// CheckPermission 检查是否有文件的权限
func CheckPermission(src string) bool {
	_, err := os.Stat(src)

	return os.IsPermission(err)
}

// IsNotExistedMkDir 判断目录是否存在，如果不存在则创建
func IsNotExistedMkDir(src string) error {
	if notExisted := CheckExisted(src); notExisted == true {
		if err := MkDir(src); err != nil {
			return err
		}
	}
	return nil
}

// MkDir 创建目录
func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	return err
}

// Open 打开文件
func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	file, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// MustOpen 打开文件，在之前要判断权限等
func MustOpen(fileName, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}

	src := dir + "/" + filePath
	perm := CheckPermission(src)
	if perm == true {
		return nil, fmt.Errorf("file.CheckPermisson Permission defined scr: %s", src)
	}

	err = IsNotExistedMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistedMkDir src: %s, err: %v", src, err)
	}

	f, err := Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("fail to OpenFile :%v", err)
	}

	return f, nil
}
