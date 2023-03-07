package article_service

import (
	"gin-blog-example/pkg/file"
	"gin-blog-example/pkg/qrcode"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
)

// ArticlePoster 海报
type ArticlePoster struct {
	PosterName string
	Article    *Article
	Qr         *qrcode.QrCode
}

// NewArticlePoster 创建一个海报
func NewArticlePoster(postName string, article *Article, qr *qrcode.QrCode) *ArticlePoster {
	return &ArticlePoster{
		PosterName: postName,
		Article:    article,
		Qr:         qr,
	}
}

func GetPosterFlag() string {
	return "poster"
}

// CheckMergedImage 检查图片是否存在
func (a *ArticlePoster) CheckMergedImage(path string) bool {
	if file.CheckExisted(path+a.PosterName) == true {
		return false
	}
	return true
}

// OpenMergedImage 打开海报
func (a *ArticlePoster) OpenMergedImage(path string) (*os.File, error) {
	f, err := file.MustOpen(a.PosterName, path)
	if err != nil {
		return nil, err
	}

	return f, nil
}

type ArticlePosterBg struct {
	Name string
	*ArticlePoster
	*Rect
	*Pt
}

type Rect struct {
	Name string
	X0   int
	Y0   int
	X1   int
	Y1   int
}

type Pt struct {
	X int
	Y int
}

// NewArticlePosterBg 新建一个海报背景
func NewArticlePosterBg(name string, ap *ArticlePoster, rect *Rect, pt *Pt) *ArticlePosterBg {
	return &ArticlePosterBg{
		Name:          name,
		ArticlePoster: ap,
		Rect:          rect,
		Pt:            pt,
	}
}

// Generate 生成海报
func (a *ArticlePosterBg) Generate() (string, string, error) {
	// 获取二维码存储路径
	fullPath := qrcode.GetQrCodeFullPath()
	// 生成二维码图像
	fileName, path, err := a.Qr.Encode(fullPath)
	if err != nil {
		return "", "", err
	}

	// 检查合并后图像（指的是存放合并后的海报）是否存在
	if !a.CheckMergedImage(path) {
		// 若不存在，则生成待合并的图像 mergedF
		mergeDF, err := a.OpenMergedImage(path)
		if err != nil {
			return "", "", err
		}
		defer mergeDF.Close()

		// 打开事先存放的背景图 bgF
		bgF, err := file.MustOpen(a.Name, path)
		if err != nil {
			return "", "", err
		}
		defer bgF.Close()

		// 打开生成的二维码图像 qrF
		qrF, err := file.MustOpen(fileName, path)
		if err != nil {
			return "", "", err
		}
		defer qrF.Close()

		// 解码 bgF 和 qrF 返回 image.Image
		bgImage, err := jpeg.Decode(bgF)
		if err != nil {
			return "", "", err
		}
		qrImage, err := jpeg.Decode(qrF)
		if err != nil {
			return "", "", err
		}

		// 创建一个新的 RGBA 图像
		jpg := image.NewRGBA(image.Rect(a.Rect.X0, a.Rect.Y0, a.Rect.X1, a.Rect.Y1))

		// 在 RGBA 图像上绘制 背景图（bgF）
		draw.Draw(jpg, jpg.Bounds(), bgImage, bgImage.Bounds().Min, draw.Over)
		// 在已绘制背景图的 RGBA 图像上，在指定 Point 上绘制二维码图像（qrF）
		draw.Draw(jpg, jpg.Bounds(), qrImage, qrImage.Bounds().Min.Sub(image.Pt(a.Pt.X, a.Pt.Y)), draw.Over)
		// 将绘制好的 RGBA 图像以 JPEG 4：2：0 基线格式写入合并后的图像文件（mergedF）
		jpeg.Encode(mergeDF, jpg, nil)
	}

	return fileName, path, nil
}
