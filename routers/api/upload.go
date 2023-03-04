package api

import (
	"gin-blog-example/pkg/e"
	"gin-blog-example/pkg/logging"
	"gin-blog-example/pkg/upload"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UploadImage 上传图片
func UploadImage(c *gin.Context) {
	code := e.SUCCESS
	data := make(map[string]string)

	// 获取上传的图片 file 是文件句柄 image 是文件头信息
	file, image, err := c.Request.FormFile("image")
	if err != nil {
		logging.Warn(err)
		code = e.ERROR
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": data,
		})
		return
	}

	if image == nil {
		code = e.INVALID_PARAMS
	} else {
		imageName := upload.GetImageName(image.Filename)
		fullPath := upload.GetImageFullPath()
		savePath := upload.GetImagePath()

		src := fullPath + imageName
		if !upload.CheckImageExt(imageName) {
			code = e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT
		} else if !upload.CheckImageSize(file) {
			code = e.ERROR_UPLOAD_CHECK_IMAGE_SIZE
		} else {
			err := upload.CheckImage(fullPath)
			if err != nil {
				logging.Warn(err)
				code = e.ERROR_UPLOAD_CHECK_IMAGE_FAIL
			} else if err := c.SaveUploadedFile(image, src); err != nil { // 保存图片
				logging.Warn(err)
				code = e.ERROR_UPLOAD_SAVE_IMAGE_FAIL
			} else {
				data["image_url"] = upload.GetImageFullUrl(imageName)
				data["image_save_url"] = savePath + imageName
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}
