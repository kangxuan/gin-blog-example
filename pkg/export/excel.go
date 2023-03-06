package export

import "gin-blog-example/settings"

// GetExcelFullUrl 获取Excel访问地址
func GetExcelFullUrl(name string) string {
	return settings.AppSetting.PrefixUrl + GetExcelPath() + name
}

// GetExcelPath 获取Excel的相对保存目录
func GetExcelPath() string {
	return settings.AppSetting.ExportSavePath
}

// GetExcelFullPath 获取Excel的绝对保存目录
func GetExcelFullPath() string {
	return settings.AppSetting.RuntimeRootPath + GetExcelPath()
}
