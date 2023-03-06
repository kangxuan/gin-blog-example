package e

var MsgFlags = map[int]string{
	SUCCESS:                         "ok",
	ERROR:                           "fail",
	INVALID_PARAMS:                  "请求参数错误",
	ERROR_EXIST_TAG:                 "已存在该标签名称",
	ERROR_NOT_EXIST_TAG:             "该标签不存在",
	ERROR_ADD_TAG_FAIL:              "添加标签失败",
	ERROR_EXIST_TAG_FAIL:            "判断标签存在错误",
	ERROR_GET_TAG_FAIL:              "获取标签列表失败",
	ERROR_COUNT_TAG_FAIL:            "获取标签数量失败",
	ERROR_UPDATE_TAG_FAIL:           "更新标签失败",
	ERROR_DELETE_TAG_FAIL:           "删除标签失败",
	ERROR_NOT_EXIST_ARTICLE:         "该文章不存在",
	ERROR_CHECK_EXIST_ARTICLE_FAIL:  "检查文章存在失败",
	ERROR_GET_ARTICLE_FAIL:          "获取文章失败",
	ERROR_COUNT_ARTICLE_FAIL:        "获取文章总数失败",
	ERROR_GET_ARTICLES_FAIL:         "获取文章列表失败",
	ERROR_UPDATE_ARTCLIE_FAIL:       "更新文章失败",
	ERROR_AUTH_CHECK_TOKEN_FAIL:     "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT:  "Token已超时",
	ERROR_AUTH_TOKEN:                "Token生成失败",
	ERROR_AUTH:                      "账号密码错误",
	ERROR_ADD_ARTCLIE_FAIL:          "添加文章失败",
	ERROR_DELETE_ARTCLIE_FAIL:       "删除文章失败",
	ERROR_UPLOAD_CHECK_IMAGE_FORMAT: "图片上传格式错误",
	ERROR_UPLOAD_CHECK_IMAGE_SIZE:   "图片上传大小超限",
	ERROR_UPLOAD_CHECK_IMAGE_FAIL:   "图片上传检查失败",
	ERROR_UPLOAD_SAVE_IMAGE_FAIL:    "图片上传失败",
	ERROR_UPLOAD_FORM_FILE_FAIL:     "图片上传FORM错误",
}

// GetMsg 获取错误信息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return "非法错误"
}
