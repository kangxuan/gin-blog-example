package upload

import (
	"fmt"
	"gin-blog-example/pkg/file"
	"gin-blog-example/pkg/logging"
	"gin-blog-example/pkg/util"
	"gin-blog-example/settings"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

// GetImagePath 获取图片相对保存目录
func GetImagePath() string {
	return settings.AppSetting.ImageSavePath
}

// GetImageName 获取图片保存的名称
func GetImageName(name string) string {
	// 取图片后缀
	ext := path.Ext(name)
	// 取文件名称
	fileName := strings.TrimPrefix(name, ext)
	// 对文件名称md5加密
	fileName = util.EncodeMD5(fileName)

	return fileName + ext
}

// GetImageFullPath 获取图片绝对保存目录
func GetImageFullPath() string {
	return settings.AppSetting.RuntimeRootPath + GetImagePath()
}

// GetImageFullUrl 获取图片完整的访问路径
func GetImageFullUrl(name string) string {
	return settings.AppSetting.PrefixUrl + GetImagePath() + name
}

// CheckImageExt 检查图片的后缀格式
func CheckImageExt(filename string) bool {
	ext := file.GetExt(filename)
	for _, v := range settings.AppSetting.ImageAllowExts {
		if strings.ToUpper(v) == strings.ToUpper(ext) {
			return true
		}
	}

	return false
}

// CheckImageSize 检查图片大小是否超限
func CheckImageSize(f multipart.File) bool {
	size, err := file.GetSize(f)
	if err != nil {
		log.Println(err)
		logging.Warn(err)
	}

	return size <= settings.AppSetting.ImageMaxSize
}

// CheckImage 检查上传图片所需（权限、文件夹）
func CheckImage(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}

	// 检查图片路径是否存在并创建
	err = file.IsNotExistedMkDir(dir + src)
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir err: %v", err)
	}

	// 检查文件权限
	perm := file.CheckPermission(src)
	if perm == true {
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	return nil
}
