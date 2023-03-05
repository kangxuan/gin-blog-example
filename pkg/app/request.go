package app

import (
	"gin-blog-example/pkg/logging"
	"github.com/astaxie/beego/validation"
)

// MarkErrors 埋点错误信息
func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		logging.Info(err.Key, err.Message)
	}
}
