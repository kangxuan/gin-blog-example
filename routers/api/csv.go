package api

import (
	"encoding/csv"
	"gin-blog-example/pkg/app"
	"gin-blog-example/pkg/e"
	"gin-blog-example/pkg/export"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

// ExportCsv 测试导出
func ExportCsv(c *gin.Context) {
	var (
		appG        = app.Gin{C: c}
		csvFileName = "test.csv"
	)

	// 创建了一个csv文件
	f, err := os.Create(export.GetExcelFullPath() + csvFileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// \xEF\xBB\xBF 是 UTF-8 BOM 的 16 进制格式，如果不标识 UTF-8 的编码格式的话，写入的汉字会显示为乱码
	f.WriteString("\\xEF\\xBB\\xBF")

	// 生成一个csv写入器
	w := csv.NewWriter(f)
	data := [][]string{
		{"1", "test1", "test1-1"},
		{"2", "test2", "test2-1"},
		{"3", "test3", "test3-1"},
	}

	// 将数据写入写入器并flush
	w.WriteAll(data)

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"csv_url": export.GetExcelFullUrl(csvFileName),
	})
}
