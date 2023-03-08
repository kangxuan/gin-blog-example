package qrcode

import (
	"gin-blog-example/pkg/file"
	"gin-blog-example/pkg/util"
	"gin-blog-example/settings"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"image/jpeg"
	"os"
)

type QrCode struct {
	URL    string                  // 二维码链接
	Width  int                     // 二维码宽度
	Height int                     // 二维码长度
	Ext    string                  // 二维码图片后缀
	Level  qr.ErrorCorrectionLevel // 二维码允许存储的数据数量
	Mode   qr.Encoding             // 二维码编码方式
}

const (
	EXT_JPG = ".jpg"
)

// NewQrCode 创建一个QrCode
func NewQrCode(url string, width, height int, level qr.ErrorCorrectionLevel, mode qr.Encoding) *QrCode {
	return &QrCode{
		URL:    url,
		Width:  width,
		Height: height,
		Level:  level,
		Mode:   mode,
		Ext:    EXT_JPG,
	}
}

// GetQrCodePath 获取二维码保存相对路径
func GetQrCodePath() string {
	return settings.AppSetting.QrCodeSavePath
}

// GetQrCodeFullPath 获取二维码保存绝对路径
func GetQrCodeFullPath() string {
	return settings.AppSetting.RuntimeRootPath + GetQrCodePath()
}

// GetQrCodeFullUrl 获取二维码的访问路径
func GetQrCodeFullUrl(name string) string {
	return settings.AppSetting.PrefixUrl + GetQrCodePath() + name
}

// GetQrCodeFileName 获取二维码的保存名称，转成md5加密的形式
func GetQrCodeFileName(value string) string {
	return util.EncodeMD5(value)
}

// GetQrCodeExt 获取二维码的后缀
func (q *QrCode) GetQrCodeExt() string {
	return q.Ext
}

// CheckEncode 检查二维码是否存在
func (q *QrCode) CheckEncode(path string) bool {
	src := path + GetQrCodeFileName(q.URL) + q.GetQrCodeExt()
	if file.CheckExisted(src) == true {
		return false
	}
	return true
}

// Encode 生成二维码
func (q *QrCode) Encode(path string) (string, string, error) {
	// 生成二维码的图片名称
	name := GetQrCodeFileName(q.URL) + q.GetQrCodeExt()
	src := path + name
	// 如果二维码不存在则创建
	if file.CheckExisted(src) == true {
		// 返回QR二维码
		code, err := qr.Encode(q.URL, q.Level, q.Mode)
		if err != nil {
			return "", "", err
		}

		// 给二维码设置高宽
		code, err = barcode.Scale(code, q.Width, q.Height)
		if err != nil {
			return "", "", err
		}

		// 生成一个文件
		f, err := file.MustOpen(name, path)
		if err != nil {
			return "", "", err
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				panic(err)
			}
		}(f)

		// 将二维码以 JPEG 4：2：0 基线格式写入文件
		err = jpeg.Encode(f, code, nil)
		if err != nil {
			return "", "", err
		}
	}
	return name, path, nil
}
