package file

import (
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
