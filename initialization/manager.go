package initialization

import (
	"manage/models"

	utils "github.com/ObrookO/go-utils"

	"github.com/astaxie/beego"

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
		key := beego.AppConfig.String("aes_key")

		encryptPass, _ := utils.AesEncrypt(ManagerPassword, key)
		if _, err := models.AddManager(models.Manager{
			Username: ManagerName,
			Nickname: ManagerNickname,
			Password: encryptPass,
		}); err != nil {
			logs.Error("Initialize Manager Failed: %v", err)
		}
	}
}
