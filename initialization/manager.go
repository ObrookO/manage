package initialization

import (
	"manage/models"
	"manage/tool"

	"github.com/astaxie/beego/logs"
)

const (
	ManagerName     = "admin"
	ManagerNickname = "Admin"
	ManagerPassword = "woshiadmin"
)

// InitializeManager 初始化管理员
func InitializeManager() {
	if !models.IsManagerExists(nil) {

		encryptPass := tool.GenerateEncryptedPassword(ManagerPassword)
		if _, err := models.AddManager(models.Manager{
			Username: ManagerName,
			Nickname: ManagerNickname,
			Password: encryptPass,
			IsAdmin:  1,
		}); err != nil {
			logs.Error("Initialize Manager Failed: %v", err)
		}
	}
}
