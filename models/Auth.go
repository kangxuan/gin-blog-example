package models

type Auth struct {
	ID       int    `gorm:"primaryKey" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// CheckAuth 检查账号密码
func CheckAuth(username, password string) bool {
	var auth Auth
	db.Where(Auth{Username: username, Password: password}).Select("id").First(&auth)
	if auth.ID > 0 {
		return true
	}
	return false
}
