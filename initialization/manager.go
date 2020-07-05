package initialization

import (
	"manage/models"
	"manage/tool"

	"github.com/astaxie/beego"

	"github.com/astaxie/beego/logs"
)

const (
	ManagerName     = "admin"
	ManagerNickname = "Admin"
)

// InitializeManager 初始化管理员
func InitializeManager() {
	if !models.IsManagerExists(nil) {

		encryptPass := tool.GenerateEncryptedPassword(beego.AppConfig.String("defaultpassword"))
		if _, err := models.AddManager(models.Manager{
			Username: ManagerName,
			Nickname: ManagerNickname,
			Email:    beego.AppConfig.String("defaultemail"),
			Password: encryptPass,
			IsAdmin:  1,
		}); err != nil {
			logs.Error("Initialize Manager Failed: %v", err)
		}
	}
}
