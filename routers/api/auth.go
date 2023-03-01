package api

import (
	"gin-blog-example/models"
	"gin-blog-example/pkg/e"
	"gin-blog-example/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type authValid struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

// GetAuth 登录
func GetAuth(c *gin.Context) {
	var auth models.Auth
	_ = c.BindJSON(&auth)

	//fmt.Println(auth.Username, auth.Password)
	valid := validation.Validation{}
	ok, _ := valid.Valid(&authValid{Username: auth.Username, Password: auth.Password})

	data := make(map[string]interface{})
	code := e.INVALID_PARAMS
	if !ok {
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
		}
	} else {
		if models.CheckAuth(auth.Username, auth.Password) {
			// 账号验证通过后生成token
			token, err := util.GenerateToken(auth.Username, auth.Password)
			if err != nil {
				code = e.ERROR_AUTH_TOKEN
				log.Println(err)
			} else {
				data["token"] = token
				code = e.SUCCESS
			}
		} else {
			code = e.ERROR_AUTH
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}
