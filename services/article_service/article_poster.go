package article_service

import (
	"gin-blog-example/pkg/file"
	"gin-blog-example/pkg/qrcode"
	"gin-blog-example/settings"
	"github.com/golang/freetype"
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

		// 绘制图片
		err = a.DrawPoster(&DrawText{
			JPG:    jpg,
			Merged: mergeDF,

			Title: "Golang Gin 系列文章",
			X0:    80,
			Y0:    160,
			Size0: 42,

			SubTitle: "----山辣",
			X1:       320,
			Y1:       220,
			Size1:    36,
		}, "msyhbd.ttc")
		if err != nil {
			return "", "", err
		}

		// 将绘制好的 RGBA 图像以 JPEG 4：2：0 基线格式写入合并后的图像文件（mergedF）
		jpeg.Encode(mergeDF, jpg, nil)
	}

	return fileName, path, nil
}

type DrawText struct {
	JPG    draw.Image
	Merged *os.File

	Title string
	X0    int
	Y0    int
	Size0 float64

	SubTitle string
	X1       int
	Y1       int
	Size1    float64
}

func (a *ArticlePosterBg) DrawPoster(d *DrawText, fontName string) error {
	// 读取字体库
	fontSource := settings.AppSetting.RuntimeRootPath + settings.AppSetting.FontSavePath + fontName
	fontSourceBytes, err := os.ReadFile(fontSource)
	if err != nil {
		return err
	}

	// 解析字体
	trueTypeFont, err := freetype.ParseFont(fontSourceBytes)
	if err != nil {
		return err
	}

	// 创建一个新的 Context
	fc := freetype.NewContext()
	// 设置屏幕每英寸的分辨率
	fc.SetDPI(72)
	//  设置用于绘制文本的字体
	fc.SetFont(trueTypeFont)
	// 以磅为单位设置字体大小
	fc.SetFontSize(d.Size0)
	// 设置剪裁矩形以进行绘制
	fc.SetClip(d.JPG.Bounds())
	// 设置目标图像
	fc.SetDst(d.JPG)
	// 设置绘制操作的源图像
	fc.SetSrc(image.Black)

	pt := freetype.Pt(d.X0, d.Y0)
	// 根据 Pt 的坐标值绘制给定的文本内容
	_, err = fc.DrawString(d.Title, pt)
	if err != nil {
		return err
	}

	fc.SetFontSize(d.Size1)
	_, err = fc.DrawString(d.SubTitle, freetype.Pt(d.X1, d.Y1))
	if err != nil {
		return err
	}

	err = jpeg.Encode(d.Merged, d.JPG, nil)
	if err != nil {
		return err
	}

	return nil
}
