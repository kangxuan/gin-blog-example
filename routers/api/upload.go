package api

import (
	"gin-blog-example/pkg/app"
	"gin-blog-example/pkg/e"
	"gin-blog-example/pkg/logging"
	"gin-blog-example/pkg/upload"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UploadImage 上传图片
func UploadImage(c *gin.Context) {
	var (
		data = make(map[string]string)
		appG = app.Gin{C: c}
	)

	// 获取上传的图片 file 是文件句柄 image 是文件头信息
	file, image, err := c.Request.FormFile("image")
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_UPLOAD_FORM_FILE_FAIL, nil)
		return
	}

	if image == nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	imageName := upload.GetImageName(image.Filename)
	fullPath := upload.GetImageFullPath()
	savePath := upload.GetImagePath()

	src := fullPath + imageName
	if !upload.CheckImageExt(imageName) {
		appG.Response(http.StatusOK, e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT, nil)
		return
	} else if !upload.CheckImageSize(file) {
		appG.Response(http.StatusOK, e.ERROR_UPLOAD_CHECK_IMAGE_SIZE, nil)
		return
	} else {
		err := upload.CheckImage(fullPath)
		if err != nil {
			logging.Warn(err)
			appG.Response(http.StatusInternalServerError, e.ERROR_UPLOAD_CHECK_IMAGE_FAIL, nil)
			return
		} else if err := c.SaveUploadedFile(image, src); err != nil { // 保存图片
			logging.Warn(err)
			appG.Response(http.StatusInternalServerError, e.ERROR_UPLOAD_SAVE_IMAGE_FAIL, nil)
			return
		} else {
			data["image_url"] = upload.GetImageFullUrl(imageName)
			data["image_save_url"] = savePath + imageName
		}
	}

	appG.Response(http.StatusOK, e.SUCCESS, data)
}
